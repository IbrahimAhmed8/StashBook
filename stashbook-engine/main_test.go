package main

import "testing"

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
			want:  "\"StashBook\"",
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
			if got := sanitizeText(tt.input); got != tt.want {
				t.Fatalf("sanitizeText(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSanitizeGameResponse(t *testing.T) {
	game := GameResponse{
		Title:       "ARK &amp; Survival",
		Platform:    "Steam",
		Description: "&quot;Build&quot; &amp; survive",
		SeoKeywords: "free &amp; fun",
	}

	sanitizeGameResponse(&game)

	if game.Title != "ARK & Survival" {
		t.Fatalf("unexpected title: %q", game.Title)
	}
	if game.Description != "\"Build\" & survive" {
		t.Fatalf("unexpected description: %q", game.Description)
	}
	if game.SeoKeywords != "free & fun" {
		t.Fatalf("unexpected seo keywords: %q", game.SeoKeywords)
	}
}
