package domain

type GatheringPointBase struct {
	Key                 int
	GatheringTypeKey    int
	GatheringPointLevel int
	GatheringItemKeys   []int
}

func (g GatheringPointBase) GetKey() int {
	return g.Key
}
