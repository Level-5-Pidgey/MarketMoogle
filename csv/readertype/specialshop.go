package readertype

import (
	"github.com/level-5-pidgey/MarketMoogle/util"
	"log"
	"strings"
)

const (
	maxPerShop = 60
)

type SpecialShop struct {
	Key      int
	ShopName string
	Windows  []ShopWindow
}

type ShopWindow struct {
	Items       []ItemReceived
	Exchange    []CostToBuy
	Quest       int
	Achievement int
	PatchNumber int
}
type ItemReceived struct {
	ItemReceived int
	Quantity     int
	Category     int
	IsHq         bool
}

type CostToBuy struct {
	CostItem       int
	Quantity       int
	IsHq           bool
	Collectability int
}

func (s SpecialShop) GetKey() int {
	return s.Key
}

func (s SpecialShop) CreateFromCsvRow(record []string) (*SpecialShop, error) {
	shopName := record[1]
	result := SpecialShop{
		Key:      util.SafeStringToInt(record[0]),
		ShopName: shopName,
		Windows:  make([]ShopWindow, 0, maxPerShop),
	}

	/*
		The structure for a special shop is as follows:
			- Item Received [0]
			- Quantity Received [0]
			- Category [0]
			- Is High Quality [0]
			- Item Received [1]
			- Quantity Received [1]
			- Category [1]
			- Is High Quality [1]
			- Item to Buy [0]
			- Quantity of Item [0]
			- Cost item is HQ [0]
			- Collectability of Cost Item [0]
			- Item to Buy [1]
			- Quantity of Item [1]
			- Cost item is HQ [1]
			- Collectability of Cost Item [1]
			- Quest [0]
			- Achievement [0]
			- Patch Number [0]
	*/

	for i := 0; i < maxPerShop; i++ {
		index := i + 2
		items := make([]ItemReceived, 0, 2)
		for ii := 0; ii < 2; ii++ {
			item := getItem(record, index, ii)

			if item.ItemReceived > 1 {
				items = append(items, item)
			}
		}

		exchanges := make([]CostToBuy, 0, 3)
		for iii := 0; iii < 3; iii++ {
			exchange := getExchange(record, shopName, index, iii)

			if exchange.CostItem > 1 && exchange.Quantity != 0 {
				exchanges = append(exchanges, exchange)
			}
		}

		if len(items) == 0 || len(exchanges) == 0 {
			continue
		}

		trade := ShopWindow{
			Items:       items,
			Exchange:    exchanges,
			Quest:       util.SafeStringToInt(record[i+(31*maxPerShop)]),
			Achievement: util.SafeStringToInt(record[i+(32*maxPerShop)]),
			PatchNumber: util.SafeStringToInt(record[i+(33*maxPerShop)]),
		}

		result.Windows = append(result.Windows, trade)
	}

	if len(result.Windows) == 0 {
		return nil, nil
	}

	return &result, nil
}

func getItem(record []string, index, offset int) ItemReceived {
	itemOffset := offset * 4

	return ItemReceived{
		ItemReceived: util.SafeStringToInt(record[index+(itemOffset*maxPerShop)]),
		Quantity:     util.SafeStringToInt(record[index+((itemOffset+1)*maxPerShop)]),
		Category:     util.SafeStringToInt(record[index+((itemOffset+2)*maxPerShop)]),
		IsHq:         util.SafeStringToBool(record[index+((itemOffset+3)*maxPerShop)]),
	}
}

func isCrafterGathererShop(shopName string) bool {
	shopNameLower := strings.ToLower(shopName)
	possibleNames := []string{
		"scrip exchange",
		"bait/",
		"master recipe",
		"materia",
		"folklore",
		"landsaint",
		"handsaint",
		"professional",
		"materials",
		"miscellany",
		"(doh)",
		"(dol)",
	}

	for _, name := range possibleNames {
		if strings.Contains(shopNameLower, name) {
			return true
		}
	}

	return false
}

func convertItem(shopName string, costItem *int) {
	result := *costItem

	if isCrafterGathererShop(shopName) {
		switch result {
		case 1:
			result = ToItemId(PoeticTomestone)
			break
		case 2:
			result = ToItemId(WhiteCraftersScrip)
			break
		case 4:
			result = ToItemId(WhiteGatherersScrip)
			break
		case 6:
			result = ToItemId(PurpleCraftersScrip)
			break
		case 7:
			result = ToItemId(PurpleGatherersScrip)
			break
		default:
			log.Printf("Unknown gatherer/crafter cost item %d for %s", result, shopName)
		}
	} else {
		switch result {
		case 1:
			result = ToItemId(PoeticTomestone)
			break
		case 2:
			result = ToItemId(UncappedTomestone)
			break
		case 3:
			result = ToItemId(CappedTomestone)
			break
		default:
			log.Printf("Unknown combat cost item %d for %s", result, shopName)
		}
	}

	*costItem = result
}

func getExchange(record []string, shopName string, index, offset int) CostToBuy {
	exchangeOffset := offset * 4
	costItem := util.SafeStringToInt(record[index+((exchangeOffset+8)*maxPerShop)])
	quantity := util.SafeStringToInt(record[index+((exchangeOffset+9)*maxPerShop)])

	// The ids here for some currencies don't line up so we have to manually fix them
	if costItem != 0 && costItem < 10 && quantity > 0 {
		convertItem(shopName, &costItem)
	}

	return CostToBuy{
		CostItem:       costItem,
		Quantity:       util.SafeStringToInt(record[index+((exchangeOffset+9)*maxPerShop)]),
		IsHq:           util.SafeStringToBool(record[index+((exchangeOffset+10)*maxPerShop)]),
		Collectability: util.SafeStringToInt(record[index+((exchangeOffset+11)*maxPerShop)]),
	}
}
