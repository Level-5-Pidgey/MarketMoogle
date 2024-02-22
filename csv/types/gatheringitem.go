package csv

type GatheringItem struct {
	Key                   int
	ItemId                int
	GatheringItemLevelKey int
	IsHidden              bool
}

func (g GatheringItem) GetKey() int {
	return g.ItemId
}
