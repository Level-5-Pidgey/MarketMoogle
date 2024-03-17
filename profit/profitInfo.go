package profitCalc

import (
	"errors"
	"fmt"
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
}

const (
	maxLevel           = 90
	noMethodFoundError = "no exchangeType method found for exchangeType"
)

func NewProfitCalculator(
	itemMap *map[int]*Item,
	currencyByObtainMethod *map[string]map[int]*Item,
	currencyByExchangeMethod *map[string]map[int]*Item,
	repo db.Repository,
) *ProfitCalculator {
	return &ProfitCalculator{
		currencyByObtainMethod:   currencyByObtainMethod,
		currencyByExchangeMethod: currencyByExchangeMethod,
		Items:                    itemMap,
		repository:               repo,
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

// GetBestSaleMethod
// Get the method of exchange that returns the most gil on this item.
// Includes selling this item on the marketboard
func (p *ProfitCalculator) GetBestSaleMethod(
	item *Item, listings *[]*db.Listing, sales *[]*db.Sale, info *PlayerInfo, gilOnly bool,
) *SaleMethod {
	var bestSale *SaleMethod

	competitionFactor := 1.0
	saleVelocity := math.Max(p.salesPerHour(sales, 7), 0.0001)

	if listings != nil {
		competitionFactor = 1.0 / math.Max(1, float64(len(*listings)))
		// If there's any market listings for this item then see what it's currently being sold for
		if (len(*listings)) > 0 {
			for _, listing := range *listings {
				// Only return values on the info's server (as that's the only place they can sell it)
				if listing.WorldId != info.HomeServer {
					continue
				}

				listingSale := SaleMethod{
					ExchangeType:      readertype.Marketboard,           // TODO put info's world name here, change this to a more complex type
					Value:             listing.Total - listing.Quantity, // 1 gil undercut per item
					Quantity:          listing.Quantity,
					ValuePer:          listing.PricePer - 1, // 1 gil undercut
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
			if sales != nil {
				if len(*sales) > 0 {
					totalSaleValue := 0
					totalQuantity := 0

					for _, sale := range *sales {
						totalSaleValue += sale.PricePer
						totalQuantity += sale.Quantity
					}

					averageSale := totalSaleValue / totalQuantity
					averageQuantity := totalQuantity / len(*sales)

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
				// Get equivalent gil value for this currency
				gilValue, _, err := p.getGilValueAndBestSaleForCurrency(exchangeType, info)

				if err != nil {
					continue
				}

				gilCost := int(float64(exchangeMethod.GetCost()) * gilValue)

				currentMethod.ExchangeType = exchangeType
				currentMethod.Value = gilCost
				currentMethod.Quantity = exchangeMethod.GetQuantity()
				currentMethod.ValuePer = gilCost / exchangeMethod.GetQuantity()
			}

			if bestSale == nil || currentMethod.ValuePer > bestSale.ValuePer {
				bestSale = &currentMethod
			}
		}
	}

	if bestSale == nil || bestSale.Value == 0 {
		return nil
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
		return cost
	}

	// Since the list of items you're buying might over-buy, we get the Quantity from the actual required item counts
	for _, item := range o.ShoppingCart.ItemsToBuy {
		if quantity, ok := o.ShoppingCart.itemsRequired[item.GetItemId()]; ok {
			cost += item.GetCostPer() * quantity
		}
	}

	return cost
}

func (o *ObtainMethod) GetCostPerItem() int {
	return o.GetCost() / o.Quantity
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
	item *Item, numRequired int, listings *[]*db.Listing, player *PlayerInfo,
) *ObtainMethod {
	var cheapestMethod *ObtainMethod

	if item.ObtainMethods != nil {
		cheapestMethod = p.nonMarketObtainMethod(item, numRequired, cheapestMethod, player)
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
		cheapestMethod = p.craftingObtainMethod(item, numRequired, listings, cheapestMethod, player)
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
	item *Item, numRequired int, cheapestMethod *ObtainMethod, info *PlayerInfo,
) *ObtainMethod {
	for _, obtainMethod := range *item.ObtainMethods {
		obtainCost := 1500

		switch obtainMethod.GetExchangeType() {
		case readertype.GrandCompanySeal:
			if info.GrandCompanyRank < obtainMethod.(exchange.GcSealExchange).RankRequired {
				continue
			}
		case readertype.Gathering:
		case readertype.Gil:
			obtainCost = obtainMethod.GetCostPerItem()
			break
		default:
			currencyObtain, err := p.getCheapestMethodToObtainCurrency(obtainMethod.GetExchangeType(), info)

			if err == nil && currencyObtain != nil {
				obtainCost = currencyObtain.GetCostPerItem() * obtainMethod.GetQuantity()
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

	// If no sales in the specified range, return 0
	if len(filteredSales) <= 1 {
		return 0
	}

	// Calculate the total gap in hours between consecutive sales
	var totalGapHours float64 = 0
	for i := 1; i < len(filteredSales); i++ {
		gap := filteredSales[i].Timestamp.Sub(filteredSales[i-1].Timestamp).Hours()
		totalGapHours += gap
	}

	// Calculate average gap in hours (total gap divided by number of gaps)
	avgGapHours := totalGapHours / float64(len(filteredSales)-1)

	// Convert average gap time into sales per hour
	if avgGapHours == 0 {
		return 0
	}
	return 1 / avgGapHours
}

func (p *ProfitCalculator) CalculateProfitForItem(item *Item, info *PlayerInfo) (*ProfitInfo, error) {
	// Pre-calculate all items that could be involved in the obtaining of this item
	itemIds := make([]int, 0, 10)
	if item.CraftingRecipes != nil {
		itemMap := p.getPossibleSubItems(nil, item, info.SkipCrystals)
		for itemId := range itemMap {
			itemIds = append(itemIds, itemId)
		}
	} else if !item.MarketProhibited {
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
			if listing.WorldId == info.HomeServer {
				listingsOnPlayerWorld = append(listingsOnPlayerWorld, listing)
			}
		}
	}

	sales, err := p.repository.GetSalesForItemOnWorld(item.Id, info.HomeServer)
	if err != nil {
		return nil, err
	}

	// Get most value created when selling the item
	bestSale := p.GetBestSaleMethod(item, &listingsOnPlayerWorld, sales, info, false)
	if bestSale == nil {
		// Sometimes there's no way to sell this item, and that's okay. We will just return early
		return nil, nil
	}

	// Get the cheapest method to obtain the item
	cheapestMethod := p.GetCheapestObtainMethod(item, bestSale.Quantity, listings, info)
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

	adjustedProfit := float64(profitMargin) * bestSale.SaleVelocity
	profitScore := (adjustedProfit * bestSale.CompetitionFactor) / effort

	return profitScore
}

func (p *ProfitCalculator) getGilValueAndBestSaleForCurrency(currency string, info *PlayerInfo) (
	float64, *SaleMethod, error,
) {
	if p.currencyByObtainMethod == nil {
		return 0, nil, errors.New("map of currencies by currency method is nil")
	}

	itemsWithObtainMethod, ok := (*p.currencyByObtainMethod)[currency]
	if !ok {
		return 0, nil, errors.New("no currency method found for currency")
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

	listings, err := p.repository.GetListingsForItemsOnWorld(itemIds, info.HomeServer)
	if err != nil {
		return 0, nil, nil
	}

	sales, err := p.repository.GetSalesForItemsOnWorld(itemIds, info.HomeServer)
	if err != nil {
		return 0, nil, nil
	}

	wg := sync.WaitGroup{}
	type scoreAndSale struct {
		score       float64
		perCurrency float64
		item        *Item
		sale        *SaleMethod
	}
	scoreChan := make(chan scoreAndSale)

	for _, item := range itemsWithObtainMethod {
		wg.Add(1)

		go func(item *Item, listings *[]*db.Listing, sales *[]*db.Sale, info *PlayerInfo) {
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

			itemSale := p.GetBestSaleMethod(item, &filteredListings, &filteredSales, info, true)

			if itemSale == nil {
				return
			}

			cost := 0
			effortFactor := 1.0
			adjustedQuantity := 1
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
			scoreChan <- scoreAndSale{
				score:       currentScore,
				perCurrency: float64(itemSale.ValuePer) / float64(adjustedQuantity),
				item:        item,
				sale:        itemSale,
			}
		}(item, listings, sales, info)
	}

	go func() {
		wg.Wait()
		close(scoreChan)
	}()

	best := scoreAndSale{
		score:       0,
		perCurrency: 0,
		item:        nil,
	}
	for result := range scoreChan {
		if result.score > best.score {
			best = result
		}
	}

	return best.perCurrency, best.sale, nil
}

func (p *ProfitCalculator) GetBestItemToSellForCurrency(currency string, info *PlayerInfo) (*SaleMethod, error) {
	_, sale, err := p.getGilValueAndBestSaleForCurrency(currency, info)

	if err != nil {
		return nil, err
	}

	return sale, nil
}

func (p *ProfitCalculator) GetGilValueForCurrency(currency string, info *PlayerInfo) (float64, error) {
	gilEquivalent, _, err := p.getGilValueAndBestSaleForCurrency(currency, info)

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

	itemIds := make([]int, 0, len(itemsWithExchangeMethod))
	for _, item := range itemsWithExchangeMethod {
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

			itemMethod := p.GetCheapestObtainMethod(item, 1, listingResults, info)

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
