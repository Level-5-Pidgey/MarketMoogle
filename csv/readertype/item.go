package readertype

import (
	"errors"
	"github.com/level-5-pidgey/MarketMoogle/util"
	"strings"
)

type Item struct {
	Id                 int
	Name               string
	Description        string
	IconId             int
	ItemLevel          int
	Rarity             int
	UiCategory         int
	SearchCategory     int
	SortCategory       int
	StackSize          int
	BuyFromVendorPrice int
	SellToVendorPrice  int
	ClassJobCategory   int
	CanBeTraded        bool
	CanDesynth         bool
	DropsFromDungeon   bool
	CanBeHq            bool
	IsCollectable      bool
	IsGlamour          bool
}

func (i Item) isOldItem() bool {
	return i.Id <= 1600 &&
		(strings.HasPrefix(i.Name, "dated") ||
			strings.HasPrefix(i.Name, "pair of dated"))
}

func (i Item) CreateFromCsvRow(record []string) (*Item, error) {
	itemId := util.SafeStringToInt(record[0])
	itemName := record[10]

	if i.isOldItem() {
		return nil, errors.New("dated item, skipped")
	}

	if itemName == "" {
		return nil, errors.New("item name is empty")
	}

	canBeTraded := !util.SafeStringToBool(record[23])
	adjustedBuyPrice := 0
	if canBeTraded {
		adjustedBuyPrice = util.SafeStringToInt(record[26])
	}

	canDesynth := false
	desynthsTo := util.SafeStringToInt(record[37])
	if desynthsTo > 0 {
		canDesynth = true
	}

	result := &Item{
		Id:                 itemId,
		Name:               record[10],
		Description:        record[9],
		IconId:             util.SafeStringToInt(record[11]),
		ItemLevel:          util.SafeStringToInt(record[12]),
		Rarity:             util.SafeStringToInt(record[13]),
		UiCategory:         util.SafeStringToInt(record[16]),
		SearchCategory:     util.SafeStringToInt(record[17]),
		SortCategory:       util.SafeStringToInt(record[19]),
		StackSize:          util.SafeStringToInt(record[21]),
		BuyFromVendorPrice: adjustedBuyPrice,
		SellToVendorPrice:  util.SafeStringToInt(record[27]),
		ClassJobCategory:   util.SafeStringToInt(record[44]),
		CanBeTraded:        canBeTraded,
		DropsFromDungeon:   util.SafeStringToBool(record[25]),
		CanBeHq:            util.SafeStringToBool(record[28]),
		CanDesynth:         canDesynth,
		IsCollectable:      util.SafeStringToBool(record[39]),
		IsGlamour:          util.SafeStringToBool(record[91]),
	}

	return result, nil
}

func (i Item) GetKey() int {
	return i.Id
}
