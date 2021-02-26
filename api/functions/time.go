package functions

import (
	"time"
)

func ToKST(utc string) string {
	t, _ := time.Parse(time.RFC3339, utc)
	t = t.Add(9 * time.Hour)
	return t.Format("2006-01-02 15:04:05")
}
