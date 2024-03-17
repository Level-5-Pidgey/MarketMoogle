package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	dc "github.com/level-5-pidgey/MarketMoogle/csv/datacollection"
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
	"github.com/level-5-pidgey/MarketMoogle/db"
	profitCalc "github.com/level-5-pidgey/MarketMoogle/profit"
	"net/http"
)

func Routes(
	items *map[int]*profitCalc.Item,
	itemsByObtainInfo *map[string]map[int]*profitCalc.Item,
	itemsByExchangeMethod *map[string]map[int]*profitCalc.Item,
	collection *dc.DataCollection,
	worlds *map[int]*readertype.World,
	db db.Repository,
) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins:   []string{"https://*", "http://*"},
				AllowedMethods:   []string{"GET"},
				AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: true,
				MaxAge:           300, // Maximum value not ignored by any of major browsers
			},
		),
	)

	// Create route controller
	controller := Controller{
		dataCollection: collection,
		worlds:         worlds,
		profitCalc:     profitCalc.NewProfitCalculator(items, itemsByObtainInfo, itemsByExchangeMethod, db),
	}

	// Item Routes
	router.Get("/api/v1/server/{worldId}/items/{itemId}/profit", controller.GetItemProfit)
	router.Get("/api/v1/server/{worldId}/items/profit", controller.GetAllItemProfit)

	// Currency
	router.Get("/api/v1/server/{worldId}/currency/{currency}/value", controller.GetGilValueOfCurrency)
	router.Get("/api/v1/server/{worldId}/currency/{currency}/best-sell", controller.GetBestItemToSellForCurrency)
	router.Get("/api/v1/server/{worldId}/currency/{currency}/min-cost", controller.GetCostToAcquireCurrency)
	router.Get("/api/v1/server/{worldId}/currency/{currency}/best-convert", controller.GetCheapestAcquisitionOfCurrency)

	return router
}
