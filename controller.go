package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	dc "github.com/level-5-pidgey/MarketMoogle/csv/datacollection"
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
	profitCalc "github.com/level-5-pidgey/MarketMoogle/profit"
	"github.com/level-5-pidgey/MarketMoogle/util"
	"net/http"
	"sort"
	"sync"
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

func (c Controller) GetItemProfit(w http.ResponseWriter, r *http.Request) {
	itemId := util.SafeStringToInt(chi.URLParam(r, "itemId"))
	queryWorldId := util.SafeStringToInt(chi.URLParam(r, "worldId"))
	dcId := c.getDcIdFromWorldId(queryWorldId)

	playerInfo := profitCalc.PlayerInfo{
		HomeServer:       queryWorldId,
		DataCenter:       dcId,
		GrandCompanyRank: readertype.Captain,
		JobLevels: map[readertype.Job]int{
			readertype.JobCarpenter:     90,
			readertype.JobBlacksmith:    90,
			readertype.JobArmourer:      90,
			readertype.JobGoldsmith:     90,
			readertype.JobLeatherworker: 90,
			readertype.JobWeaver:        90,
			readertype.JobAlchemist:     90,
			readertype.JobCulinarian:    90,
			readertype.JobMiner:         90,
			readertype.JobBotanist:      90,
			readertype.JobFisher:        90,
			readertype.JobPaladin:       90,
		},
	}
	itemMap := *c.profitCalc.Items
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

func (c Controller) GetAllItemProfit(w http.ResponseWriter, r *http.Request) {
	worldId := util.SafeStringToInt(chi.URLParam(r, "worldId"))
	dcId := c.getDcIdFromWorldId(worldId)

	playerInfo := profitCalc.PlayerInfo{
		HomeServer:       worldId,
		DataCenter:       dcId,
		GrandCompanyRank: readertype.Captain,
		JobLevels: map[readertype.Job]int{
			readertype.JobCarpenter:     90,
			readertype.JobBlacksmith:    90,
			readertype.JobArmourer:      90,
			readertype.JobGoldsmith:     90,
			readertype.JobLeatherworker: 90,
			readertype.JobWeaver:        90,
			readertype.JobAlchemist:     90,
			readertype.JobCulinarian:    90,
			readertype.JobMiner:         90,
			readertype.JobBotanist:      90,
			readertype.JobFisher:        90,
			readertype.JobPaladin:       90,
		},
	}

	var wg sync.WaitGroup
	resultsChan := make(chan *profitCalc.ProfitInfo)
	errorsChan := make(chan error)

	for _, item := range *c.profitCalc.Items {
		if item.MarketProhibited {
			continue
		}

		wg.Add(1)

		go func(item *profitCalc.Item, playerInfo *profitCalc.PlayerInfo) {
			defer wg.Done()

			profitInfo, err := c.profitCalc.CalculateProfitForItem(item, playerInfo)

			if err != nil {
				errorsChan <- err
				return
			} else {
				if profitInfo == nil || profitInfo.SaleMethod.ValuePer > 1000000 {
					return
				}

				resultsChan <- profitInfo
			}
		}(item, &playerInfo)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
		close(errorsChan)
	}()

	result := make([]*profitCalc.ProfitInfo, 0)
	errors := make([]error, 0)

	for {
		select {
		case resultInfo, ok := <-resultsChan:
			if !ok {
				resultsChan = nil
			} else {
				result = append(result, resultInfo)
			}

		case err, ok := <-errorsChan:
			if !ok {
				errorsChan = nil
			} else {
				errors = append(errors, err)
			}
		}

		if resultsChan == nil && errorsChan == nil {
			break
		}
	}

	if len(errors) > 0 {
		fmt.Printf("Multiple (%d) errors occurred: ", len(errors))
		for index, err := range errors {
			fmt.Printf("Error #%d: %v\n", index+1, err)
		}

		util.ErrorJSON(
			w,
			fmt.Errorf("multiple (%d) errors occurred", len(errors)),
			http.StatusInternalServerError,
		)

		return
	}

	sort.Slice(
		result, func(i, j int) bool {
			if result[i].ProfitScore == result[j].ProfitScore {
				return result[i].ObtainMethod.GetCost() > result[j].ObtainMethod.GetCost()
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

func (c Controller) GetCurrencyProfit(w http.ResponseWriter, r *http.Request) {
	serverId := util.SafeStringToInt(chi.URLParam(r, "worldId"))
	currency := chi.URLParam(r, "currency")
	dcId := c.getDcIdFromWorldId(serverId)

	playerInfo := profitCalc.PlayerInfo{
		HomeServer:       serverId,
		DataCenter:       dcId,
		GrandCompanyRank: readertype.Captain,
		JobLevels: map[readertype.Job]int{
			readertype.JobCarpenter:     90,
			readertype.JobBlacksmith:    90,
			readertype.JobArmourer:      90,
			readertype.JobGoldsmith:     90,
			readertype.JobLeatherworker: 90,
			readertype.JobWeaver:        90,
			readertype.JobAlchemist:     90,
			readertype.JobCulinarian:    90,
			readertype.JobMiner:         90,
			readertype.JobBotanist:      90,
			readertype.JobFisher:        90,
			readertype.JobPaladin:       90,
		},
	}

	exchangeType := profitCalc.FromString(currency)
	value, err := c.profitCalc.GetSellValueForCurrency(exchangeType, &playerInfo)

	if err != nil {
		util.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	err = util.WriteJSON(w, http.StatusOK, value)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusNotFound)
	}
}

func (c Controller) GetCostToAcquireCurrency(w http.ResponseWriter, r *http.Request) {
	serverId := util.SafeStringToInt(chi.URLParam(r, "worldId"))
	currency := chi.URLParam(r, "currency")
	dcId := c.getDcIdFromWorldId(serverId)

	playerInfo := profitCalc.PlayerInfo{
		HomeServer:       serverId,
		DataCenter:       dcId,
		GrandCompanyRank: readertype.Captain,
		JobLevels: map[readertype.Job]int{
			readertype.JobCarpenter:     90,
			readertype.JobBlacksmith:    90,
			readertype.JobArmourer:      90,
			readertype.JobGoldsmith:     90,
			readertype.JobLeatherworker: 90,
			readertype.JobWeaver:        90,
			readertype.JobAlchemist:     90,
			readertype.JobCulinarian:    90,
			readertype.JobMiner:         90,
			readertype.JobBotanist:      90,
			readertype.JobFisher:        90,
			readertype.JobPaladin:       90,
		},
	}

	exchangeType := profitCalc.FromString(currency)
	obtainMethod, err := c.profitCalc.GetCheapestWayToObtainCurrency(exchangeType, &playerInfo)

	if err != nil {
		util.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	err = util.WriteJSON(w, http.StatusOK, obtainMethod)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusNotFound)
	}
}
