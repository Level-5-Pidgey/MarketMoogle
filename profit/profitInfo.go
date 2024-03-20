package profitCalc

import (
	"errors"
	"fmt"
	cache "github.com/go-pkgz/expirable-cache"
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
	"github.com/level-5-pidgey/MarketMoogle/db"
	"github.com/level-5-pidgey/MarketMoogle/profit/exchange"
	"math"
	"sort"
	"sync"
	"time"
)

type ProfitCalculator struct {
	Items                    *map[int]*Item
	currencyByObtainMethod   *map[string]map[int]*Item
	currencyByExchangeMethod *map[string]map[int]*Item
	repository               db.Repository
	cache                    cache.Cache
}

const (
	maxLevel             = 90
	noMethodFoundError   = "no exchangeType method found"
	competitionThreshold = 3.0
	salesDayRange        = 7
)

func NewProfitCalculator(
	itemMap *map[int]*Item,
	currencyByObtainMethod *map[string]map[int]*Item,
	currencyByExchangeMethod *map[string]map[int]*Item,
	repo db.Repository,
	cache cache.Cache,
) *ProfitCalculator {
	return &ProfitCalculator{
		currencyByObtainMethod:   currencyByObtainMethod,
		currencyByExchangeMethod: currencyByExchangeMethod,
		Items:                    itemMap,
		repository:               repo,
		cache:                    cache,
	}
}

type SaleMethod struct {
	// What's the method to sell this item?
	ExchangeType string
	// How much currency you're getting from this sale
	Value int
	// How many items you need to sell to get the currency
	Quantity int
	// How much currency are you getting per required item
	ValuePer          int
	SaleVelocity      float64
	CompetitionFactor float64
}

func calculateCompetitionFactor(numOfListings int) float64 {
	const sensitivity = 0.45
	return 1 / (1 + math.Exp(sensitivity*float64(numOfListings)-competitionThreshold))
}

func calculateWeightedMedian(sales *[]*db.Sale) int {
	now := time.Now().UTC()
	maxDiff := float64(0)

	if sales == nil {
		return 0
	}

	salesLength := len(*sales)
	if salesLength == 0 {
		return 0
	}

	if salesLength == 1 {
		return (*sales)[0].PricePer
	}

	type weightedSale struct {
		price     int
		score     float64
		timestamp time.Time
	}
	weightedSales := make([]weightedSale, 0, salesLength)

	for _, sale := range *sales {
		diff := now.Sub(sale.Timestamp).Hours()
		if diff > maxDiff {
			maxDiff = diff
		}
	}

	for _, sale := range *sales {
		diff := now.Sub(sale.Timestamp).Hours()

		weightedSales = append(
			weightedSales, weightedSale{
				price:     sale.PricePer,
				score:     float64(sale.PricePer) * (1 - diff/maxDiff),
				timestamp: sale.Timestamp,
			},
		)
	}

	sort.Slice(
		weightedSales, func(i, j int) bool {
			return weightedSales[i].score < weightedSales[j].score
		},
	)

	return weightedSales[len(weightedSales)/2].price
}

