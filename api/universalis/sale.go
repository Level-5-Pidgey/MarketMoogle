package universalis

type Sale struct {
	MarketInfo `bson:",inline"`
	BuyerName  string `bson:"buyerName"`
	Timestamp  int64  `bson:"timestamp"`
}
