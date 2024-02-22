package csv

type GilShopItem struct {
	Key    int
	ItemId int
}

func (g GilShopItem) GetKey() int {
	return g.ItemId
}
