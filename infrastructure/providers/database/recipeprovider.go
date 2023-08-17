/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (recipeprovider.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package database

import (
	interfaces "MarketMoogle/business/database"
	schema "MarketMoogle/core/graph/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type RecipeDatabaseProvider struct {
	db             interfaces.Client
	collectionName string
}

func NewRecipeDatabaseProvider(dbClient interfaces.Client) *RecipeDatabaseProvider {
	return &RecipeDatabaseProvider{
		db:             dbClient,
		collectionName: "recipes",
	}
}

func (recipeProv RecipeDatabaseProvider) InsertRecipe(ctx context.Context, recipeToAdd *schema.Recipe) (
	*schema.Recipe, error,
) {
	_, err := recipeProv.db.InsertOne(ctx, recipeProv.collectionName, recipeToAdd)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return recipeToAdd, nil
}

func (recipeProv RecipeDatabaseProvider) FindRecipesByItemId(ctx context.Context, itemId int) (
	[]*schema.Recipe, error,
) {
	return recipeProv.findRecipesBy(ctx, bson.M{"itemresultid": itemId})
}

func (recipeProv RecipeDatabaseProvider) FindRecipeByObjectId(
	ctx context.Context, objectIdString string,
) (*schema.Recipe, error) {
	objectID, err := primitive.ObjectIDFromHex(objectIdString)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	findResult, err := recipeProv.db.FindOne(ctx, recipeProv.collectionName, bson.M{"_id": objectID})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	recipe := schema.Recipe{}
	err = findResult.Decode(&recipe)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &recipe, nil
}

func (recipeProv RecipeDatabaseProvider) GetAllRecipes(ctx context.Context) ([]*schema.Recipe, error) {
	return recipeProv.findRecipesBy(ctx, bson.M{})
}

func (recipeProv RecipeDatabaseProvider) findRecipesBy(ctx context.Context, filter bson.M) ([]*schema.Recipe, error) {
	collection, err := recipeProv.db.GetCollection(recipeProv.collectionName)

	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, filter)

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cursor, ctx)

	var recipes []*schema.Recipe
	for cursor.Next(ctx) {
		recipe := schema.Recipe{}
		err = cursor.Decode(&recipe)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		recipes = append(recipes, &recipe)
	}

	return recipes, nil
}
