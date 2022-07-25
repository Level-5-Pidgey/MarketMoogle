/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (marketboardprovider.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package db

import (
	"MarketMoogleAPI/core/apitypes/universalis"
	schema "MarketMoogleAPI/core/graph/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MarketboardDatabaseProvider struct {
	db                *DatabaseClient
	listingsPerServer int
	collectionName    string
}

func NewMarketboardDatabaseProvider(dbClient *DatabaseClient) *MarketboardDatabaseProvider {
	return &MarketboardDatabaseProvider{
		db:                dbClient,
		listingsPerServer: 8,
		collectionName:    "marketboard",
	}
}

func (mbProv MarketboardDatabaseProvider) CreateMarketEntry(ctx context.Context, entryFromApi *schema.MarketboardEntry) (*schema.MarketboardEntry, error) {
	_, err := mbProv.db.InsertOne(ctx, mbProv.collectionName, entryFromApi)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return entryFromApi, nil
}

func (mbProv MarketboardDatabaseProvider) ReplaceMarketEntry(ctx context.Context, itemId int, dataCenter string, newEntry *universalis.MarketQuery, currentTimestamp *string) error {
	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"itemid": itemId, "datacenter": dataCenter}

	updatedEntry := schema.MarketboardEntry{
		ItemID:              itemId,
		LastUpdateTime:      *currentTimestamp,
		MarketEntries:       newEntry.GetMarketEntries(mbProv.listingsPerServer),
		MarketHistory:       newEntry.GetItemHistory(mbProv.listingsPerServer),
		DataCenter:          dataCenter,
		CurrentAveragePrice: newEntry.CurrentAveragePrice,
		CurrentMinPrice:     &newEntry.MinPrice,
		RegularSaleVelocity: newEntry.RegularSaleVelocity,
		HqSaleVelocity:      newEntry.HqSaleVelocity,
		NqSaleVelocity:      newEntry.NqSaleVelocity,
	}

	_, err := mbProv.db.ReplaceOne(ctx, mbProv.collectionName, filter, updatedEntry, opts)

	return err
}

func (mbProv MarketboardDatabaseProvider) FindMarketboardEntryByObjectId(ctx context.Context, objectId string) (*schema.MarketboardEntry, error) {
	objectID, err := primitive.ObjectIDFromHex(objectId)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	findResult, err := mbProv.db.FindOne(ctx, mbProv.collectionName, bson.M{"_id": objectID})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	marketboardEntry := schema.MarketboardEntry{}
	err = findResult.Decode(&marketboardEntry)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &marketboardEntry, nil
}

func (mbProv MarketboardDatabaseProvider) FindMarketboardEntriesByItemId(ctx context.Context, itemId int) ([]*schema.MarketboardEntry, error) {
	return mbProv.findMarketboardEntriesBy(ctx, bson.M{"itemid": itemId})
}

func (mbProv MarketboardDatabaseProvider) FindItemEntryOnDc(ctx context.Context, itemId int, dataCenter string) (*schema.MarketboardEntry, error) {
	collection := mbProv.db.client.Database(mbProv.db.databaseName).Collection(mbProv.collectionName)
	cursor, err := collection.Find(ctx, bson.M{"itemid": itemId, "datacenter": dataCenter})

	marketEntry := schema.MarketboardEntry{}
	for cursor.Next(ctx) {
		err = cursor.Decode(&marketEntry)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}
	
	//Return nil if nothing was found
	if marketEntry.ItemID == 0 {
		return nil, nil
	}
	
	return &marketEntry, nil
}

func (mbProv MarketboardDatabaseProvider) GetAllMarketboardEntries(ctx context.Context) ([]*schema.MarketboardEntry, error) {
	return mbProv.findMarketboardEntriesBy(ctx, bson.M{})
}

func (mbProv MarketboardDatabaseProvider) findMarketboardEntriesBy(ctx context.Context, filter bson.M) ([]*schema.MarketboardEntry, error) {
	collection := mbProv.db.client.Database(mbProv.db.databaseName).Collection(mbProv.collectionName)
	cursor, err := collection.Find(ctx, filter)

	var marketboardEntries []*schema.MarketboardEntry
	for cursor.Next(ctx) {
		marketEntry := schema.MarketboardEntry{}
		err = cursor.Decode(&marketEntry)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		marketboardEntries = append(marketboardEntries, &marketEntry)
	}

	return marketboardEntries, nil
}
