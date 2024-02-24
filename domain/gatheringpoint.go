package domain

type GatheringPoint struct {
	Key                  int
	GatheringTypeId      int
	GatheringPointBaseId int
	TerritoryTypeId      int
	PlaceNameId          int
}

func (g GatheringPoint) GetKey() int {
	return g.GatheringPointBaseId
}
