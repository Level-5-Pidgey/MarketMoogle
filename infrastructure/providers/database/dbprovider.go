/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (dbprovider.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type DatabaseClient struct {
	client        *mongo.Client
	databaseName  string
	dbCredentials *options.Credential
	dbUri         string
}

func NewDatabaseClient(dbName string, uri string, credentials options.Credential) *DatabaseClient {
	clientConnection, err := mongo.NewClient(options.Client().
		ApplyURI(uri).
		SetAuth(credentials))
	if err != nil {
		log.Fatal(err)
	}

	err = clientConnection.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	//Verify connection
	err = clientConnection.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return &DatabaseClient{
		client:        clientConnection,
		databaseName:  dbName,
		dbCredentials: &credentials,
		dbUri:         uri,
	}
}

func (dbClient DatabaseClient) CollectionExists(ctx context.Context, collectionName string, database *mongo.Database) (bool, error) {
	databaseCollections, err := database.ListCollectionNames(ctx, bson.D{})

	if err != nil {
		return false, err
	}

	for _, name := range databaseCollections {
		if name == collectionName {
			return true, nil
		}
	}

	return false, nil
}

func (dbClient DatabaseClient) UpsertCollection(ctx context.Context, collectionName string, database *mongo.Database) error {
	result, err := dbClient.CollectionExists(ctx, collectionName, database)

	if result && err == nil {
		return nil
	}

	return database.CreateCollection(ctx, collectionName)
}

func (dbClient DatabaseClient) FindOne(ctx context.Context, collectionName string, filter interface{}, opts ...*options.FindOneOptions) (*mongo.SingleResult, error) {
	collection := dbClient.client.Database(dbClient.databaseName).Collection(collectionName)
	result := collection.FindOne(ctx, filter, opts...)

	return result, result.Err()
}

func (dbClient DatabaseClient) InsertOne(ctx context.Context, collectionName string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	collection := dbClient.client.Database(dbClient.databaseName).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, document, opts...)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return res, nil
}

func (dbClient DatabaseClient) UpdateOne(ctx context.Context, collectionName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	collection := dbClient.client.Database(dbClient.databaseName).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.UpdateOne(ctx, filter, update, opts...)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return res, nil
}

func (dbClient DatabaseClient) ReplaceOne(ctx context.Context, collectionName string, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	collection := dbClient.client.Database(dbClient.databaseName).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.ReplaceOne(ctx, filter, replacement, opts...)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return res, nil
}

func (dbClient DatabaseClient) CreateIndex(ctx context.Context, collectionName string, keys interface{}, opts *options.IndexOptions) error {
	model := mongo.IndexModel{
		Keys:    keys,
		Options: opts,
	}
	
	_, err := dbClient.client.Database(dbClient.databaseName).Collection(collectionName).Indexes().CreateOne(ctx, model)
	
	return err
}
