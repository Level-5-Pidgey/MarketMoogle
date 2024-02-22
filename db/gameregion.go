package db

type GameWorld struct {
	Id   int
	Name string
}

type DataCenter struct {
	Id     int
	Name   string
	Worlds map[int]GameWorld
}

type GameRegion struct {
	Id          int
	DataCenters map[int]DataCenter
}
