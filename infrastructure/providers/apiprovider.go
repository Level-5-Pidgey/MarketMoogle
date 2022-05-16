/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (apiprovider.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package providers

import (
	"MarketMoogleAPI/core/apitypes/xivapi"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

func MakePaginatedRequest(contentType string, page *int) (*xivapi.PaginatedContent, error) {
	url := fmt.Sprintf("https://xivapi.com/%s?page=%d", contentType, *page)

	return MakeApiRequest[xivapi.PaginatedContent](url)
}

func MakeXivApiContentRequest[T any](contentType string, id *int) (*T, error) {
	url := fmt.Sprintf("https://xivapi.com/%s/%d", contentType, *id)

	return MakeApiRequest[T](url)
}

func MakeApiRequest[T any](urlString string) (*T, error) {
	resp, requestError := http.Get(urlString)
	if requestError != nil {
		log.Fatal(requestError)
		return nil, requestError
	}

	body, readAllError := ioutil.ReadAll(resp.Body)
	if readAllError != nil {
		log.Fatal(readAllError)
		return nil, readAllError
	}

	var responseObject T
	var empty T
	err := json.Unmarshal(body, &responseObject)
	if err != nil {
		return nil, err
	}

	// Check if the response object is empty
	if reflect.DeepEqual(responseObject, empty) {
		return nil, errors.New("response object is empty")
	}

	return &responseObject, nil
}
