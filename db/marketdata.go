package db

type MarketData struct {
	UnresolvedItems []int                  `json:"unresolvedItems"`
	Items           map[string]ItemDetails `json:"items"`
}

type ItemDetails struct {
	ItemId         int       `json:"itemID"`
	LastUploadTime int64     `json:"lastUploadTime"`
	Listings       []Listing `json:"listings"`
	Sales          []Sale    `json:"entries"`
	DataCenterName string    `json:"dcName"`
}

type MarketInfo struct {
	PricePer      int  `json:"pricePerUnit"`
	Quantity      int  `json:"quantity"`
	Total         int  `json:"total"`
	IsHighQuality bool `json:"hq"`
	ItemId        int  `json:"itemID"`
	WorldId       int  `json:"worldID"`
	RegionId      int  `json:"-"`
	DataCenterId  int  `json:"-"`
}