// GetBestSaleMethod
// Get the method of exchange that returns the most gil on this item.
// Includes selling this item on the marketboard
func (p *ProfitCalculator) GetBestSaleMethod(
	item *Item, listings *[]*db.Listing, sales *[]*db.Sale, serverId int, gilOnly bool,
) *SaleMethod {
	cacheKey := fmt.Sprintf("sale_%d_%d", item.Id, serverId)
	if val, found := p.cache.Get(cacheKey); found {
		saleMethod := val.(SaleMethod)

		return &saleMethod
	}
	
	var bestSale *SaleMethod

	competitionFactor := 1.0
	saleVelocity := math.Max(p.salesPerHour(sales, salesDayRange), 0.001)

	if listings != nil {
		// If there's any market listings for this item then see what it's currently being sold for
		if (len(*listings)) > 0 {
			competitionFactor = calculateCompetitionFactor(len(*listings))

			// Get weighted median price for this item (so we can discard overly expensive recipes)
			medianPrice := calculateWeightedMedian(sales)
			for _, listing := range *listings {
				// Only return values on the info's server (as that's the only place they can sell it)
				if listing.WorldId != serverId {
					continue
				}

				// Skip outrageous listing prices (that a player wouldn't realistically buy)
				if listing.PricePer > medianPrice {
					continue
				}

				adjustedPrice := listing.PricePer - 1 // 1 gil undercut
				if adjustedPrice < 1 {
					adjustedPrice = 1
				}

				listingSale := SaleMethod{
					ExchangeType:      readertype.Marketboard,           // TODO put info's world name here, change this to a more complex type
					Value:             adjustedPrice * listing.Quantity, // 1 gil undercut per item
					Quantity:          listing.Quantity,
					ValuePer:          adjustedPrice,
					SaleVelocity:      saleVelocity,
					CompetitionFactor: competitionFactor,
				}

				// Players will (usually) only buy the cheapest listing, so we only update if this is the cheapest
				if bestSale == nil || listingSale.ValuePer < bestSale.ValuePer {
					bestSale = &listingSale
				}
			}
		}

		/*
			If there's no market listings for this item on the info's home world we can generate an average value
			from recent sales
		*/
		if bestSale == nil || bestSale.ValuePer == 0 {
			if sales != nil && len(*sales) > 0 {
				filteredSales := make([]*db.Sale, 0, len(*sales))
				daysAgo := time.Now().AddDate(0, 0, -salesDayRange).UTC()
				for _, sale := range *sales {
					if sale.WorldId == serverId && sale.Timestamp.After(daysAgo) {
						filteredSales = append(filteredSales, sale)
					}
				}

				if len(filteredSales) > 1 {
					// Get average price
					totalSaleValue := 0
					totalQuantity := 0

					for _, sale := range filteredSales {
						totalSaleValue += sale.PricePer
						totalQuantity += sale.Quantity
					}

					averageSale := totalSaleValue / totalQuantity
					averageQuantity := totalQuantity / len(filteredSales)

					historySale := SaleMethod{
						ExchangeType:      readertype.Marketboard, // TODO put info's world name here, change this to a more complex type
						Value:             averageSale * averageQuantity,
						Quantity:          averageQuantity,
						ValuePer:          averageSale,
						SaleVelocity:      saleVelocity,
						CompetitionFactor: competitionFactor,
					}

					if bestSale == nil || historySale.ValuePer > bestSale.ValuePer {
						bestSale = &historySale
					}
				}
			}
		}
	}

	// Get profit from exchanges
	if item.ExchangeMethods != nil {
		exchangeMethods := *item.ExchangeMethods

		for _, exchangeMethod := range exchangeMethods {
			currentMethod := SaleMethod{
				ExchangeType:      exchangeMethod.GetExchangeType(),
				Value:             0,
				Quantity:          0,
				ValuePer:          0,
				SaleVelocity:      saleVelocity,
				CompetitionFactor: competitionFactor,
			}

			switch exchangeMethod.GetExchangeType() {
			case readertype.Gil:
				currentMethod.ExchangeType = readertype.Gil
				currentMethod.Value = exchangeMethod.GetCost()
				currentMethod.Quantity = exchangeMethod.GetQuantity()
				currentMethod.ValuePer = exchangeMethod.GetCost() / exchangeMethod.GetQuantity()
			default:
				if gilOnly {
					continue
				}

				exchangeType := exchangeMethod.GetExchangeType()
				gilValue := 0.0

				cacheKey := fmt.Sprintf("cv_%s_%d", exchangeType, serverId)
				if val, found := p.cache.Get(cacheKey); found {
					gilValue = val.(float64)
				} else {
					// Cache miss: calculate the gil value
					var err error
					gilValue, _, err = p.getGilValueAndBestSaleForCurrency(exchangeType, serverId)
					if err != nil {
						// Handle error, possibly continue to next item
						continue
					}

					// Consider caching the newly obtained gilValue here if applicable
				}

				gilCost := int(float64(exchangeMethod.GetCost()) * gilValue)

				currentMethod.ExchangeType = exchangeType
				currentMethod.Value = gilCost
				currentMethod.Quantity = exchangeMethod.GetQuantity()
				currentMethod.ValuePer = gilCost / exchangeMethod.GetQuantity()
			}

			if bestSale == nil || currentMethod.ValuePer >= bestSale.ValuePer {
				bestSale = &currentMethod
			}
		}
	}

	if bestSale == nil || bestSale.Value == 0 {
		return nil
	}

	// Cache the value of the currency
	if _, exists := p.cache.Peek(cacheKey); !exists {
		p.cache.Set(cacheKey, *bestSale, time.Minute*10)
	}

	return bestSale
}

