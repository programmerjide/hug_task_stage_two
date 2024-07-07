package utils

import (
	"strings"
	"unicode"
)

// IsZero checks if a numeric value is zero.
func IsZero[T comparable](value T) bool {
	var zero T
	return value == zero
}

// IsEmpty checks if a string is empty.
func IsEmpty(s string) bool {
	return len(s) == 0
}

// IsNotEmpty checks if a string is not empty.
func IsNotEmpty(s string) bool {
	return len(s) > 0
}

// Contains checks if a string contains a substring.
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// HasPrefix checks if a string starts with a prefix.
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// HasSuffix checks if a string ends with a suffix.
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// ToLower converts a string to lowercase.
func ToLower(s string) string {
	return strings.ToLower(s)
}

// ToUpper converts a string to uppercase.
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// Trim removes leading and trailing spaces from a string.
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// IsAlpha checks if a string contains only alphabetic characters.
func IsAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// IsNumeric checks if a string contains only numeric characters.
func IsNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsAlphanumeric checks if a string contains only alphanumeric characters.
func IsAlphanumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// Join concatenates a slice of strings into a single string with a separator.
func Join(elems []string, separator string) string {
	return strings.Join(elems, separator)
}

// Split splits a string into a slice of substrings separated by a specified separator.
func Split(s, separator string) []string {
	return strings.Split(s, separator)
}

// Replace replaces all occurrences of old in a string with new.
func Replace(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// ContainsOnlySpaces checks if a string contains only spaces.
func ContainsOnlySpaces(s string) bool {
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}
