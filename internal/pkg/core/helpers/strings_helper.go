package helpers

import "strings"

// ReplaceLast replaces the last occurrence of pattern in str with replacement.
// If pattern is not found in str, str is returned unchanged.
func ReplaceLast(str, pattern, replacement string) string {
	if pattern == "" {
		return str
	}
	lastIndex := strings.LastIndex(str, pattern)
	if lastIndex == -1 {
		return str
	}
	return str[:lastIndex] + replacement + str[lastIndex+len(pattern):]
}
