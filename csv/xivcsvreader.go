package csv

import (
	csvEncoding "encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type XivCsvReader interface {
	GetReaderType() string

	ProcessCsv() (results interface{}, err error)
}

type GenericXivCsvReader[T readerType[T]] struct {
	RowsToSkip int
	FileName   string
}

func (csvReader GenericXivCsvReader[T]) fileExists(filePath string) bool {
	_, err := os.Stat(filePath)

	return !errors.Is(err, os.ErrNotExist)
}

func (csvReader GenericXivCsvReader[T]) getCsvPath() string {
	return fmt.Sprintf("./dl/%s.csv", strings.ToLower(csvReader.FileName))
}

func (csvReader GenericXivCsvReader[T]) getReader() (*csvEncoding.Reader, io.ReadCloser, error) {
	if !csvReader.fileExists(csvReader.getCsvPath()) {
		// Make folder if this is the first time downloading anything
		newPath := filepath.Dir(csvReader.getCsvPath())
		err := os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating folder: %v", err)
		}

		out, err := os.Create(csvReader.getCsvPath())
		defer func(out *os.File) {
			err := out.Close()
			if err != nil {
				log.Fatalf("Error closing file: %v", err)
			}
		}(out)

		resp, err := http.Get(
			fmt.Sprintf(
				"https://raw.githubusercontent.com/xivapi/ffxiv-datamining/master/csv/%s.csv",
				csvReader.FileName,
			),
		)
		if err != nil {
			log.Fatalf("Error getting file from internet: %v", err)
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatalf("Error closing response body: %v", err)
			}
		}(resp.Body)

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatalf("Error copying file from internet: %v", err)
		}

		fmt.Printf("Downloaded %s.csv.\n", csvReader.FileName)
	}

	f, err := os.Open(csvReader.getCsvPath())
	if err != nil {
		log.Fatalf("Error closing response body: %v", err)
	}

	reader := csvEncoding.NewReader(f)
	return reader, io.ReadCloser(f), nil
}

func (csvReader GenericXivCsvReader[T]) skipHeaderRows(reader *csvEncoding.Reader) error {
	if csvReader.RowsToSkip == 0 {
		return nil
	}

	for i := 0; i < csvReader.RowsToSkip; i++ {
		if _, err := reader.Read(); err != nil {
			if err == io.EOF {
				break
			}

			return err
		}
	}

	return nil
}
