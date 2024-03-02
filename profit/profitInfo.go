package profitCalc

import (
	"fmt"
	"github.com/level-5-pidgey/MarketMoogle/db"
	"log"
	"math"
	"sort"
	"time"
)

type ProfitCalculator struct {
	ItemMap *map[int]*Item
	Db      db.Repository
}

type PlayerInfo struct {
	// TODO add levels and gear to restrict options

	HomeServer int

	DataCenter int

	GrandCompanyRank int
}

const (
	maxLevel = 90
)

func NewProfitCalculator(itemMap *map[int]*Item, repo db.Repository) *ProfitCalculator {
	return &ProfitCalculator{
		ItemMap: itemMap,
		Db:      repo,
	}
}

type SaleMethod struct {
	// TODO add the currency you get for exchanging this item

	// What's the method to sell this item?
	ExchangeType string
	// How much currency you're getting from this sale
	Value int
	// How many items you need to sell to get the currency
	Quantity int
	// How much currency are you getting per required item
	ValuePer int
}

// GetBestSaleMethod
// Get the method of exchange that returns the most gil (or currency) on this item.
// Includes selling this item on the marketboard
func (p *ProfitCalculator) GetBestSaleMethod(
	item *Item, listings *[]*db.Listing, sales *[]*db.Sale, player *PlayerInfo,
) *SaleMethod {
	var bestSale *SaleMethod

	if listings != nil {
		derefListings := *listings
		// If there's any market listings for this item then see what it's currently being sold for
		if (len(derefListings)) > 0 {
			for _, listing := range derefListings {
				// Only return values on the player's server (as that's the only place they can sell it)
				if listing.WorldId != player.HomeServer {
					continue
				}

				listingSale := SaleMethod{
					ExchangeType: "Marketboard", // TODO put player's world name here
					Value:        listing.Total - listing.Quantity,
					Quantity:     listing.Quantity,
					ValuePer:     listing.PricePer - 1, // 1 gil undercut
				}

				// Players will (usually) only buy the cheapest listing, so we only update if this is the cheapest
				if bestSale == nil || listingSale.ValuePer < bestSale.ValuePer {
					bestSale = &listingSale
				}
			}
		}

		/*
			If there's no market listings for this item on the player's home world we can generate an average value
			from recent sales
		*/
		if bestSale == nil || bestSale.ValuePer == 0 {
			if sales != nil {
				saleLen := len(*sales)

				if saleLen > 0 {
					totalSaleValue := 0
					totalQuantity := 0

					for _, sale := range *sales {
						totalSaleValue += sale.PricePer
						totalQuantity += sale.Quantity
					}

					averageSale := totalSaleValue / totalQuantity
					averageQuantity := totalQuantity / saleLen

					historySale := SaleMethod{
						ExchangeType: "Marketboard", // TODO put player's world name here
						Value:        averageSale * averageQuantity,
						Quantity:     averageQuantity,
						ValuePer:     averageSale,
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
		for _, exchange := range exchangeMethods {
			currentMethod := SaleMethod{
				ExchangeType: "None",
				Value:        0,
				Quantity:     0,
				ValuePer:     0,
			}

			switch exchange.(type) {
			case GilExchange:
				currentMethod.ExchangeType = fmt.Sprintf("Sell to %s", exchange.(GilExchange).NpcName)
				currentMethod.Value = exchange.GetCost()
				currentMethod.Quantity = exchange.GetQuantity()
				currentMethod.ValuePer = exchange.GetCost() / exchange.GetQuantity()
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

type ObtainInfo struct {
	ItemsRequired []*PurchaseInfo

	// TODO expand this into an object (with a type enum and human readable value)
	ObtainMethod string

	Cost int

	CostPerItem int

	ResultQuantity int

	EffortFactor float64
}

type PurchaseInfo struct {
	ItemId int

	Quantity int

	Server int

	BuyFrom string
}

func isEasierToObtain(curr, new *ObtainInfo) bool {
	if curr == nil {
		return true
	}

	currEffortCost := float64(curr.Cost) * curr.EffortFactor
	newEffortCost := float64(new.Cost) * new.EffortFactor

	if currEffortCost == newEffortCost {
		return len(new.ItemsRequired) < len(curr.ItemsRequired)
	}

	return newEffortCost < currEffortCost
}

func (p *ProfitCalculator) GetCostToObtain(
	item *Item, numRequired int, listings *[]*db.Listing, player *PlayerInfo,
) *ObtainInfo {
	var cheapestMethod *ObtainInfo

	if item.ObtainMethods != nil {
		cheapestMethod = getNonMarketObtainMethod(item, numRequired, cheapestMethod, player)
	}

	if !item.MarketProhibited && listings != nil {
		var filteredListings []*db.Listing
		for _, listing := range *listings {
			if listing.ItemId == item.Id {
				filteredListings = append(filteredListings, listing)
			}
		}

		if len(filteredListings) != 0 {
			cheapestMethod = getMarketObtainMethod(item, cheapestMethod, numRequired, &filteredListings, player)
		}
	}

	if item.CraftingRecipes != nil {
		cheapestMethod = p.getCraftingObtainMethod(item, numRequired, cheapestMethod, player)
	}

	return cheapestMethod
}

func (p *ProfitCalculator) getIngredientsForRecipe(
	itemsAndQuantities *map[int]int, numRequired int, recipe *RecipeInfo, skipCrystals bool,
) *map[int]int {
	if itemsAndQuantities == nil {
		newMap := make(map[int]int)
		itemsAndQuantities = &newMap
	}

	for _, ingredient := range recipe.RecipeIngredients {
		ingredientItem, ok := (*p.ItemMap)[ingredient.ItemId]

		if !ok {
			continue
		}

		// TODO remove this magic number, get correct value dynamically from csv load
		if ingredientItem.UiCategory == 59 && skipCrystals {
			continue
		}

		if ingredientItem.CraftingRecipes != nil {
			for _, craftingRecipe := range *ingredientItem.CraftingRecipes {
				// For each crafting recipe the ingredient has
				// Calculate the number of times it takes to get this item
				numIngredientRequired := numRequired * craftingRecipe.Yield

				itemsAndQuantities = p.getIngredientsForRecipe(
					itemsAndQuantities,
					numIngredientRequired,
					&craftingRecipe,
					skipCrystals,
				)
			}
		} else {
			// There is no sub recipe, we should be updating the itemsAndQuantities map
			if _, ok := (*itemsAndQuantities)[ingredient.ItemId]; ok {
				(*itemsAndQuantities)[ingredient.ItemId] += (ingredient.Quantity * numRequired) / recipe.Yield
			} else {
				(*itemsAndQuantities)[ingredient.ItemId] = (ingredient.Quantity * numRequired) / recipe.Yield
			}
		}
	}

	return itemsAndQuantities
}

func (p *ProfitCalculator) getCraftingObtainMethod(
	item *Item, numRequired int, cheapestMethod *ObtainInfo, player *PlayerInfo,
) *ObtainInfo {
	for _, craftingRecipe := range *item.CraftingRecipes {
		// Get all the possible ingredients needed for this recipe, recursively
		itemsRequiredForRecipe := p.getIngredientsForRecipe(
			nil,
			numRequired/craftingRecipe.Yield,
			&craftingRecipe,
			true,
		)

		recipeCost := ObtainInfo{
			ItemsRequired:  []*PurchaseInfo{},
			ResultQuantity: craftingRecipe.Yield,
			EffortFactor:   recipeEffort(&craftingRecipe, item),
			ObtainMethod:   craftingRecipe.CraftType, // TODO expand with type
		}

		canCraft := true
		for itemId, quantity := range *itemsRequiredForRecipe {
			ingredientItem, ok := (*p.ItemMap)[itemId]

			if !ok {
				continue
			}

			var ingredientListings *[]*db.Listing = nil
			if !ingredientItem.MarketProhibited {
				ingredientResults, err := p.Db.GetListingsForItemOnDataCenter(ingredientItem.Id, player.DataCenter)
				ingredientListings = ingredientResults
				if err != nil {
					log.Printf(
						"Error getting listings for item %d on data center %d: %s",
						ingredientItem.Id, player.DataCenter, err,
					)
				}
			}

			ingredientCost := p.GetCostToObtain(
				ingredientItem,
				quantity,
				ingredientListings,
				player,
			)

			if ingredientCost == nil {
				canCraft = false
				break
			}

			recipeCost.ItemsRequired = append(
				recipeCost.ItemsRequired,
				ingredientCost.ItemsRequired...,
			)
			recipeCost.Cost += ingredientCost.Cost

			if !isEasierToObtain(cheapestMethod, &recipeCost) {
				canCraft = false
				break
			}
		}

		recipeCost.CostPerItem = recipeCost.Cost / craftingRecipe.Yield

		if isEasierToObtain(cheapestMethod, &recipeCost) && canCraft {
			cheapestMethod = &recipeCost
		}
	}

	return cheapestMethod
}

func getMarketObtainMethod(
	item *Item, cheapestMethod *ObtainInfo, numRequired int, listings *[]*db.Listing, player *PlayerInfo,
) *ObtainInfo {
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

	purchasePlan := ObtainInfo{
		ItemsRequired:  []*PurchaseInfo{},
		ObtainMethod:   "Market",
		Cost:           0,
		CostPerItem:    0,
		ResultQuantity: 0,
		EffortFactor:   0.99,
	}

	for _, listing := range *listings {
		if numRequired <= 0 {
			break
		}

		numRequired -= listing.Quantity
		purchasePlan.ItemsRequired = append(
			purchasePlan.ItemsRequired, &PurchaseInfo{
				ItemId:   item.Id,
				Quantity: listing.Quantity,
				Server:   listing.WorldId,
				BuyFrom:  listing.RetainerName,
			},
		)

		if player.HomeServer == listing.WorldId {
			purchasePlan.EffortFactor += 0.01
		} else {
			purchasePlan.EffortFactor += 0.06
		}
		purchasePlan.Cost += listing.Total
		purchasePlan.ResultQuantity += listing.Quantity
		purchasePlan.CostPerItem = purchasePlan.Cost / purchasePlan.ResultQuantity
	}

	if isEasierToObtain(cheapestMethod, &purchasePlan) {
		cheapestMethod = &purchasePlan
	}

	return cheapestMethod
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

func getNonMarketObtainMethod(
	item *Item, numRequired int, cheapestMethod *ObtainInfo, player *PlayerInfo,
) *ObtainInfo {
	for _, obtainMethod := range *item.ObtainMethods {
		// TODO use a type switch and calculate equivalent gil costs for currency exchanges
		switch obtainMethod.(type) {
		case GcSealExchange:
			if player.GrandCompanyRank < obtainMethod.(GcSealExchange).RankRequired {
				continue
			}
		}

		numOfExchanges := 1
		totalEffort := obtainMethod.GetEffortFactor()
		for (numOfExchanges * obtainMethod.GetQuantity()) < numRequired {
			numOfExchanges++
			totalEffort *= 1.01
		}
		totalCost := numOfExchanges * obtainMethod.GetCost()
		totalQuantity := numOfExchanges * obtainMethod.GetQuantity()

		currentMethod := ObtainInfo{
			ItemsRequired: []*PurchaseInfo{
				{
					ItemId:   item.Id,
					Quantity: numOfExchanges,
					Server:   player.HomeServer,
				},
			},
			Cost:           totalCost,
			CostPerItem:    totalCost / totalQuantity,
			ResultQuantity: totalQuantity,
			EffortFactor:   totalEffort,
			ObtainMethod:   obtainMethod.GetObtainType(),
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
	ObtainMethod *ObtainInfo
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
	// Get market listings for item if this item is sellable
	var listings *[]*db.Listing = nil
	var listingsOnPlayerWorld []*db.Listing
	if !item.MarketProhibited {
		ingredientResults, err := p.Db.GetListingsForItemOnDataCenter(item.Id, info.DataCenter)
		listings = ingredientResults
		if err != nil {
			return nil, err
		}

		for _, listing := range *listings {
			if listing.WorldId == info.HomeServer {
				listingsOnPlayerWorld = append(listingsOnPlayerWorld, listing)
			}
		}
	}

	sales, err := p.Db.GetSalesByItemAndWorldId(item.Id, info.HomeServer)
	if err != nil {
		return nil, err
	}

	// Get most value created when selling the item
	bestSale := p.GetBestSaleMethod(item, listings, sales, info)
	if bestSale == nil {
		// Sometimes there's no way to sell this item, and that's okay. We will just return early
		return nil, nil
	}

	// Get cheapest method to obtain the item
	cheapestMethod := p.GetCostToObtain(item, bestSale.Quantity, listings, info)
	if cheapestMethod == nil {
		return nil, nil
	}

	// Calculate other variables and a "profit score"
	profitMargin := bestSale.Value - cheapestMethod.Cost

	salesPerHour := math.Max(p.salesPerHour(sales, 7), 0.0001)
	adjustedProfit := float64(profitMargin) * salesPerHour
	competitionFactor := 1.0 / math.Max(1, float64(len(listingsOnPlayerWorld)))
	profitScore := math.Round((adjustedProfit * competitionFactor) / cheapestMethod.EffortFactor)

	// Return info
	return &ProfitInfo{
		ItemId:       item.Id,
		ObtainMethod: cheapestMethod,
		SaleMethod:   bestSale,
		ProfitScore:  profitScore,
	}, nil
}
