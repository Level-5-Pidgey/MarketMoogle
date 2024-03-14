package exchange

import (
	"bytes"
	"fmt"
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
)

type GilExchange struct {
	CurrencyExchange
}

func NewGilExchange(price int, npc, shop string) GilExchange {
	return GilExchange{
		CurrencyExchange: CurrencyExchange{
			CurrencyType: readertype.Gil,
			Price:        price,
			Quantity:     1,
			Npc:          npc,
			ShopName:     shop,
		},
	}
}
func (gilExchange GilExchange) GetObtainDescription() string {
	var buffer bytes.Buffer
	buffer.WriteString("Buy ")

	if gilExchange.Npc != "" {
		buffer.WriteString(fmt.Sprintf("from %s", gilExchange.Npc))
	} else {
		buffer.WriteString("from vendor")
	}

	if gilExchange.ShopName != "" {
		buffer.WriteString(fmt.Sprintf(" (%s)", gilExchange.ShopName))
	}

	return buffer.String()
}
