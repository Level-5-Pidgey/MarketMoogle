package csv

import (
	csvEncoding "encoding/csv"
	"errors"
	"fmt"
	csv "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	"io"
	"log"
	"sync"
)

type GroupedXivApiCsvReader[T csv.ReaderType] struct {
	csv.GroupedCsvReader[T]
	csv.XivApiCsvInfo[T]
}

func (gcr GroupedXivApiCsvReader[T]) ProcessCsv() (results interface{}, err error) {
	defer func() {
		if err == nil && results == nil {
			err = errors.New("no results found")
		}
	}()

	reader, readCloser, err := gcr.GetReader()
	if err != nil {
		return nil, fmt.Errorf("unable to get reader: %w", err)
	}

	defer func(closer io.ReadCloser) {
		err := closer.Close()
		if err != nil {
			log.Fatalf("Error closing response body: %v", err)
		}
	}(readCloser)

	if err = gcr.SkipHeaderRows(reader); err != nil {
		return nil, fmt.Errorf("couldn't skip header rows: %w", err)
	}

	if results, err = gcr.ReadCsvData(reader); err != nil {
		return nil, fmt.Errorf("couldn't read csv data: %w", err)
	}

	return results, nil
}

func (gcr GroupedXivApiCsvReader[T]) ReadCsvData(reader *csvEncoding.Reader) (map[int][]T, error) {
	var wg sync.WaitGroup
	extractionChannel := make(chan T)
	results := make(map[int][]T)

	// Append results in channel to final array
	go func() {
		for result := range extractionChannel {
			valuesInResult := results[result.GetKey()]

			if valuesInResult == nil {
				valuesInResult = make([]T, 0)
			}

			valuesInResult = append(valuesInResult, result)

			results[result.GetKey()] = valuesInResult
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

			item, err := gcr.ProcessRow(record)

			if item == nil || err != nil {
				return
			}

			extractionChannel <- *item
		}(record)
	}

	wg.Wait()
	close(extractionChannel)

	return results, nil
}

func (gcr GroupedXivApiCsvReader[T]) GetReaderType() string {
	return gcr.FileName
}
