/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (classrecipe.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package xivapi

import schema "MarketMoogleAPI/core/graph/model"

type ClassRecipe struct {
	AmountIngredient0 int `json:"AmountIngredient0"`
	AmountIngredient1 int `json:"AmountIngredient1"`
	AmountIngredient2 int `json:"AmountIngredient2"`
	AmountIngredient3 int `json:"AmountIngredient3"`
	AmountIngredient4 int `json:"AmountIngredient4"`
	AmountIngredient5 int `json:"AmountIngredient5"`
	AmountIngredient6 int `json:"AmountIngredient6"`
	AmountIngredient7 int `json:"AmountIngredient7"`
	AmountIngredient8 int `json:"AmountIngredient8"`
	AmountIngredient9 int `json:"AmountIngredient9"`
	AmountResult      int `json:"AmountResult"`
	CanHq             int `json:"CanHq"`
	CanQuickSynth     int `json:"CanQuickSynth"`
	CraftType         struct {
		ID           int    `json:"ID"`
		MainPhysical int    `json:"MainPhysical"`
		Name         string `json:"Name"`
		NameDe       string `json:"Name_de"`
		NameEn       string `json:"Name_en"`
		NameFr       string `json:"Name_fr"`
		NameJa       string `json:"Name_ja"`
		SubPhysical  int    `json:"SubPhysical"`
	} `json:"CraftType"`
	CraftTypeTarget          string                 `json:"CraftTypeTarget"`
	CraftTypeTargetID        int                    `json:"CraftTypeTargetID"`
	DifficultyFactor         int                    `json:"DifficultyFactor"`
	DurabilityFactor         int                    `json:"DurabilityFactor"`
	ExpRewarded              int                    `json:"ExpRewarded"`
	ID                       int                    `json:"ID"`
	IsExpert                 int                    `json:"IsExpert"`
	IsSecondary              int                    `json:"IsSecondary"`
	ItemIngredient0          GameItem               `json:"ItemIngredient0"`
	ItemIngredient1          GameItem               `json:"ItemIngredient1"`
	ItemIngredient2          GameItem               `json:"ItemIngredient2"`
	ItemIngredient3          GameItem               `json:"ItemIngredient3"`
	ItemIngredient4          GameItem               `json:"ItemIngredient4"`
	ItemIngredient5          GameItem               `json:"ItemIngredient5"`
	ItemIngredient6          GameItem               `json:"ItemIngredient6"`
	ItemIngredient7          GameItem               `json:"ItemIngredient7"`
	ItemIngredient8          GameItem               `json:"ItemIngredient8"`
	ItemIngredient9          GameItem               `json:"ItemIngredient9"`
	ItemResult               GameItem               `json:"ItemResult"`
	MaterialQualityFactor    int                    `json:"MaterialQualityFactor"`
	Number                   int                    `json:"Number"`
	PatchNumber              int                    `json:"PatchNumber"`
	QualityFactor            int                    `json:"QualityFactor"`
	QuickSynthControl        int                    `json:"QuickSynthControl"`
	QuickSynthCraftsmanship  int                    `json:"QuickSynthCraftsmanship"`
	RecipeLevelTable         RecipeLevelInformation `json:"RecipeLevelTable"`
	RecipeLevelTableTarget   string                 `json:"RecipeLevelTableTarget"`
	RecipeLevelTableTargetID int                    `json:"RecipeLevelTableTargetID"`
	RecipeNotebookList       int                    `json:"RecipeNotebookList"`
	RequiredControl          int                    `json:"RequiredControl"`
	RequiredCraftsmanship    int                    `json:"RequiredCraftsmanship"`
	SecretRecipeBook         SecretRecipeBook       `json:"SecretRecipeBook"`
	SecretRecipeBookTargetID int                    `json:"SecretRecipeBookTargetID"`
}

type ItemsAndQuant struct {
	ItemID   int
	Quantity int
}

func (r ClassRecipe) ConvertToSchemaRecipe(craftType *schema.CraftType) schema.Recipe {

	var items []*schema.RecipeContents
	for _, itemAndQuant := range r.GetRecipeItemsAndQuant() {
		contents := schema.RecipeContents{
			ItemID: itemAndQuant.ItemID,
			Count:  itemAndQuant.Quantity,
		}

		items = append(items, &contents)
	}

	return schema.Recipe{
		RecipeID:               r.ID,
		ItemResultID:           r.ItemResult.ID,
		ResultQuantity:         r.AmountResult,
		CraftedBy:              *craftType,
		RecipeLevel:            &r.RecipeLevelTable.ClassJobLevel,
		MasteryStars:           &r.RecipeLevelTable.Stars,
		RecipeItems:            items,
		SuggestedControl:       &r.RequiredControl,
		SuggestedCraftsmanship: &r.RequiredCraftsmanship,
		Durability:             &r.RecipeLevelTable.Durability,
	}
}

func (r ClassRecipe) GetRecipeItemsAndQuant() []ItemsAndQuant {
	var result []ItemsAndQuant

	if r.AmountIngredient0 != 0 {
		itemInfo := ItemsAndQuant{
			ItemID:   r.ItemIngredient0.ID,
			Quantity: r.AmountIngredient0,
		}
		result = append(result, itemInfo)
	}

	if r.AmountIngredient1 != 0 {
		itemInfo := ItemsAndQuant{
			ItemID:   r.ItemIngredient1.ID,
			Quantity: r.AmountIngredient1,
		}
		result = append(result, itemInfo)
	}

	if r.AmountIngredient2 != 0 {
		itemInfo := ItemsAndQuant{
			ItemID:   r.ItemIngredient2.ID,
			Quantity: r.AmountIngredient2,
		}
		result = append(result, itemInfo)
	}

	if r.AmountIngredient3 != 0 {
		itemInfo := ItemsAndQuant{
			ItemID:   r.ItemIngredient3.ID,
			Quantity: r.AmountIngredient3,
		}
		result = append(result, itemInfo)
	}

	if r.AmountIngredient4 != 0 {
		itemInfo := ItemsAndQuant{
			ItemID:   r.ItemIngredient4.ID,
			Quantity: r.AmountIngredient4,
		}
		result = append(result, itemInfo)
	}

	if r.AmountIngredient5 != 0 {
		itemInfo := ItemsAndQuant{
			ItemID:   r.ItemIngredient5.ID,
			Quantity: r.AmountIngredient5,
		}
		result = append(result, itemInfo)
	}

	if r.AmountIngredient6 != 0 {
		itemInfo := ItemsAndQuant{
			ItemID:   r.ItemIngredient6.ID,
			Quantity: r.AmountIngredient6,
		}
		result = append(result, itemInfo)
	}

	if r.AmountIngredient7 != 0 {
		itemInfo := ItemsAndQuant{
			ItemID:   r.ItemIngredient7.ID,
			Quantity: r.AmountIngredient7,
		}
		result = append(result, itemInfo)
	}

	if r.AmountIngredient8 != 0 {
		itemInfo := ItemsAndQuant{
			ItemID:   r.ItemIngredient8.ID,
			Quantity: r.AmountIngredient8,
		}
		result = append(result, itemInfo)
	}

	if r.AmountIngredient9 != 0 {
		itemInfo := ItemsAndQuant{
			ItemID:   r.ItemIngredient9.ID,
			Quantity: r.AmountIngredient9,
		}
		result = append(result, itemInfo)
	}

	return result
}
