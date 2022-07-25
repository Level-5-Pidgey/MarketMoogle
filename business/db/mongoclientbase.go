/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (mongoclientbase.go) is part of MarketMoogleAPI and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClientBase interface {
	Connect() (*mongo.Client, error)
	
	CollectionExists(ctx context.Context, database *mongo.Database, collectionName string) (bool, error)
	UpsertCollection(ctx context.Context, database *mongo.Database, collectionName string) error

	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*mongo.SingleResult, error)
	FindAll(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]interface{}, error)

	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)

	ReplaceOne(ctx context.Context, filter interface{}, options ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
}