type ObtainMethod struct {
	ShoppingCart ShoppingCart `json:"ItemsToBuy"`

	// TODO expand this into an object (with a type enum and human readable value)
	ObtainMethod string

	Quantity int

	EffortFactor float64
}

func (o *ObtainMethod) GetCost() int {
	cost := 0
	if o.ShoppingCart.itemsRequired == nil || o.ShoppingCart.ItemsToBuy == nil {
		return 0
	}

	// Since the list of items you're buying might over-buy, we get the Quantity from the actual required item counts
	for _, item := range o.ShoppingCart.ItemsToBuy {
		if quantity, ok := o.ShoppingCart.itemsRequired[item.GetItemId()]; ok {
			cost += item.GetCostPer() * quantity
		}
	}

	return cost
}

func (o *ObtainMethod) GetCostPerItem() float64 {
	return float64(o.GetCost() / o.Quantity)
}

type PurchaseInfo struct {
	ItemId int

	Quantity int

	Server int

	BuyFrom string
}

func isEasierToObtain(curr, new *ObtainMethod) bool {
	if curr == nil {
		return true
	}

	currEffortCost := float64(curr.GetCost()) * curr.EffortFactor
	newEffortCost := float64(new.GetCost()) * new.EffortFactor

	if currEffortCost == newEffortCost {
		return len(new.ShoppingCart.itemsRequired) < len(curr.ShoppingCart.itemsRequired)
	}

	return newEffortCost < currEffortCost
}

func (p *ProfitCalculator) GetCheapestObtainMethod(
	item *Item, numRequired int, listings *[]*db.Listing, player *PlayerInfo, skipExchanges bool,
) *ObtainMethod {
	var cheapestMethod *ObtainMethod

	if item.ObtainMethods != nil {
		cheapestMethod = p.nonMarketObtainMethod(item, numRequired, cheapestMethod, player, skipExchanges)
	}

	if !item.MarketProhibited && listings != nil {
		var filteredListings []*db.Listing
		for _, listing := range *listings {
			if listing.ItemId == item.Id {
				filteredListings = append(filteredListings, listing)
			}
		}

		if len(filteredListings) != 0 {
			cheapestMethod = marketObtainMethod(item, cheapestMethod, numRequired, &filteredListings, player)
		}
	}

	if item.CraftingRecipes != nil {
		cheapestMethod = p.craftingObtainMethod(item, numRequired, listings, cheapestMethod, player, skipExchanges)
	}

	return cheapestMethod
}

func (p *ProfitCalculator) getPossibleSubItems(
	itemsAndQuantities map[int]struct{}, item *Item, skipCrystals bool,
) map[int]struct{} {
	if itemsAndQuantities == nil {
		newMap := make(map[int]struct{})
		itemsAndQuantities = newMap
	}

	for _, recipe := range *item.CraftingRecipes {
		for _, ingredient := range recipe.RecipeIngredients {
			ingredientItem, ok := (*p.Items)[ingredient.ItemId]

			if !ok {
				continue
			}

			// TODO remove this magic number, get correct value dynamically from csv load
			if ingredientItem.UiCategory == 59 && skipCrystals {
				continue
			}

			if ingredientItem.CraftingRecipes != nil {
				itemsAndQuantities = p.getPossibleSubItems(
					itemsAndQuantities,
					ingredientItem,
					skipCrystals,
				)
			} else {
				// There is no sub recipe, we should be updating the itemsAndQuantities map
				if _, ok := itemsAndQuantities[ingredient.ItemId]; !ok {
					itemsAndQuantities[ingredient.ItemId] = struct{}{} // Empty struct with size of 0
				}
			}
		}
	}

	return itemsAndQuantities
}

