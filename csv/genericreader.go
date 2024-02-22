package csv

type GenericCsvReader struct {
	FileName   string
	RowsToSkip int
	ProcessRow func([]string) (*interface{}, error)
}
