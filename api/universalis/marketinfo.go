package universalis

type MarketInfo struct {
	PricePerUnit int  `bson:"pricePerUnit"`
	Quantity     int  `bson:"quantity"`
	Total        int  `bson:"total"`
	Hq           bool `bson:"hq"`
	WorldId      int  `bson:"worldID"`
}
