package db

type Repository interface {
	Connect(connectionInfo string) error
	CreatePartitions() error

	CreateListing(listing Listing) (*Listing, error)
	CreateListings(listings *[]Listing) error
	GetListingsByItemAndWorldId(itemId, worldId int) (*[]*Listing, error)
	GetListingsForItemOnDataCenter(itemId, dataCenterId int) (*[]*Listing, error)
	DeleteListingByUniversalisId(listingId int) error
	DeleteListings(universalisListingId []string) error

	// Market RecentHistory

	GetSalesByItemAndWorldId(itemId, worldId int) (*[]*Sale, error)
	GetSalesForItemOnDataCenter(itemId, dataCenterId int) (*[]*Sale, error)
	CreateSale(sale Sale) (*Sale, error)
	CreateSales(sales *[]Sale) error
	DeleteSaleById(saleId int) error
	// DeleteSales(universalisSalesId []string) error
}
