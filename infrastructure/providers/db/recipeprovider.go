/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (recipeprovider.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package db

import (
	"MarketMoogleAPI/core/apitypes/xivapi"
	schema "MarketMoogleAPI/core/graph/model"
	"MarketMoogleAPI/core/util"
	"MarketMoogleAPI/infrastructure/providers"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

const recipeCollection = "recipes"

func (db DbProvider) CreateRecipesFromApi(itemID *int) (*[]*schema.Recipe, error) {
	//Get different obtain info for item
	prov := providers.XivApiProvider{}

	itemRecipeOut := util.Async(func() *xivapi.RecipeLookup {
		recipe, err := prov.GetRecipeIdByItemId(itemID)

		if err != nil {
			log.Fatal(err)
		}

		return recipe
	})

	//Turn into item object
	itemRecipe := <-itemRecipeOut

	recipeDict := itemRecipe.GetRecipes()
	var recipes []*schema.Recipe
	for key, value := range recipeDict {
		recipe := value.ConvertToSchemaRecipe(&key)
		recipes = append(recipes, &recipe)
	}

	if len(recipes) > 0 {
		return db.SaveRecipes(&recipes)
	}

	return &recipes, nil
}

func (db DbProvider) SaveRecipes(input *[]*schema.Recipe) (*[]*schema.Recipe, error) {
	collection := db.client.Database(db.databaseName).Collection(recipeCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var inputInterface []interface{}
	for _, v := range *input {
		inputInterface = append(inputInterface, v)
	}

	_, err := collection.InsertMany(ctx, inputInterface)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return input, nil
}

func (db DbProvider) FindRecipesByItemId(ctx context.Context, ItemId int) ([]*schema.Recipe, error) {
	collection := db.client.Database(db.databaseName).Collection(recipeCollection)
	cursor, err := collection.Find(ctx, bson.M{"itemresultid": ItemId})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var recipes []*schema.Recipe
	for cursor.Next(ctx) {
		var recipe *schema.Recipe
		err := cursor.Decode(&recipe)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (db DbProvider) FindRecipeByObjectId(ID string) (*schema.Recipe, error) {
	objectID, err := primitive.ObjectIDFromHex(ID)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	collection := db.client.Database(db.databaseName).Collection(db.databaseName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := collection.FindOne(ctx, bson.M{"_id": objectID})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	recipe := schema.Recipe{}
	err = result.Decode(&recipe)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &recipe, nil
}

func (db DbProvider) GetAllRecipes(ctx context.Context) ([]*schema.Recipe, error) {
	collection := db.client.Database(db.databaseName).Collection(recipeCollection)
	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var recipes []*schema.Recipe
	for cursor.Next(ctx) {
		var recipe *schema.Recipe
		err := cursor.Decode(&recipe)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}
