/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (crafttype.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package enum

import (
	"fmt"
	"io"
	"strconv"
)

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
