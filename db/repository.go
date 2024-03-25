package db

type Repository interface {
	GetListingsForItemsOnWorld(itemIds []int, worldId int) (*[]Listing, error)
	GetListingsForItemsOnDataCenter(itemIds []int, dataCenterId int) (*[]Listing, error)

	// Market RecentHistory

	GetSalesForItemOnDataCenter(itemId, dataCenterId int) (*[]Sale, error)
	GetSalesForItemsOnWorld(itemIds []int, worldId int) (*[]Sale, error)
	GetSalesForItemsOnDataCenter(itemIds []int, dataCenterId int) (*[]Sale, error)
}
