package main

import (
	"github.com/go-chi/chi/v5"
	dc "github.com/level-5-pidgey/MarketMoogle/csv/datacollection"
	"github.com/level-5-pidgey/MarketMoogle/domain"
	profitCalc "github.com/level-5-pidgey/MarketMoogle/profit"
	"github.com/level-5-pidgey/MarketMoogle/util"
	"net/http"
	"sort"
)

type Controller struct {
	dataCollection *dc.DataCollection
	serverMap      *map[int]domain.GameRegion
	profitCalc     *profitCalc.ProfitCalculator
}

func (c Controller) GetProfitInfo(w http.ResponseWriter, r *http.Request) {
	itemId := util.SafeStringToInt(chi.URLParam(r, "itemId"))
	serverId := util.SafeStringToInt(chi.URLParam(r, "worldId"))
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
		util.ErrorJSON(w, nil, http.StatusNotFound)
		return
	}

	profitInfo, err := c.profitCalc.CalculateProfitForItem(item, &playerInfo)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err = util.WriteJSON(w, http.StatusOK, profitInfo)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusNotFound)
	}
}

func (c Controller) GetAllProfitInfo(w http.ResponseWriter, r *http.Request) {
	serverId := util.SafeStringToInt(chi.URLParam(r, "worldId"))

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
			util.ErrorJSON(w, err, http.StatusInternalServerError)
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

	err := util.WriteJSON(w, http.StatusOK, top25)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusNotFound)
	}
}
