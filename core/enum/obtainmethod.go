/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (obtainmethod.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package enum

import (
	"fmt"
	"io"
	"strconv"
)

type ObtainMethod string

const (
	ObtainMethodCrafted  ObtainMethod = "CRAFTED"
	ObtainMethodDropped  ObtainMethod = "DROPPED"
	ObtainMethodGathered ObtainMethod = "GATHERED"
	ObtainMethodSold     ObtainMethod = "SOLD"
)

var AllObtainMethod = []ObtainMethod{
	ObtainMethodCrafted,
	ObtainMethodDropped,
	ObtainMethodGathered,
	ObtainMethodSold,
}

func (e ObtainMethod) IsValid() bool {
	switch e {
	case ObtainMethodCrafted, ObtainMethodDropped, ObtainMethodGathered, ObtainMethodSold:
		return true
	}
	return false
}

func (e ObtainMethod) String() string {
	return string(e)
}

func (e *ObtainMethod) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ObtainMethod(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ObtainMethod", str)
	}
	return nil
}

func (e ObtainMethod) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
