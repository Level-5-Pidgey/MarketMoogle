package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/domain"
	"github.com/level-5-pidgey/MarketMoogle/util"
)

func NewRegionReader() *csv.UngroupedXivApiCsvReader[csvType.DataCenter] {
	csvColumns := map[string]int{
		"id":     0,
		"name":   1,
		"group":  2,
		"isTest": 3,
	}

	return &csv.UngroupedXivApiCsvReader[csvType.DataCenter]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.DataCenter]{
			FileName:   "WorldDCGroupType",
			RowsToSkip: 4,
			ProcessRow: func(record []string) (*csvType.DataCenter, error) {
				regionName := record[csvColumns["name"]]
				isTest := util.SafeStringToBool(record[csvColumns["isTest"]])

				if regionName == "" || isTest {
					return nil, nil
				}

				return &csvType.DataCenter{
					Key:    util.SafeStringToInt(record[csvColumns["id"]]),
					Name:   regionName,
					Group:  util.SafeStringToInt(record[csvColumns["group"]]),
					IsTest: isTest,
				}, nil
			},
		},
	}
}
