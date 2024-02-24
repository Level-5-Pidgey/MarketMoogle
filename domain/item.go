package domain

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

func (g Item) GetKey() int {
	return g.Id
}
