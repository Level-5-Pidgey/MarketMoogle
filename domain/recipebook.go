package domain

type RecipeBook struct {
	Id         int
	BookItemId int
	BookName   string
}

func (r RecipeBook) GetKey() int {
	return r.Id
}
