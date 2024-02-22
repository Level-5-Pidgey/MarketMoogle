package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/level-5-pidgey/MarketMoogleApi/csv"
	dc "github.com/level-5-pidgey/MarketMoogleApi/csv/datacollection"
	"github.com/level-5-pidgey/MarketMoogleApi/db"
	profitCalc "github.com/level-5-pidgey/MarketMoogleApi/profit"
	"net/http"
	"sort"
)

type Controller struct {
	dataCollection *dc.DataCollection
	serverMap      *map[int]db.GameRegion
	profitCalc     *profitCalc.ProfitCalculator
}

func (c Controller) GetProfitInfo(w http.ResponseWriter, r *http.Request) {
	itemId := csv.SafeStringToInt(chi.URLParam(r, "itemId"))
	serverId := csv.SafeStringToInt(chi.URLParam(r, "worldId"))
	dcId := 0

	for _, region := range *c.serverMap {
		for _, dataCenter := range region.DataCenters {
			for _, world := range dataCenter.Worlds {
				if world.Id == serverId {
					dcId = dataCenter.Id
				}
			}
		}
	}

	playerInfo := profitCalc.PlayerInfo{
		HomeServer:       serverId,
		DataCenter:       dcId,
		GrandCompanyRank: 99,
	}
	itemMap := *c.profitCalc.ItemMap
	item, ok := itemMap[itemId]

	if !ok {
		ErrorJSON(w, nil, http.StatusNotFound)
		return
	}

	profitInfo, err := c.profitCalc.CalculateProfitForItem(item, &playerInfo)
	if err != nil {
		ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err = WriteJSON(w, http.StatusOK, profitInfo)
	if err != nil {
		ErrorJSON(w, err, http.StatusNotFound)
	}
}

func (c Controller) GetAllProfitInfo(w http.ResponseWriter, r *http.Request) {
	serverId := csv.SafeStringToInt(chi.URLParam(r, "worldId"))

	dcId := 0

	for _, region := range *c.serverMap {
		for _, dataCenter := range region.DataCenters {
			for _, world := range dataCenter.Worlds {
				if world.Id == serverId {
					dcId = dataCenter.Id
				}
			}
		}
	}

	result := make([]*profitCalc.ProfitInfo, 0)
	for _, item := range *c.profitCalc.ItemMap {
		if item.MarketProhibited {
			continue
		}

		playerInfo := profitCalc.PlayerInfo{
			HomeServer:       serverId,
			DataCenter:       dcId,
			GrandCompanyRank: 99,
		}

		profitInfo, err := c.profitCalc.CalculateProfitForItem(item, &playerInfo)
		if err != nil {
			ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		if profitInfo == nil {
			continue
		}

		if profitInfo.SaleMethod.ValuePer > 1000000 {
			continue
		}

		result = append(result, profitInfo)
	}

	sort.Slice(
		result, func(i, j int) bool {
			if result[i].ProfitScore == result[j].ProfitScore {
				return result[i].ObtainMethod.Cost > result[j].ObtainMethod.Cost
			}

			return result[i].ProfitScore > result[j].ProfitScore
		},
	)

	top25 := result[0:25]

	err := WriteJSON(w, http.StatusOK, top25)
	if err != nil {
		ErrorJSON(w, err, http.StatusNotFound)
	}
}
