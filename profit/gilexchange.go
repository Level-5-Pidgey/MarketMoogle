package profitCalc

type GilExchange struct {
	TokenExchange
	NpcName string
}

func (gilExchange GilExchange) GetObtainType() string {
	return "Buy with Gil"
}

func (gilExchange GilExchange) GetCost() int {
	return gilExchange.Value
}

func (gilExchange GilExchange) GetQuantity() int {
	return gilExchange.Quantity
}

func (gilExchange GilExchange) GetCostPerItem() int {
	return gilExchange.GetCost() / gilExchange.GetQuantity()
}

func (gilExchange GilExchange) GetEffortFactor() float64 {
	return 0.85
}
