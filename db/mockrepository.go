package db

type MockRepository struct {
	listings map[int]*Listing
	sales    map[int]*Sale
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		listings: make(map[int]*Listing),
		sales:    make(map[int]*Sale),
	}
}

func (r *MockRepository) CreateListing(listing Listing) (*Listing, error) {
	r.listings[listing.Id] = &listing

	return &listing, nil
}

func (r *MockRepository) CreateListings(listings *[]Listing) error {
	for _, listing := range *listings {
		r.listings[listing.Id] = &listing
	}

	return nil
}

func (r *MockRepository) GetListingsByItemAndWorldId(itemId, worldId int) (*[]*Listing, error) {
	result := make([]*Listing, 0)
	for _, listing := range r.listings {
		if listing.ItemId == itemId && listing.WorldId == worldId {
			result = append(result, listing)
		}
	}

	return &result, nil
}

func (r *MockRepository) GetListingsForItemOnDataCenter(itemId, dataCenterId int) (*[]*Listing, error) {
	result := make([]*Listing, 0)
	for _, listing := range r.listings {
		// We're just gonna pretend that there's only 1 data center within tests
		if listing.ItemId == itemId {
			result = append(result, listing)
		}
	}

	return &result, nil
}

func (r *MockRepository) DeleteListingByUniversalisId(listingUniversalisId int) error {
	if _, ok := r.listings[listingUniversalisId]; ok {
		delete(r.listings, listingUniversalisId)
		return nil
	}

	return nil
}

func (r *MockRepository) DeleteListings(universalisListingIds []string) error {
	for _, listing := range r.listings {
		for _, listingUniversalisId := range universalisListingIds {
			if listing.UniversalisId == listingUniversalisId {
				delete(r.listings, listing.Id)
			}
		}
	}

	return nil
}

// Market RecentHistory

func (r *MockRepository) GetSalesByItemAndWorldId(itemId, worldId int) (*[]*Sale, error) {
	result := make([]*Sale, 0)
	for _, sale := range r.sales {
		if sale.ItemId == itemId && sale.WorldId == worldId {
			result = append(result, sale)
		}
	}

	return &result, nil
}

func (r *MockRepository) GetSalesForItemOnDataCenter(itemId, dataCenterId int) (*[]*Sale, error) {
	result := make([]*Sale, 0)
	for _, sale := range r.sales {
		// We're just gonna pretend that there's only 1 data center within tests
		if sale.ItemId == itemId {
			result = append(result, sale)
		}
	}

	return &result, nil
}

func (r *MockRepository) CreateSale(sale Sale) (*Sale, error) {
	r.sales[sale.Id] = &sale

	return &sale, nil
}

func (r *MockRepository) CreateSales(sales *[]Sale) error {
	for _, sale := range *sales {
		r.sales[sale.Id] = &sale
	}

	return nil
}

func (r *MockRepository) DeleteSaleById(saleId int) error {
	if _, ok := r.sales[saleId]; ok {
		delete(r.sales, saleId)
		return nil
	}

	return nil
}

func (r *MockRepository) Connect(connectionInfo string) error {
	return nil
}
