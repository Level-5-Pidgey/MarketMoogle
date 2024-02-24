package domain

type GatheringType struct {
	Key    int
	Name   string
	IconId int
}

func (g GatheringType) GetKey() int {
	return g.Key
}
