/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (timestamp.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package custom

import (
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"io"
	"strconv"
	"time"
)

func MarshalTimestamp(t time.Time) graphql.Marshaler {
	if t.IsZero() {
		return graphql.Null
	}

	return graphql.WriterFunc(func(writer io.Writer) {
		io.WriteString(writer, strconv.Quote(t.Format(time.StampMilli)))
	})
}

func UnmarshalTimestamp(v interface{}) (time.Time, error) {
	if tmpStr, ok := v.(string); ok {
		return time.Parse(time.StampMilli, tmpStr)
	}
	return time.Time{}, errors.New("time should be StampMilli formatted string")
}
