package profitCalc

import "github.com/level-5-pidgey/MarketMoogle/csv/readertype"

type PlayerInfo struct {
	HomeServer int

	DataCenter int

	GrandCompanyRank readertype.GrandCompanyRank

	JobLevels map[readertype.Job]int
}
