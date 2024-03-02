package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type DataCenter struct {
	Key    int
	Name   string
	Group  int
	IsTest bool
}

func (r DataCenter) CreateFromCsvRow(record []string) (*DataCenter, error) {
	regionName := record[1]
	isTest := util.SafeStringToBool(record[3])

	if regionName == "" || isTest {
		return nil, nil
	}

	return &DataCenter{
		Key:    util.SafeStringToInt(record[0]),
		Name:   regionName,
		Group:  util.SafeStringToInt(record[2]),
		IsTest: isTest,
	}, nil
}

func (r DataCenter) GetKey() int {
	return r.Key
}
