package domain

type ItemSearchCategory struct {
	Key           int
	Name          string
	IconId        int
	CategoryValue int
	ClassJobId    int
}

func (i ItemSearchCategory) GetKey() int {
	return i.Key
}
