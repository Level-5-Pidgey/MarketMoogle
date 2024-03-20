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

func (r *MockRepository) GetListingsForItemsOnWorld(itemIds []int, worldId int) (*[]*Listing, error) {
	result := make([]*Listing, 0)
	for _, itemId := range itemIds {
		for _, listing := range r.listings {
			if listing.ItemId == itemId && listing.WorldId == worldId {
				result = append(result, listing)
			}
		}
	}

	return &result, nil
}

func (r *MockRepository) GetListingsForItemsOnDataCenter(itemIds []int, dataCenterId int) (*[]*Listing, error) {
	result := make([]*Listing, 0)
	for _, itemId := range itemIds {
		for _, listing := range r.listings {
			// We're just gonna pretend that there's only 1 data center within tests
			if listing.ItemId == itemId {
				result = append(result, listing)
			}
		}
	}

	return &result, nil
}

// Market RecentHistory

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

func (r *MockRepository) GetSalesForItemsOnWorld(itemIds []int, worldId int) (*[]*Sale, error) {
	result := make([]*Sale, 0)
	for _, itemId := range itemIds {
		for _, sale := range r.sales {
			if sale.ItemId == itemId && sale.WorldId == worldId {
				result = append(result, sale)
			}
		}
	}

	return &result, nil
}

func (r *MockRepository) GetSalesForItemsOnDataCenter(itemIds []int, dataCenterId int) (*[]*Sale, error) {
	result := make([]*Sale, 0)
	for _, itemId := range itemIds {
		for _, sale := range r.sales {
			// We're just gonna pretend that there's only 1 data center within tests
			if sale.ItemId == itemId {
				result = append(result, sale)
			}
		}
	}

	return &result, nil
}
