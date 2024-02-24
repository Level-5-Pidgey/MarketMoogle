package domain

type TerritoryType struct {
	Key      int
	RegionId int
	PlaceId  int
	MapId    int
}

func (t TerritoryType) GetKey() int {
	return t.Key
}
