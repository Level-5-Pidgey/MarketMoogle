package db

type Listing struct {
	MarketInfo     `json:",inline"`
	ListingId      string `json:"listingID"`
	RetainerId     string `json:"retainerId"`
	RetainerName   string `json:"retainerName"`
	RetainerCity   int    `json:"retainerCity"`
	Tax            int    `json:"tax"`
	LastReviewTime int64  `json:"lastReviewTime"`
}
