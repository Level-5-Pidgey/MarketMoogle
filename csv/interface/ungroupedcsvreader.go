package csv

import csvEncoding "encoding/csv"

type UngroupedCsvReader[T ReaderType] interface {
	GenericCsvReader

	ReadCsvData(reader *csvEncoding.Reader) (map[int]T, error)
}
