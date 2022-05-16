/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (dbprovider.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type DbProvider struct {
	databaseName string
	client       *mongo.Client
}

func NewDbProvider() *DbProvider {
	return Connect()
}

func Connect() *DbProvider {
	var credentials options.Credential
	credentials.AuthSource = "admin"
	credentials.Password = "access123!"
	credentials.Username = "root"

	hostname := os.Getenv("MONGO_HOST")

	if hostname == "" {
		hostname = "localhost"
	}

	client, err := mongo.NewClient(
		options.Client().
			ApplyURI(fmt.Sprintf("mongodb://%s:27017", hostname)).
			SetAuth(credentials))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	return &DbProvider{
		databaseName: "sanctuary",
		client:       client,
	}
}

func (db DbProvider) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.client.Disconnect(ctx)
}

func (db DbProvider) UpsertCollectionAndIndex(collectionName string, index string) bool {
	database := db.client.Database(db.databaseName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.upsertCollection(ctx, database, collectionName)

	if err != nil {
		log.Fatal(err)
	}

	collectionToAddIndex := database.Collection(collectionName)
	collectionToAddIndex.Indexes().List(ctx)

	return false
}

func (db DbProvider) upsertCollection(ctx context.Context, database *mongo.Database, collectionName string) error {
	if db.CollectionExists(ctx, database, collectionName) {
		return nil
	}

	return database.CreateCollection(ctx, collectionName)
}

func (db DbProvider) CollectionExists(ctx context.Context, database *mongo.Database, collectionName string) bool {
	databaseCollections, err := database.ListCollectionNames(ctx, bson.D{})

	if err != nil {
		log.Fatal(err)
		return false
	}

	for _, name := range databaseCollections {
		if name == collectionName {
			return true
		}
	}

	return false
}
