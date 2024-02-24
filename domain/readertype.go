package domain

type ReaderType interface {
	GetKey() int

	CreateFromCsvRow(record []string) (ReaderType, error)
}
