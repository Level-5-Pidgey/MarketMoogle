package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type GatheringPoint struct {
	Key                  int
	GatheringTypeId      int
	GatheringPointBaseId int
	TerritoryTypeId      int
	PlaceNameId          int
}

func (g GatheringPoint) CreateFromCsvRow(record []string) (*GatheringPoint, error) {
	return &GatheringPoint{
		Key:                  util.SafeStringToInt(record[0]),
		GatheringTypeId:      util.SafeStringToInt(record[1]),
		GatheringPointBaseId: util.SafeStringToInt(record[3]),
		TerritoryTypeId:      util.SafeStringToInt(record[7]),
		PlaceNameId:          util.SafeStringToInt(record[8]),
	}, nil
}

func (g GatheringPoint) GetKey() int {
	return g.GatheringPointBaseId
}
