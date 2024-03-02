package db

import "time"

type Sale struct {
	Id            int    `json:"sales_id"`
	ItemId        int    `json:"item_id"`
	WorldId       int    `json:"world_id"`
	PricePer      int    `json:"price_per_unit"`
	Quantity      int    `json:"quantity"`
	TotalPrice    int    `json:"total_price"`
	IsHighQuality bool   `json:"is_high_quality"`
	BuyerName     string `json:"buyer_name"`
	// Unix timestamp of sale
	Timestamp time.Time `json:"sale_time"`
}
