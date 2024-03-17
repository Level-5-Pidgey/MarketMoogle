package util

import "strconv"

func SafeStringToBool(conv string) bool {
	if b, err := strconv.ParseBool(conv); err != nil {
		return false
	} else {
		return b
	}
}

func SafeStringToInt(conv string) int {
	if i, err := strconv.Atoi(conv); err != nil {
		return 0
	} else {
		return i
	}
}

func SafeStringToFloat(conv string) float64 {
	if f, err := strconv.ParseFloat(conv, 64); err != nil {
		return 0
	} else {
		return f
	}
}
