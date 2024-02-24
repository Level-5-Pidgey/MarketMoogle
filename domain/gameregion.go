package domain

type World struct {
	Id             int
	Name           string
	RegionId       int
	RegionName     string
	DataCenterId   int
	DataCenterName string
}

func (w World) GetKey() int {
	return w.Id
}
