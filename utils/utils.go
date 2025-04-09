package utils

import (
	"strings"
	"unicode"
)

// ToPascalCase converts a snake_case string to PascalCase
func ToPascalCase(str string) string {
	words := strings.Split(str, "_")
	for i, word := range words {
		if len(word) > 0 {
			r := []rune(word)
			r[0] = unicode.ToUpper(r[0])
			words[i] = string(r)
		}
	}
	return strings.Join(words, "")
}

// IsReadOnlyQuery checks if a SQL query is read-only
func IsReadOnlyQuery(query string) bool {
	lowerQuery := strings.ToLower(strings.TrimSpace(query))
	return strings.HasPrefix(lowerQuery, "select") &&
		!strings.Contains(lowerQuery, "delete") &&
		!strings.Contains(lowerQuery, "insert") &&
		!strings.Contains(lowerQuery, "update") &&
		!strings.Contains(lowerQuery, "drop") &&
		!strings.Contains(lowerQuery, "alter") &&
		!strings.Contains(lowerQuery, "create")
}