func (p *ProfitCalculator) craftingObtainMethod(
	item *Item, numRequired int, listings *[]*db.Listing, cheapestMethod *ObtainMethod, player *PlayerInfo,
	skipExchanges bool,
) *ObtainMethod {
	for _, craftingRecipe := range *item.CraftingRecipes {
		// Check if the player is capable of crafting this recipe
		if jobLevel, ok := player.JobLevels[craftingRecipe.JobRequired]; ok {
			if jobLevel < craftingRecipe.RecipeLevel {
				continue
			}
		} else {
			// If they don't have the job at all, skip as well
			continue
		}

		recipeCost := ObtainMethod{
			ShoppingCart: ShoppingCart{
				ItemsToBuy:    []ShoppingItem{},
				itemsRequired: make(map[int]int),
			},
			Quantity:     craftingRecipe.Yield,
			EffortFactor: recipeEffort(&craftingRecipe, item),
			ObtainMethod: fmt.Sprintf("Craft with %s", craftingRecipe.JobRequired), // TODO expand with type
		}

		canCraft := true
		for _, ingredient := range craftingRecipe.RecipeIngredients {
			ingredientItem, ok := (*p.Items)[ingredient.ItemId]
			if !ok {
				continue
			}

			if player.SkipCrystals && ingredientItem.UiCategory == 59 && ingredientItem.Id < 20 {
				continue
			}

			ingredientQuantity := (ingredient.Quantity*numRequired + craftingRecipe.Yield - 1) / craftingRecipe.Yield

			// Get the best way to obtain this ingredient
			ingredientObtain := p.GetCheapestObtainMethod(
				ingredientItem,
				ingredientQuantity,
				listings,
				player,
				skipExchanges,
			)

			// Automatically skip out of unobtainable or expensive ingredients
			if ingredientObtain == nil || !isEasierToObtain(cheapestMethod, ingredientObtain) {
				canCraft = false
				break
			}

			// Merge shopping carts together
			recipeCost.ShoppingCart.mergeWith(ingredientObtain.ShoppingCart)
		}

		if canCraft && isEasierToObtain(cheapestMethod, &recipeCost) {
			cheapestMethod = &recipeCost
		}
	}

	return cheapestMethod
}

func marketObtainMethod(
	item *Item, cheapestMethod *ObtainMethod, numRequired int, listings *[]*db.Listing, player *PlayerInfo,
) *ObtainMethod {
	sortListings := func(listings []*db.Listing) {
		sort.Slice(
			listings, func(i, ii int) bool {
				listingAEffortCost := calculateListingEffortCost(listings[i], player.HomeServer)
				listingBEffortCost := calculateListingEffortCost(listings[ii], player.HomeServer)

				// Tiebreaker logic
				if listingAEffortCost == listingBEffortCost {
					if listings[i].Total == listings[ii].Total {
						return listings[i].Id < listings[ii].Id
					}

					return listings[i].PricePer < listings[ii].PricePer
				}

				return listingAEffortCost < listingBEffortCost
			},
		)
	}

	sortListings(*listings)

	purchasePlan := ObtainMethod{
		ShoppingCart: ShoppingCart{
			ItemsToBuy: []ShoppingItem{},
			itemsRequired: map[int]int{
				item.Id: numRequired,
			},
		},
		ObtainMethod: "Market",
		Quantity:     0,
		EffortFactor: 0.99,
	}

	for _, listing := range *listings {
		if numRequired <= 0 {
			break
		}

		if alreadyBoughtListing(cheapestMethod, listing) {
			continue
		}

		numRequired -= listing.Quantity
		purchasePlan.ShoppingCart.ItemsToBuy = append(
			purchasePlan.ShoppingCart.ItemsToBuy,
			ShoppingListing{
				ItemId:       item.Id,
				Quantity:     listing.Quantity,
				RetainerName: listing.RetainerName,
				listingId:    listing.Id,
				worldId:      listing.WorldId,
				CostPer:      listing.PricePer,
			},
		)

		if player.HomeServer == listing.WorldId {
			purchasePlan.EffortFactor += 0.01
		} else {
			purchasePlan.EffortFactor += 0.06
		}

		purchasePlan.Quantity += listing.Quantity
	}

	if isEasierToObtain(cheapestMethod, &purchasePlan) {
		cheapestMethod = &purchasePlan
	}

	return cheapestMethod
}

