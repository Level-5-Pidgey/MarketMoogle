package domain

import "time"

type Listing struct {
	Id            int       `json:"listing_id"`
	UniversalisId string    `json:"universalis_listing_id"`
	ItemId        int       `json:"item_id"`
	WorldId       int       `json:"world_id"`
	RegionId      int       `json:"region_id"`
	DataCenterId  int       `json:"data_center_id"`
	PricePer      int       `json:"price_per_unit"`
	Quantity      int       `json:"quantity"`
	Total         int       `json:"total_price"`
	IsHighQuality bool      `json:"is_high_quality"`
	RetainerName  string    `json:"retainer_name"`
	RetainerCity  int       `json:"retainer_city"`
	LastReview    time.Time `json:"last_review_time"`
}
