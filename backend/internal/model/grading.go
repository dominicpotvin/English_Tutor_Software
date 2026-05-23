package model

import "strings"

// Normalize canonicalises a free-text answer for comparison: lower-cased,
// trimmed, internal whitespace collapsed and trailing punctuation removed.
func Normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.Join(strings.Fields(s), " ")
	return strings.TrimRight(s, " .!?;,")
}

// Grade reports whether the given answer is correct for the exercise. The
// stored answer may list several acceptable values separated by "|".
func Grade(e Exercise, given string) bool {
	g := Normalize(given)
	if g == "" {
		return false
	}
	for _, acceptable := range strings.Split(e.Answer, "|") {
		if Normalize(acceptable) == g {
			return true
		}
	}
	return false
}
