package main

import (
	"github.com/go-chi/chi/v5"
	dc "github.com/level-5-pidgey/MarketMoogle/csv/datacollection"
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
	profitCalc "github.com/level-5-pidgey/MarketMoogle/profit"
	"github.com/level-5-pidgey/MarketMoogle/util"
	"net/http"
	"sort"
)

type Controller struct {
	dataCollection *dc.DataCollection
	worlds         *map[int]*readertype.World
	profitCalc     *profitCalc.ProfitCalculator
}

func (c Controller) getDcIdFromWorldId(queryWorldId int) int {
	dcId := 0

	// Get datacenter id from the player's world id
	for _, world := range *c.worlds {
		if world.Id == queryWorldId {
			dcId = world.DataCenterId
		}
	}
	return dcId
}

func (c Controller) GetProfitInfo(w http.ResponseWriter, r *http.Request) {
	itemId := util.SafeStringToInt(chi.URLParam(r, "itemId"))
	queryWorldId := util.SafeStringToInt(chi.URLParam(r, "worldId"))
	dcId := c.getDcIdFromWorldId(queryWorldId)

	playerInfo := profitCalc.PlayerInfo{
		HomeServer:       queryWorldId,
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
	worldId := util.SafeStringToInt(chi.URLParam(r, "worldId"))
	dcId := c.getDcIdFromWorldId(worldId)

	result := make([]*profitCalc.ProfitInfo, 0)
	for _, item := range *c.profitCalc.ItemMap {
		if item.MarketProhibited {
			continue
		}

		playerInfo := profitCalc.PlayerInfo{
			HomeServer:       worldId,
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
