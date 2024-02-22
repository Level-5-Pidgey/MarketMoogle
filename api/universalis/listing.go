package universalis

type Listing struct {
	MarketInfo     `bson:",inline"`
	ListingId      string `bson:"listingID"`
	RetainerId     string `bson:"retainerId"`
	RetainerName   string `bson:"retainerName"`
	RetainerCity   int    `bson:"retainerCity"`
	Tax            int    `bson:"tax"`
	LastReviewTime int    `bson:"lastReviewTime"`
}
