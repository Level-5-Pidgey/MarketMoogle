package csv

import (
	"github.com/level-5-pidgey/MarketMoogleApi/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogleApi/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogleApi/csv/types"
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
				isTest := csv.SafeStringToBool(record[csvColumns["isTest"]])

				if regionName == "" || isTest {
					return nil, nil
				}

				return &csvType.DataCenter{
					Key:    csv.SafeStringToInt(record[csvColumns["id"]]),
					Name:   regionName,
					Group:  csv.SafeStringToInt(record[csvColumns["group"]]),
					IsTest: isTest,
				}, nil
			},
		},
	}
}
