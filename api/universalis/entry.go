package universalis

import (
	"github.com/level-5-pidgey/MarketMoogleApi/db"
	"time"
)

type Entry struct {
	Event         string    `bson:"event"`
	Item          int       `bson:"item"`
	World         int       `bson:"world"`
	Listings      []Listing `bson:"listings"`
	Sales         []Sale    `bson:"sales"`
	RecentHistory []Sale    `bson:"recentHistory"` // Because for some godforsaken reason universalis uses a different term for the s
}

func (entry *Entry) getSaleHistory() []Sale {
	if len(entry.RecentHistory) == 0 {
		return entry.Sales
	}

	return entry.RecentHistory
}

func (entry *Entry) ConvertToDbListings() *[]db.Listing {
	result := make([]db.Listing, len(entry.Listings))
	now := time.Now()

	for listingIndex, listing := range entry.Listings {
		listingWorld := listing.WorldId
		if entry.World != 0 {
			listingWorld = entry.World
		}

		result[listingIndex] = db.Listing{
			UniversalisId: listing.ListingId,
			ItemId:        entry.Item,
			WorldId:       listingWorld,
			PricePer:      listing.PricePerUnit,
			Quantity:      listing.Quantity,
			Total:         listing.Total,
			IsHighQuality: listing.Hq,
			RetainerName:  listing.RetainerName,
			RetainerCity:  listing.RetainerCity,
			LastReview:    now,
		}
	}

	return &result
}

func (entry *Entry) ConvertToDbSales() *[]db.Sale {
	sales := entry.getSaleHistory()
	result := make([]db.Sale, len(sales))

	for saleIndex, sale := range sales {
		saleWorld := sale.WorldId
		if entry.World != 0 {
			saleWorld = entry.World
		}

		result[saleIndex] = db.Sale{
			ItemId:        entry.Item,
			WorldId:       saleWorld,
			PricePer:      sale.PricePerUnit,
			Quantity:      sale.Quantity,
			TotalPrice:    sale.Total,
			IsHighQuality: sale.Hq,
			BuyerName:     sale.BuyerName,
			Timestamp:     time.Unix(sale.Timestamp, 0),
		}
	}

	return &result
}
