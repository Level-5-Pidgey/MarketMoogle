package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type World struct {
	Id             int
	Name           string
	RegionId       int
	RegionName     string
	DataCenterId   int
	DataCenterName string
}

func (w World) CreateFromCsvRow(record []string) (*World, error) {
	isPublic := util.SafeStringToBool(record[6])

	if !isPublic {
		return nil, nil
	}

	return &World{
		Id:           util.SafeStringToInt(record[0]),
		Name:         record[2],
		RegionId:     util.SafeStringToInt(record[3]),
		DataCenterId: util.SafeStringToInt(record[5]),
	}, nil
}

func (w World) GetKey() int {
	return w.Id
}