func alreadyBoughtListing(cheapestMethod *ObtainMethod, listing *db.Listing) bool {
	if cheapestMethod == nil {
		return false
	}

	for _, cartItem := range cheapestMethod.ShoppingCart.ItemsToBuy {
		if cartItem.GetHash() == listing.UniversalisId {
			return true
		}
	}

	return false
}

func calculateListingEffortCost(listing *db.Listing, playerServer int) float64 {
	listingScore := 0.99

	if listing.WorldId == playerServer {
		listingScore += 0.01
	} else {
		listingScore += 0.06
	}

	return math.Round(float64(listing.Total) * listingScore)
}

func (p *ProfitCalculator) nonMarketObtainMethod(
	item *Item, numRequired int, cheapestMethod *ObtainMethod, info *PlayerInfo, skipExchanges bool,
) *ObtainMethod {
	for _, obtainMethod := range *item.ObtainMethods {
		obtainCost := 1500.0

		switch obtainMethod.GetExchangeType() {
		case readertype.GrandCompanySeal:
			if info.GrandCompanyRank < obtainMethod.(exchange.GcSealExchange).RankRequired {
				continue
			}
		case readertype.Gathering:
			obtainCost = obtainMethod.GetCostPerItem()
			break
		case readertype.Gil:
			obtainCost = obtainMethod.GetCostPerItem()
			break
		default:
			// Prevent recursive currency searches
			if skipExchanges {
				break
			}

			currencyObtain, err := p.getCheapestMethodToObtainCurrency(
				obtainMethod.GetExchangeType(),
				info,
			)

			if err == nil {
				obtainCost = currencyObtain.GetCostPerItem() * float64(obtainMethod.GetQuantity())
			}
		}

		numOfExchanges := 1
		totalEffort := obtainMethod.GetEffortFactor()
		for (numOfExchanges * obtainMethod.GetQuantity()) < numRequired {
			numOfExchanges++
			totalEffort *= 1.01
		}
		totalQuantity := numOfExchanges * obtainMethod.GetQuantity()

		currentMethod := ObtainMethod{
			ShoppingCart: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					LocalItem{
						ItemId:       item.Id,
						Quantity:     totalQuantity,
						ObtainedFrom: obtainMethod.GetObtainDescription(), // TODO add npc name here
						CostPer:      obtainCost,
					},
				},
				itemsRequired: map[int]int{
					item.Id: numRequired,
				},
			},
			Quantity:     totalQuantity,
			EffortFactor: totalEffort,
			ObtainMethod: obtainMethod.GetExchangeType(),
		}

		if isEasierToObtain(cheapestMethod, &currentMethod) {
			cheapestMethod = &currentMethod
		}
	}

	return cheapestMethod
}

