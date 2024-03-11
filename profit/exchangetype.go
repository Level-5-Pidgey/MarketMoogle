package profitCalc

import "strings"

type ExchangeType int

const (
	ExchangeTypeDefault = iota
	ExchangeTypeGil
	ExchangeTypeGcSeal
	ExchangeTypeGathering
	ExchangeTypeMarketboard
)

func FromString(s string) ExchangeType {
	switch strings.ToLower(s) {
	case "gil":
		return ExchangeTypeGil
	case "gcseal":
		return ExchangeTypeGcSeal
	case "gathering":
		return ExchangeTypeGathering
	case "marketboard":
		return ExchangeTypeMarketboard
	default:
		return ExchangeTypeDefault
	}
}

func (e ExchangeType) String() string {
	switch e {
	case ExchangeTypeGil:
		return "Buy/Sell with Gil from NPC"
	case ExchangeTypeGcSeal:
		return "Exchange Grand Company Seals"
	case ExchangeTypeGathering:
		return "Gathering"
	case ExchangeTypeMarketboard:
		return "Marketboard"
	default:
		return "Unknown"
	}
}
