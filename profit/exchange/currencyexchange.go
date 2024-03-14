package exchange

import (
	"bytes"
	"fmt"
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
)

type CurrencyExchange struct {
	CurrencyType readertype.Currency
	ShopName     string
	Npc          string
	Price        int
	Quantity     int
}

func (c CurrencyExchange) GetExchangeType() string {
	return c.CurrencyType.String()
}

func (c CurrencyExchange) GetObtainDescription() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Exchange %s ", c.CurrencyType.GetPlural()))

	if c.Npc != "" {
		buffer.WriteString(fmt.Sprintf("from %s", c.Npc))
	} else {
		buffer.WriteString("from NPC")
	}

	if c.ShopName != "" {
		buffer.WriteString(fmt.Sprintf(" (%s)", c.ShopName))
	}

	return buffer.String()
}

func (c CurrencyExchange) GetCost() int {
	return c.Price
}

func (c CurrencyExchange) GetQuantity() int {
	return c.Quantity
}

func (c CurrencyExchange) GetCostPerItem() int {
	return c.Price / c.Quantity
}

func (c CurrencyExchange) GetEffortFactor() float64 {
	return c.CurrencyType.GetEffort()
}