func recipeEffort(recipe *RecipeInfo, item *Item) float64 {
	result := 1.0
	// Various effort penalties and bonuses depending on recipe requirements
	if recipe.SecretRecipeBook != "" {
		result *= 1.02
	}

	if recipe.SpecializationRequired {
		result *= 1.05
	}

	if recipe.RecipeLevel <= (maxLevel-10) &&
		!recipe.IsExpert &&
		item.CanBeHq {
		result *= 0.95
	}

	if !item.CanBeHq {
		result *= 0.95
	}

	if recipe.IsExpert {
		result *= 1.2
	}

	return result
}

func combinePurchaseInfo(slice1, slice2 []*PurchaseInfo) []*PurchaseInfo {
	combined := make([]*PurchaseInfo, 0)
	itemMap := make(map[string]*PurchaseInfo)

	for _, pInfo := range append(slice1, slice2...) {
		key := fmt.Sprintf("%v-%d", pInfo.ItemId, pInfo.Server) // Assuming Item has an ID field
		if existing, found := itemMap[key]; found {
			existing.Quantity += pInfo.Quantity
		} else {
			newItem := *pInfo // Make a copy to avoid modifying the original slice
			itemMap[key] = &newItem
			combined = append(combined, &newItem)
		}
	}

	return combined
}

type ProfitInfo struct {
	ItemId       int
	ObtainMethod *ObtainMethod
	SaleMethod   *SaleMethod
	ProfitScore  float64
}

func (p *ProfitCalculator) salesPerHour(sales *[]*db.Sale, dayRange int) float64 {
	if sales == nil || len(*sales) == 0 {
		return 0
	}

	daysAgo := time.Now().AddDate(0, 0, -dayRange).UTC()
	var filteredSales []*db.Sale
	for _, sale := range *sales {
		if sale.Timestamp.After(daysAgo) {
			filteredSales = append(filteredSales, sale)
		}
	}

	return float64(len(filteredSales)) / float64(dayRange)
}

func (p *ProfitCalculator) CalculateProfitForItem(item *Item, info *PlayerInfo) (*ProfitInfo, error) {
	// Pre-calculate all items that could be involved in the obtaining of this item
	itemIds := make([]int, 0, 10)
	if item.CraftingRecipes != nil {
		itemMap := p.getPossibleSubItems(nil, item, info.SkipCrystals)
		for itemId := range itemMap {
			itemIds = append(itemIds, itemId)
		}
	}

	// Add the main item itself if not market prohibited
	if !item.MarketProhibited {
		itemIds = append(itemIds, item.Id)
	}

	// Get market listings for item if this item is sellable
	var listings *[]*db.Listing = nil
	var listingsOnPlayerWorld []*db.Listing
	if len(itemIds) > 0 {
		listingResults, err := p.repository.GetListingsForItemsOnDataCenter(itemIds, info.DataCenter)
		listings = listingResults
		if err != nil {
			return nil, err
		}

		for _, listing := range *listings {
			if listing.WorldId == info.HomeServer && listing.ItemId == item.Id {
				listingsOnPlayerWorld = append(listingsOnPlayerWorld, listing)
			}
		}
	}

	sales, err := p.repository.GetSalesForItemOnDataCenter(item.Id, info.DataCenter)
	if err != nil {
		return nil, err
	}

	// Get most value created when selling the item
	bestSale := p.GetBestSaleMethod(item, &listingsOnPlayerWorld, sales, info.HomeServer, false)
	if bestSale == nil {
		// Sometimes there's no way to sell this item, and that's okay. We will just return early
		return nil, nil
	}

	// Get the cheapest method to obtain the item
	cheapestMethod := p.GetCheapestObtainMethod(item, bestSale.Quantity, listings, info, false)
	if cheapestMethod == nil {
		return nil, nil
	}

	// Return info
	return &ProfitInfo{
		ItemId:       item.Id,
		ObtainMethod: cheapestMethod,
		SaleMethod:   bestSale,
		ProfitScore:  calculateProfitScore(bestSale, cheapestMethod.GetCost(), cheapestMethod.EffortFactor),
	}, nil
}

