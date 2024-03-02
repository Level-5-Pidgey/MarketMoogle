package csv

type readerType[T any] interface {
	GetKey() int

	CreateFromCsvRow(record []string) (*T, error)
}
