package csv

import (
	csvEncoding "encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
)

type GroupedXivCsvReader[T readerType[T]] struct {
	GenericXivCsvReader[T]
}

func (gcr GroupedXivCsvReader[T]) readCsvData(reader *csvEncoding.Reader) (map[int][]*T, error) {
	var wg sync.WaitGroup
	extractionChannel := make(chan *T)
	results := make(map[int][]*T)

	// Append results in channel to final array
	go func() {
		for result := range extractionChannel {
			resultKey := (*result).GetKey()
			valuesInResult := results[resultKey]

			if valuesInResult == nil {
				valuesInResult = make([]*T, 0)
			}

			valuesInResult = append(valuesInResult, result)

			results[resultKey] = valuesInResult
		}
	}()

	// Read CSV
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		wg.Add(1)

		go func(record []string) {
			defer wg.Done()

			var t T
			item, err := t.CreateFromCsvRow(record)

			if item == nil || err != nil {
				return
			}

			extractionChannel <- item
		}(record)
	}

	wg.Wait()
	close(extractionChannel)

	return results, nil
}

func (gcr GroupedXivCsvReader[T]) ProcessCsv() (results interface{}, err error) {
	defer func() {
		if err == nil && results == nil {
			err = errors.New("no results found")
		}
	}()

	reader, readCloser, err := gcr.getReader()
	if err != nil {
		return nil, fmt.Errorf("unable to get reader: %w", err)
	}

	defer func(closer io.ReadCloser) {
		err := closer.Close()
		if err != nil {
			log.Fatalf("Error closing response body: %v", err)
		}
	}(readCloser)

	if err = gcr.skipHeaderRows(reader); err != nil {
		return nil, fmt.Errorf("couldn't skip header rows: %w", err)
	}

	if results, err = gcr.readCsvData(reader); err != nil {
		return nil, fmt.Errorf("couldn't read csv data: %w", err)
	}

	return results, nil
}

func (gcr GroupedXivCsvReader[T]) GetReaderType() string {
	return gcr.FileName
}
