/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (itemprovider.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package db

import (
	"MarketMoogleAPI/core/apitypes/xivapi"
	schema "MarketMoogleAPI/core/graph/model"
	"MarketMoogleAPI/core/util"
	"MarketMoogleAPI/infrastructure/providers"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"log"
	"time"
)

const itemCollection = "items"

func (db DbProvider) SaveItem(input *schema.Item) (*schema.Item, error) {
	collection := db.client.Database(db.databaseName).Collection(itemCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, input)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return input, nil
}

func (db DbProvider) FindItemByObjectId(ID string) (*schema.Item, error) {
	objectID, err := primitive.ObjectIDFromHex(ID)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	collection := db.client.Database(db.databaseName).Collection(itemCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := collection.FindOne(ctx, bson.M{"_id": objectID})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	item := schema.Item{}
	result.Decode(&item)

	return &item, nil
}

func (db DbProvider) FindItemByItemId(ctx context.Context, ItemId int) (*schema.Item, error) {
	collection := db.client.Database(db.databaseName).Collection(itemCollection)
	result := collection.FindOne(ctx, bson.M{"itemid": ItemId})

	item := schema.Item{}
	err := result.Decode(&item)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &item, nil
}

func (db DbProvider) GetAllItems(ctx context.Context) ([]*schema.Item, error) {
	collection := db.client.Database(db.databaseName).Collection(itemCollection)
	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var items []*schema.Item
	for cursor.Next(ctx) {
		var item *schema.Item
		err := cursor.Decode(&item)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (db DbProvider) GetItemFromApi(itemId *int) (*schema.Item, error) {
	//Get different obtain info for item
	prov := providers.XivApiProvider{}

	//Run the api methods asynchronously
	gameItemOut := util.Async(func() *xivapi.GameItem {
		gameItem, err := prov.GetGameItemById(itemId)

		if err != nil {
			log.Fatal(err)
			return nil
		}

		return gameItem
	})

	//Turn into item object
	gameItem := <-gameItemOut

	blankItem := xivapi.GameItem{}
	if *gameItem == blankItem {
		return nil, errors.New("could not find the item specified")
	}

	newItem := schema.Item{
		ItemID:            gameItem.ID,
		Name:              gameItem.Name,
		Description:       &gameItem.Description,
		CanBeHq:           gameItem.CanBeHq == 1,
		IconID:            gameItem.IconID,
		SellToVendorValue: &gameItem.PriceMid,
	}

	//Pass to save method and return result
	return &newItem, nil
}

func (db DbProvider) SaveItemFromApi(itemID *int) (*schema.Item, error) {
	item, err := db.GetItemFromApi(itemID)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db.SaveItem(item)
}

func (db DbProvider) UpdateVendorSellPrice(itemId *int, newPrice *int) error {
	collection := db.client.Database(db.databaseName).Collection(itemCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx,
		bson.M{"itemid": *itemId},
		bson.D{
			{"$set", bson.D{{"buyfromvendorvalue", *newPrice}}},
		})

	return err
}