func calculateProfitScore(bestSale *SaleMethod, cost int, effort float64) float64 {
	profitMargin := bestSale.Value - cost

	if profitMargin < 0 {
		profitMargin = 0.0
	}

	if bestSale.ExchangeType != "Marketboard" {
		return float64(profitMargin) / effort
	}

	adjustedProfit := float64(profitMargin) * bestSale.SaleVelocity
	profitScore := (adjustedProfit * bestSale.CompetitionFactor) / effort

	return profitScore
}

func (p *ProfitCalculator) getGilValueAndBestSaleForCurrency(currency string, serverId int) (
	float64, *SaleMethod, error,
) {
	if p.currencyByObtainMethod == nil {
		return 0, nil, errors.New("map of currencies by currency method is nil")
	}

	itemsWithObtainMethod, ok := (*p.currencyByObtainMethod)[currency]
	if !ok {
		return 0, nil, errors.New("no currency method found for currency")
	}

	type scoreAndSale struct {
		score       float64
		perCurrency float64
		item        *Item
		sale        *SaleMethod
	}

	itemIds := make([]int, 0, len(itemsWithObtainMethod))
	for _, item := range itemsWithObtainMethod {
		if item.MarketProhibited {
			continue
		}

		itemIds = append(itemIds, item.Id)
	}

	if len(itemIds) == 0 {
		return 0, nil, nil
	}

	listings, err := p.repository.GetListingsForItemsOnWorld(itemIds, serverId)
	if err != nil {
		return 0, nil, nil
	}

	sales, err := p.repository.GetSalesForItemsOnWorld(itemIds, serverId)
	if err != nil {
		return 0, nil, nil
	}

	wg := sync.WaitGroup{}

	scoreChan := make(chan scoreAndSale)

	for _, item := range itemsWithObtainMethod {
		wg.Add(1)

		go func(item *Item, listings *[]*db.Listing, sales *[]*db.Sale, serverId int) {
			defer wg.Done()

			filteredSales := make([]*db.Sale, 0, 100)
			filteredListings := make([]*db.Listing, 0, 100)
			for _, sale := range *sales {
				if sale.ItemId != item.Id {
					continue
				}

				filteredSales = append(filteredSales, sale)
			}

			for _, listing := range *listings {
				if listing.ItemId != item.Id {
					continue
				}

				filteredListings = append(filteredListings, listing)
			}

			itemSale := p.GetBestSaleMethod(item, &filteredListings, &filteredSales, serverId, true)

			if itemSale == nil {
				return
			}

			cost := 0
			effortFactor := 1.0
			adjustedQuantity := 1.0
			for _, obtainMethod := range *item.ObtainMethods {
				if obtainMethod.GetExchangeType() != currency {
					continue
				}

				cost = obtainMethod.GetCost()
				effortFactor = obtainMethod.GetEffortFactor()
				adjustedQuantity = obtainMethod.GetCostPerItem()

				break
			}

			currentScore := calculateProfitScore(itemSale, cost, effortFactor)
			score := scoreAndSale{
				score:       currentScore,
				perCurrency: float64(itemSale.ValuePer) / adjustedQuantity,
				item:        item,
				sale:        itemSale,
			}

			scoreChan <- score
		}(item, listings, sales, serverId)
	}

	go func() {
		wg.Wait()
		close(scoreChan)
	}()

	best := scoreAndSale{}
	for result := range scoreChan {
		if result.score > best.score {
			best = result
		}
	}

	// Cache the value of the currency
	cacheKey := fmt.Sprintf("cv_%s_%d", currency, serverId)
	if _, exists := p.cache.Peek(cacheKey); !exists {
		p.cache.Set(cacheKey, best.perCurrency, time.Minute*10)
	}

	return best.perCurrency, best.sale, nil
}

func (p *ProfitCalculator) GetBestItemToSellForCurrency(currency string, serverId int) (*SaleMethod, error) {
	_, sale, err := p.getGilValueAndBestSaleForCurrency(currency, serverId)

	if err != nil {
		return nil, err
	}

	return sale, nil
}

