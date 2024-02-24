package domain

type GatheringItemLevel struct {
	Key   int
	Level int
	Stars int
}

func (g GatheringItemLevel) GetKey() int {
	return g.Key
}
