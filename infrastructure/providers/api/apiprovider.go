/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (apiprovider.go) is part of MarketMoogleAPI and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package api

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

func MakePaginatedRequest(contentType string, page int, apiKey string) (*xivapi.PaginatedContent, error) {
	queryOptions := fmt.Sprintf("?page=%d", page)
	if apiKey != "" {
		queryOptions = fmt.Sprintf("%s&page=%d", apiKey, page)
	}
	
	url := fmt.Sprintf("https://xivapi.com/%s%s", contentType, queryOptions)

	return MakeApiRequest[xivapi.PaginatedContent](url)
}

func MakeXivApiContentRequest[T any](contentType string, id int, apiKeyString string) (*T, error) {
	url := fmt.Sprintf("https://xivapi.com/%s/%d%s", contentType, id, apiKeyString)

	return MakeApiRequest[T](url)
}

func MakeApiRequest[T any](urlString string) (*T, error) {
	resp, requestError := http.Get(urlString)
	if requestError != nil {
		log.Fatal(requestError)
		return nil, requestError
	}
	
	//API has a DNS problem or is offline, cancel unmarshalling
	if resp.StatusCode == 522 {
		return nil, errors.New("522 code returned from api request")
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
