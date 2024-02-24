package domain

type DataCenter struct {
	Key    int
	Name   string
	Group  int
	IsTest bool
}

func (r DataCenter) GetKey() int {
	return r.Key
}
