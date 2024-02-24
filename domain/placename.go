package domain

type PlaceName struct {
	Key  int
	Name string
}

func (p PlaceName) GetKey() int {
	return p.Key
}
