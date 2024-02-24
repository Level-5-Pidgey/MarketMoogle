package csv

import (
	csvEncoding "encoding/csv"
	"github.com/level-5-pidgey/MarketMoogle/domain"
)

type UngroupedCsvReader[T domain.ReaderType] interface {
	GenericCsvReader

	ReadCsvData(reader *csvEncoding.Reader) (map[int]T, error)
}
