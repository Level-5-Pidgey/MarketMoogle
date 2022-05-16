/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (timestamp.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package util

import (
	"strconv"
	"time"
)

func ConvertTimestampStringToTime(timestampString string) (time.Time, error) {
	i, err := strconv.ParseInt(timestampString, 10, 64)

	if err != nil {
		return time.Time{}, err
	}

	tm := time.Unix(0, i)
	return tm, nil
}

func ConvertTimeToTimestampString(dateTime time.Time) string {
	timestamp := strconv.FormatInt(dateTime.UTC().UnixNano(), 10)
	return timestamp
}

func GetCurrentTimestampString() string {
	return ConvertTimeToTimestampString(time.Now().UTC())
}
