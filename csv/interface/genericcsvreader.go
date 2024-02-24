package csv

import (
	csvEncoding "encoding/csv"
	"errors"
	"fmt"
	"github.com/level-5-pidgey/MarketMoogle/domain"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type GenericCsvReader interface {
	GetReaderType() string

	ProcessCsv() (interface{}, error)
}

type XivApiCsvInfo[T domain.ReaderType] struct {
	FileName   string
	RowsToSkip int
	ProcessRow func([]string) (*T, error)
}

// SkipHeaderRows skips a specified number of rows in a CSV file.
//
// It takes a csv.Reader as input and skips the first x rows defined in the reader
// implementation. If the number of rows to skip is 0, it returns an error.
//
// Parameters:
// - reader: a pointer to the csv.Reader that represents the CSV file to read.
//
// Returns an error indicating if any issues occurred during row skipping.
func (gcr XivApiCsvInfo[T]) SkipHeaderRows(reader *csvEncoding.Reader) error {
	if gcr.RowsToSkip == 0 {
		return errors.New("'RowsToSkip' must be defined")
	}

	// Skips the first x rows defined in the reader implementation
	for i := 0; i < gcr.RowsToSkip; i++ {
		if _, err := reader.Read(); err != nil {
			if err == io.EOF {
				break
			}

			return err
		}
	}

	return nil
}

func (gcr XivApiCsvInfo[T]) getCsvPath() string {
	return fmt.Sprintf("./dl/%s.csv", strings.ToLower(gcr.FileName))
}

// GetReader retrieves a csv.Reader for the UngroupedXivApiCsvReader type.
//
// It makes an HTTP GET request to the `FileName` specified in the UngroupedXivApiCsvReader type.
// The response body is closed automatically using a deferred function.
// Returns:
// The function returns a *csv.Reader and an error.
func (gcr XivApiCsvInfo[T]) GetReader() (*csvEncoding.Reader, io.ReadCloser, error) {
	if !fileExists(gcr.getCsvPath()) {
		// Make folder if this is the first time downloading anything
		newPath := filepath.Dir(gcr.getCsvPath())
		err := os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating folder: %v", err)
		}

		out, err := os.Create(gcr.getCsvPath())
		defer func(out *os.File) {
			err := out.Close()
			if err != nil {
				log.Fatalf("Error closing file: %v", err)
			}
		}(out)

		resp, err := http.Get(
			fmt.Sprintf(
				"https://raw.githubusercontent.com/xivapi/ffxiv-datamining/master/csv/%s.csv",
				gcr.FileName,
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

		fmt.Printf("Downloaded %s.csv.\n", gcr.FileName)
	}

	f, err := os.Open(gcr.getCsvPath())
	if err != nil {
		log.Fatalf("Error closing response body: %v", err)
	}

	reader := csvEncoding.NewReader(f)
	return reader, io.ReadCloser(f), nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)

	return !errors.Is(err, os.ErrNotExist)
}
