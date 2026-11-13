package sanitize_test

import (
	"testing"

	"stashbook-engine/internal/model"
	"stashbook-engine/internal/sanitize"
)

func TestSanitizeText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "ampersand entity",
			input: "Tom &amp; Jerry",
			want:  "Tom & Jerry",
		},
		{
			name:  "quote entity",
			input: "&quot;StashBook&quot;",
			want:  `"StashBook"`,
		},
		{
			name:  "apostrophe entity",
			input: "It&#39;s free",
			want:  "It's free",
		},
		{
			name:  "trims surrounding whitespace",
			input: "  Free game on Steam.  ",
			want:  "Free game on Steam.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitize.SanitizeText(tt.input); got != tt.want {
				t.Fatalf("SanitizeText(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSanitizeGame(t *testing.T) {
	game := model.GameResponse{
		Title:       "ARK &amp; Survival",
		Platform:    "Steam",
		Description: "&quot;Build&quot; &amp; survive",
		SeoKeywords: "free &amp; fun",
	}

	sanitize.SanitizeGame(&game)

	if game.Title != "ARK & Survival" {
		t.Fatalf("unexpected title: %q", game.Title)
	}
	if game.Description != `"Build" & survive` {
		t.Fatalf("unexpected description: %q", game.Description)
	}
	if game.SeoKeywords != "free & fun" {
		t.Fatalf("unexpected seo keywords: %q", game.SeoKeywords)
	}
}

func TestSanitizeBook(t *testing.T) {
	book := model.BookResponse{
		Title:       "Go &amp; Friends",
		Author:      "Ada &quot;The Dev&quot;",
		Description: "Learn &amp; build.",
	}

	sanitize.SanitizeBook(&book)

	if book.Title != "Go & Friends" {
		t.Fatalf("unexpected title: %q", book.Title)
	}
	if book.Author != `Ada "The Dev"` {
		t.Fatalf("unexpected author: %q", book.Author)
	}
	if book.Description != "Learn & build." {
		t.Fatalf("unexpected description: %q", book.Description)
	}
}
