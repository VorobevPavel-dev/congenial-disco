package utility

import (
	"strings"
)

func StringIsIn(value string, slice []string) bool {
	for _, item := range slice {
		if strings.Compare(item, value) == 0 {
			return true
		}
	}
	return false
}
