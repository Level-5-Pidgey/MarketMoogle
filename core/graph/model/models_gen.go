// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package schema

import (
	"fmt"
	"io"
	"strconv"
)

type Job struct {
	ID    string `json:"Id"`
	JobID int    `json:"JobId"`
	Level int    `json:"Level"`
}

type JobInput struct {
	JobID int `json:"JobId"`
	Level int `json:"Level"`
}

type MarketEntry struct {
	ServerID     int     `json:"ServerId"`
	Server       string  `json:"Server"`
	Quantity     int     `json:"Quantity"`
	TotalPrice   int     `json:"TotalPrice"`
	PricePer     int     `json:"PricePer"`
	Hq           bool    `json:"Hq"`
	IsCrafted    bool    `json:"IsCrafted"`
	RetainerName *string `json:"RetainerName"`
}

type Recipe struct {
	RecipeID               int               `json:"RecipeId"`
	ItemResultID           int               `json:"ItemResultId"`
	ResultQuantity         int               `json:"ResultQuantity"`
	CraftedBy              CraftType         `json:"CraftedBy"`
	RecipeLevel            *int              `json:"RecipeLevel"`
	MasteryStars           *int              `json:"MasteryStars"`
	RecipeItems            []*RecipeContents `json:"RecipeItems"`
	SuggestedControl       *int              `json:"SuggestedControl"`
	SuggestedCraftsmanship *int              `json:"SuggestedCraftsmanship"`
	Durability             *int              `json:"Durability"`
}

type RecipePurchaseInformation struct {
	Item            *Item  `json:"Item"`
	ServerToBuyFrom string `json:"ServerToBuyFrom"`
	BuyFromVendor   bool   `json:"BuyFromVendor"`
	Quantity        int    `json:"Quantity"`
}

type RecipeResaleInformation struct {
	Profit          int                          `json:"Profit"`
	ItemsToPurchase []*RecipePurchaseInformation `json:"ItemsToPurchase"`
	CraftLevel      int                          `json:"CraftLevel"`
	CraftType       CraftType                    `json:"CraftType"`
	ItemCost        int                          `json:"ItemCost"`
}

type User struct {
	ID              int     `json:"Id"`
	LodestoneID     *int    `json:"LodestoneId"`
	Jobs            []*Job  `json:"Jobs"`
	DataCenter      *string `json:"DataCenter"`
	Server          *string `json:"Server"`
	PortraitAddress *string `json:"PortraitAddress"`
	IsPremium       *bool   `json:"IsPremium"`
}

type UserInput struct {
	LodestoneID     *int        `json:"LodestoneId"`
	Jobs            []*JobInput `json:"Jobs"`
	DataCenter      *string     `json:"DataCenter"`
	Server          *string     `json:"Server"`
	PortraitAddress *string     `json:"PortraitAddress"`
	IsPremium       *bool       `json:"IsPremium"`
}

type CraftType string

const (
	CraftTypeCarpenter     CraftType = "CARPENTER"
	CraftTypeBlacksmith    CraftType = "BLACKSMITH"
	CraftTypeArmourer      CraftType = "ARMOURER"
	CraftTypeGoldsmith     CraftType = "GOLDSMITH"
	CraftTypeLeatherworker CraftType = "LEATHERWORKER"
	CraftTypeWeaver        CraftType = "WEAVER"
	CraftTypeAlchemist     CraftType = "ALCHEMIST"
	CraftTypeCulinarian    CraftType = "CULINARIAN"
)

var AllCraftType = []CraftType{
	CraftTypeCarpenter,
	CraftTypeBlacksmith,
	CraftTypeArmourer,
	CraftTypeGoldsmith,
	CraftTypeLeatherworker,
	CraftTypeWeaver,
	CraftTypeAlchemist,
	CraftTypeCulinarian,
}

func (e CraftType) IsValid() bool {
	switch e {
	case CraftTypeCarpenter, CraftTypeBlacksmith, CraftTypeArmourer, CraftTypeGoldsmith, CraftTypeLeatherworker, CraftTypeWeaver, CraftTypeAlchemist, CraftTypeCulinarian:
		return true
	}
	return false
}

func (e CraftType) String() string {
	return string(e)
}

func (e *CraftType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CraftType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CraftType", str)
	}
	return nil
}

func (e CraftType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
