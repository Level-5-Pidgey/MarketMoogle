/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (itemprovider.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package db

import (
	schema "MarketMoogleAPI/core/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"log"
)

type ItemDatabaseProvider struct {
	db *DatabaseClient
	collectionName string
}

func NewItemDataBaseProvider(dbClient *DatabaseClient) *ItemDatabaseProvider {
	return &ItemDatabaseProvider{
		db: dbClient,
		collectionName: "items",
	}
}

func (itemProv ItemDatabaseProvider) InsertItem(ctx context.Context, input *schema.Item) (*schema.Item, error) {
	_, err := itemProv.db.InsertOne(ctx, itemProv.collectionName, input)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return input, nil
}

func (itemProv ItemDatabaseProvider) FindItemByObjectId(ctx context.Context, objectIdString string) (*schema.Item, error) {
	objectID, err := primitive.ObjectIDFromHex(objectIdString)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	findResult, err := itemProv.db.FindOne(ctx, itemProv.collectionName, bson.M{"_id": objectID})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	item := schema.Item{}
	err = findResult.Decode(&item)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &item, nil
}

func (itemProv ItemDatabaseProvider) FindItemByItemId(ctx context.Context, itemId int) (*schema.Item, error) {
	findResult, err := itemProv.db.FindOne(ctx, itemProv.collectionName, bson.M{"itemid": itemId})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	item := schema.Item{}
	err = findResult.Decode(&item)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &item, nil
}

func (itemProv ItemDatabaseProvider) GetAllItems(ctx context.Context) ([]*schema.Item, error) {
	collection := itemProv.db.client.Database(itemProv.db.databaseName).Collection(itemProv.collectionName)
	cursor, err := collection.Find(ctx, bson.M{})

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cursor, ctx)

	var items []*schema.Item
	for cursor.Next(ctx) {
		item := schema.Item{}
		err = cursor.Decode(&item)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		items = append(items, &item)
	}

	return items, nil
}

func (itemProv ItemDatabaseProvider) UpdateVendorSellPrice(ctx context.Context, itemId *int, newPrice *int) error {
	_, err := itemProv.db.UpdateOne(
		ctx,
		itemProv.collectionName,
		bson.M{"itemid": *itemId},
		bson.D{
			{"$set", bson.D{
				{"buyfromvendorvalue", *newPrice},
			}},
		})

	return err
}
