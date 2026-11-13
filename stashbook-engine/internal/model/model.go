// Package model defines the data types shared by the StashBook engine:
// the game and book records stored in SQLite and published as JSON.
package model

// GameResponse is a single game offer as stored in the database and
// published in public/games.json.
type GameResponse struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Platform      string  `json:"platform"`
	ImageURL      string  `json:"image"`
	OriginalPrice float64 `json:"original_price"`
	IsExpired     bool    `json:"is_expired"`
	Description   string  `json:"description"`
	EndDate       string  `json:"end_date"`
	GameURL       string  `json:"game_url"`
	SeoKeywords   string  `json:"seo_keywords"`
}

// BookResponse is a single free book as stored in the database and
// published in public/books.json.
type BookResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	ImageURL    string `json:"image"`
	Description string `json:"description"`
	BookURL     string `json:"book_url"`
	SeoKeywords string `json:"seo_keywords"`
}
