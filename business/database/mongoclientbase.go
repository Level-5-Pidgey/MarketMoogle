/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (mongoclientbase.go) is part of MarketMoogleAPI and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client interface {
	GetDatabaseName() string
	GetDatabase(databaseName string) (*mongo.Database, error)
	CollectionExists(ctx context.Context, collectionName string, database *mongo.Database) (bool, error)
	UpsertCollection(ctx context.Context, collectionName string, database *mongo.Database) error
	GetCollection(collectionName string) (*mongo.Collection, error)
	GetCollectionOnDatabase(collectionName string, databaseName string) (*mongo.Collection, error)
	FindOne(
		ctx context.Context, collectionName string, filter interface{}, opts ...*options.FindOneOptions,
	) (*mongo.SingleResult, error)
	InsertOne(
		ctx context.Context, collectionName string, document interface{}, opts ...*options.InsertOneOptions,
	) (*mongo.InsertOneResult, error)
	ReplaceOne(
		ctx context.Context, collectionName string, filter interface{}, replacement interface{},
		opts ...*options.ReplaceOptions,
	) (*mongo.UpdateResult, error)
	UpdateOne(
		ctx context.Context, collectionName string, filter interface{}, update interface{},
		opts ...*options.UpdateOptions,
	) (*mongo.UpdateResult, error)
	CreateIndex(ctx context.Context, collectionName string, keys interface{}, opts *options.IndexOptions) error
}
