package utils

import (
	"strings"
)

func A1ToInt(val string) int64 {
	upperVal := strings.ToUpper(val)
	idx := 0
	for i := range upperVal {
		idx += 26*i + (int(upperVal[i]) - 65)
	}
	return int64(idx)
}
