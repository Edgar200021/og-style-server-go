package utils

import "strings"

func IsContainsSubstring(query, substring string) bool {

	return strings.Contains(strings.ToLower(query), strings.ToLower(substring))
}
