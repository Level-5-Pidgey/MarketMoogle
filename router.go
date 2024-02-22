package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	dc "github.com/level-5-pidgey/MarketMoogleApi/csv/datacollection"
	"github.com/level-5-pidgey/MarketMoogleApi/db"
	profitCalc "github.com/level-5-pidgey/MarketMoogleApi/profit"
	"net/http"
)

func Routes(
	items *map[int]*profitCalc.Item, collection *dc.DataCollection, servers *map[int]db.GameRegion,
	db db.Repository,
) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins:   []string{"https://*", "http://*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
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
		serverMap:      servers,
		profitCalc:     profitCalc.NewProfitCalculator(items, db),
	}

	// Routes here
	router.Get("/api/v1/server/{worldId}/items/{itemId}/profit", controller.GetProfitInfo)
	router.Get("/api/v1/server/{worldId}/items/profit", controller.GetAllProfitInfo)

	return router
}