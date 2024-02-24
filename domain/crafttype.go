package domain

type CraftType struct {
	Key  int
	Name string
}

func (c CraftType) GetKey() int {
	return c.Key
}
