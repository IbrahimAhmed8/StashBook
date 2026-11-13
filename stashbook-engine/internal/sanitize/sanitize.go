// Package sanitize cleans text coming from storefront APIs before it is
// published: API responses often contain HTML entities such as &amp;.
package sanitize

import (
	"html"
	"strings"

	"stashbook-engine/internal/model"
)

// SanitizeText decodes HTML entities and trims surrounding whitespace.
func SanitizeText(value string) string {
	return strings.TrimSpace(html.UnescapeString(value))
}

// SanitizeGame cleans every text field of a game record in place.
func SanitizeGame(game *model.GameResponse) {
	game.Title = SanitizeText(game.Title)
	game.Platform = SanitizeText(game.Platform)
	game.ImageURL = SanitizeText(game.ImageURL)
	game.Description = SanitizeText(game.Description)
	game.EndDate = SanitizeText(game.EndDate)
	game.GameURL = SanitizeText(game.GameURL)
	game.SeoKeywords = SanitizeText(game.SeoKeywords)
}

// SanitizeBook cleans every text field of a book record in place.
func SanitizeBook(book *model.BookResponse) {
	book.Title = SanitizeText(book.Title)
	book.Author = SanitizeText(book.Author)
	book.ImageURL = SanitizeText(book.ImageURL)
	book.Description = SanitizeText(book.Description)
	book.BookURL = SanitizeText(book.BookURL)
	book.SeoKeywords = SanitizeText(book.SeoKeywords)
}
