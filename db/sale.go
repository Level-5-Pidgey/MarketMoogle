package db

type Sale struct {
	MarketInfo `json:",inline"`
	BuyerName  string `json:"buyerName"`
	Timestamp  int64  `json:"timestamp"`
}
