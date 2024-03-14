package readertype

type Type string

const (
	DefaultExchangeType = ""
	Gathering           = "Gathering"
	Marketboard         = "Marketboard"
)

func (t Type) String() string {
	return string(t)
}
