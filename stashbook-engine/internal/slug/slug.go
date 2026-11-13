// Package slug converts game titles to URL-safe slugs for sitemap URLs.
package slug

import (
	"regexp"
	"strings"
)

var (
	slugNonAlphaNum = regexp.MustCompile(`[^a-z0-9\s-]`)
	slugSpaces      = regexp.MustCompile(`[\s_]+`)
	slugHyphens     = regexp.MustCompile(`-+`)
)

// Slugify lowercases text and turns it into hyphen-separated segments,
// dropping any character that is not a letter, digit, space or hyphen.
func Slugify(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))
	text = slugNonAlphaNum.ReplaceAllString(text, "")
	text = slugSpaces.ReplaceAllString(text, "-")
	text = slugHyphens.ReplaceAllString(text, "-")
	return strings.Trim(text, "-")
}
