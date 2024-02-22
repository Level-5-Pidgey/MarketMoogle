package csv

type GcScripShopItem struct {
	Key                      int
	ItemId                   int
	GrandCompanyRankRequired int
	AmountRequired           int
}

func (g GcScripShopItem) GetKey() int {
	return g.ItemId
}
