/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (recipelookup.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package xivapi

import (
	schema "MarketMoogleAPI/core/graph/model"
)

type RecipeLookup struct {
	ID          int         `json:"ID"`
	ALC         ClassRecipe `json:"ALC"`
	ALCTargetID int         `json:"ALCTargetID"`
	ARM         ClassRecipe `json:"ARM"`
	ARMTargetID int         `json:"ARMTargetID"`
	BSM         ClassRecipe `json:"BSM"`
	BSMTargetID int         `json:"BSMTargetID"`
	CRP         ClassRecipe `json:"CRP"`
	CRPTargetID int         `json:"CRPTargetID"`
	CUL         ClassRecipe `json:"CUL"`
	CULTargetID int         `json:"CULTargetID"`
	GSM         ClassRecipe `json:"GSM"`
	GSMTargetID int         `json:"GSMTargetID"`
	LTW         ClassRecipe `json:"LTW"`
	LTWTargetID int         `json:"LTWTargetID"`
	WVR         ClassRecipe `json:"WVR"`
	WVRTargetID int         `json:"WVRTargetID"`
}

func (r RecipeLookup) GetRecipes() map[schema.CraftType]ClassRecipe {
	result := make(map[schema.CraftType]ClassRecipe)

	//Alchemist
	if r.ALCTargetID != 0 {
		result[schema.CraftTypeAlchemist] = r.ALC
	}

	//Armourer
	if r.ARMTargetID != 0 {
		result[schema.CraftTypeArmourer] = r.ARM
	}

	//Blacksmith
	if r.BSMTargetID != 0 {
		result[schema.CraftTypeBlacksmith] = r.BSM
	}

	//Carpenter
	if r.CRPTargetID != 0 {
		result[schema.CraftTypeCarpenter] = r.CRP
	}

	//Culinarian
	if r.CULTargetID != 0 {
		result[schema.CraftTypeCulinarian] = r.CUL
	}

	//Goldsmith
	if r.GSMTargetID != 0 {
		result[schema.CraftTypeGoldsmith] = r.GSM
	}

	//Leatherworker
	if r.LTWTargetID != 0 {
		result[schema.CraftTypeLeatherworker] = r.LTW
	}

	//Weaver
	if r.WVRTargetID != 0 {
		result[schema.CraftTypeWeaver] = r.WVR
	}

	return result
}

func (r RecipeLookup) GetRecipeItems() map[schema.CraftType][]ItemsAndQuant {
	result := make(map[schema.CraftType][]ItemsAndQuant)

	//Iterate through map and then get all
	recipes := r.GetRecipes()
	for key, value := range recipes {
		result[key] = value.GetRecipeItemsAndQuant()
	}

	return result
}
