package readertype

import (
	"github.com/level-5-pidgey/MarketMoogle/util"
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
	result := SpecialShop{
		Key:      util.SafeStringToInt(record[0]),
		ShopName: record[1],
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

	for i := 2; i < maxPerShop; i++ {
		items := make([]ItemReceived, 0, 2)
		for ii := 0; ii < 2; ii++ {
			item := getItem(record, i, ii)

			if item.ItemReceived > 1 {
				items = append(items, item)
			}
		}

		exchanges := make([]CostToBuy, 0, 3)
		for iii := 0; iii < 3; iii++ {
			exchange := getExchange(record, i, iii)

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
	return ItemReceived{
		ItemReceived: util.SafeStringToInt(record[index+(offset*maxPerShop)]),
		Quantity:     util.SafeStringToInt(record[index+((offset+1)*maxPerShop)]),
		Category:     util.SafeStringToInt(record[index+((offset+2)*maxPerShop)]),
		IsHq:         util.SafeStringToBool(record[index+((offset+3)*maxPerShop)]),
	}
}

func getExchange(record []string, index, offset int) CostToBuy {
	costItem := util.SafeStringToInt(record[index+((offset+8)*maxPerShop)])

	// The ids here for some currencies don't line up so we have to manually fix them
	if costItem < 10 {
		switch costItem {
		case 2:
			costItem = ToItemId(UncappedTomestone)
		case 3:
			costItem = ToItemId(CappedTomestone)
		case 6:
			costItem = ToItemId(PurpleCraftersScrip)
		case 7:
			costItem = ToItemId(PurpleGatherersScrip)
		}
	}

	return CostToBuy{
		CostItem:       costItem,
		Quantity:       util.SafeStringToInt(record[index+((offset+9)*maxPerShop)]),
		IsHq:           util.SafeStringToBool(record[index+((offset+10)*maxPerShop)]),
		Collectability: util.SafeStringToInt(record[index+((offset+11)*maxPerShop)]),
	}
}
