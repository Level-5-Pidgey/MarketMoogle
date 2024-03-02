package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type TerritoryType struct {
	Key      int
	RegionId int
	PlaceId  int
	MapId    int
}

func (t TerritoryType) CreateFromCsvRow(record []string) (*TerritoryType, error) {
	return &TerritoryType{
		Key:      util.SafeStringToInt(record[0]),
		RegionId: util.SafeStringToInt(record[4]),
		PlaceId:  util.SafeStringToInt(record[5]),
		MapId:    util.SafeStringToInt(record[7]),
	}, nil
}

func (t TerritoryType) GetKey() int {
	return t.Key
}