func (p *ProfitCalculator) GetGilValueForCurrency(currency string, serverId int) (float64, error) {
	if p.cache != nil {
		cacheKey := fmt.Sprintf("cv_%s_%d", currency, serverId)
		if val, found := p.cache.Get(cacheKey); found {
			return val.(float64), nil
		}
	}

	gilEquivalent, _, err := p.getGilValueAndBestSaleForCurrency(currency, serverId)

	if err != nil {
		return 0.0, err
	}

	return gilEquivalent, nil
}

func (p *ProfitCalculator) getCheapestMethodToObtainCurrency(currency string, info *PlayerInfo) (
	*ObtainMethod, error,
) {
	if p.currencyByExchangeMethod == nil {
		return nil, errors.New("map of currencies by obtain method is nil")
	}

	itemsWithExchangeMethod, ok := (*p.currencyByExchangeMethod)[currency]
	if !ok {
		return nil, errors.New(noMethodFoundError)
	}

	if p.cache != nil {
		cacheKey := fmt.Sprintf("cc_%s_%d", currency, info.HomeServer)
		if val, found := p.cache.Get(cacheKey); found {
			cachedMethod := val.(ObtainMethod)
			return &cachedMethod, nil
		}
	}

	itemIds := make([]int, 0, len(itemsWithExchangeMethod))
	for _, item := range itemsWithExchangeMethod {
		if item.MarketProhibited {
			continue
		}

		itemIds = append(itemIds, item.Id)
	}

	listingResults, err := p.repository.GetListingsForItemsOnDataCenter(itemIds, info.DataCenter)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	type obtainWithCost struct {
		method *ObtainMethod
		cost   float64
	}
	obtainChan := make(chan obtainWithCost)

	for _, item := range itemsWithExchangeMethod {
		wg.Add(1)

		go func(item *Item, listings *[]*db.Listing, info *PlayerInfo) {
			defer wg.Done()

			itemMethod := p.GetCheapestObtainMethod(item, 1, listingResults, info, true)

			if itemMethod == nil {
				return
			}

			costPerToken := 0.0
			for _, exchangeMethod := range *item.ExchangeMethods {
				if exchangeMethod.GetExchangeType() != currency {
					continue
				}

				costPerToken = float64(itemMethod.GetCost()) / float64(exchangeMethod.GetCost())

				break
			}

			obtainChan <- obtainWithCost{
				method: itemMethod,
				cost:   costPerToken,
			}
		}(item, listingResults, info)
	}

	go func() {
		wg.Wait()
		close(obtainChan)
	}()

	bestCost := 0.0
	var bestMethod *ObtainMethod
	for result := range obtainChan {
		if bestMethod == nil || result.cost < bestCost {
			bestCost = result.cost
			bestMethod = result.method
		}
	}

	// Cache the value of the currency
	cacheKey := fmt.Sprintf("cc_%s_%d", currency, info.HomeServer)
	if _, exists := p.cache.Peek(cacheKey); !exists {
		p.cache.Set(cacheKey, *bestMethod, time.Minute*10)
	}

	return bestMethod, nil
}

func (p *ProfitCalculator) GetMinCostOfCurrency(currency string, info *PlayerInfo) (int, error) {
	cheapestMethod, err := p.getCheapestMethodToObtainCurrency(currency, info)

	// It's okay if there's no way to obtain this currency
	if err != nil && err.Error() == noMethodFoundError {
		if err.Error() == noMethodFoundError {
			return 0, nil
		}

		return 0, err
	}

	return cheapestMethod.GetCost(), nil
}

func (p *ProfitCalculator) GetCheapestAcquisitionMethodForCurrency(currency string, info *PlayerInfo) (
	*ObtainMethod, error,
) {
	cheapestMethod, err := p.getCheapestMethodToObtainCurrency(currency, info)

	if err != nil {
		return nil, err
	}

	if cheapestMethod == nil {
		return nil, errors.New("could not find a cheap method to obtain this currency")
	}

	return cheapestMethod, nil
}
