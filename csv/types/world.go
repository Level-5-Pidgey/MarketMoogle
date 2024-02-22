package csv

type World struct {
	Key        int
	Name       string
	Region     int
	DataCenter int
	IsPublic   bool
}

func (w World) GetKey() int {
	return w.Key
}
