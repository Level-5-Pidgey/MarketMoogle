package db

import "github.com/level-5-pidgey/MarketMoogle/domain"

type Repository interface {
	Connect(connectionInfo string) error
	CreatePartitions() error

	CreateListing(listing domain.Listing) (*domain.Listing, error)
	CreateListings(listings *[]domain.Listing) error
	GetListingsByItemAndWorldId(itemId, worldId int) (*[]*domain.Listing, error)
	GetListingsForItemOnDataCenter(itemId, dataCenterId int) (*[]*domain.Listing, error)
	DeleteListingByUniversalisId(listingId int) error
	DeleteListings(universalisListingId []string) error

	// Market RecentHistory

	GetSalesByItemAndWorldId(itemId, worldId int) (*[]*domain.Sale, error)
	GetSalesForItemOnDataCenter(itemId, dataCenterId int) (*[]*domain.Sale, error)
	CreateSale(sale domain.Sale) (*domain.Sale, error)
	CreateSales(sales *[]domain.Sale) error
	DeleteSaleById(saleId int) error
	// DeleteSales(universalisSalesId []string) error
}
