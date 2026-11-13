package assets_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"stashbook-engine/internal/assets"
	"stashbook-engine/internal/model"
	"stashbook-engine/internal/store"
)

func seedStore(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })

	s.SaveGame(model.GameResponse{
		ID: "epic_1", Title: "Test Game", Platform: "Epic Games",
		ImageURL: "https://example.com/img.png", OriginalPrice: 19.99,
		Description: "Tom &amp; Jerry save the day.",
		EndDate:     "Limited Time", GameURL: "https://example.com/game",
	})
	s.SaveGame(model.GameResponse{
		ID: "steam_2", Title: "Bare Game", Platform: "Steam",
	})
	s.SaveBook(model.BookResponse{
		ID: "dbook_1", Title: "Go in Action", Author: "Ada Lovelace",
		ImageURL: "https://example.com/cover.png",
		BookURL:  "https://example.com/book",
	})
	return s
}

func TestWriteGamesJSONAppliesDefaultsAndSanitizes(t *testing.T) {
	s := seedStore(t)
	outDir := t.TempDir()

	if err := assets.WriteGamesJSON(s, outDir); err != nil {
		t.Fatalf("WriteGamesJSON: %v", err)
	}

	raw, err := os.ReadFile(filepath.Join(outDir, "games.json"))
	if err != nil {
		t.Fatalf("read games.json: %v", err)
	}

	var games []model.GameResponse
	if err := json.Unmarshal(raw, &games); err != nil {
		t.Fatalf("unmarshal games.json: %v", err)
	}
	if len(games) != 2 {
		t.Fatalf("got %d games, want 2", len(games))
	}

	byID := map[string]model.GameResponse{}
	for _, g := range games {
		byID[g.ID] = g
	}

	withDesc := byID["epic_1"]
	if withDesc.Description != "Tom & Jerry save the day." {
		t.Fatalf("description not sanitized: %q", withDesc.Description)
	}
	if withDesc.SeoKeywords == "" {
		t.Fatal("expected generated SEO keywords for game without any")
	}
	if withDesc.IsExpired {
		t.Fatal("newly saved game is expired, want active")
	}

	bare := byID["steam_2"]
	if bare.Description == "" {
		t.Fatal("expected generated description for game without any")
	}
	if bare.SeoKeywords == "" {
		t.Fatal("expected generated SEO keywords for game without any")
	}
}

func TestWriteBooksJSON(t *testing.T) {
	s := seedStore(t)
	outDir := t.TempDir()

	if err := assets.WriteBooksJSON(s, outDir); err != nil {
		t.Fatalf("WriteBooksJSON: %v", err)
	}

	raw, err := os.ReadFile(filepath.Join(outDir, "books.json"))
	if err != nil {
		t.Fatalf("read books.json: %v", err)
	}

	var books []model.BookResponse
	if err := json.Unmarshal(raw, &books); err != nil {
		t.Fatalf("unmarshal books.json: %v", err)
	}
	if len(books) != 1 {
		t.Fatalf("got %d books, want 1", len(books))
	}
	if books[0].Title != "Go in Action" {
		t.Fatalf("unexpected book title: %q", books[0].Title)
	}
	if books[0].Description == "" || books[0].SeoKeywords == "" {
		t.Fatal("expected generated description and SEO keywords for book without any")
	}
}

func TestWriteSitemapXML(t *testing.T) {
	s := seedStore(t)
	outDir := t.TempDir()

	if err := assets.WriteSitemapXML(s, outDir); err != nil {
		t.Fatalf("WriteSitemapXML: %v", err)
	}

	raw, err := os.ReadFile(filepath.Join(outDir, "sitemap.xml"))
	if err != nil {
		t.Fatalf("read sitemap.xml: %v", err)
	}
	sitemap := string(raw)

	for _, want := range []string{
		"https://stashbook.online/",
		"https://stashbook.online/books",
		"https://stashbook.online/iconic",
		"https://stashbook.online/about",
		"https://stashbook.online/contact",
		"https://stashbook.online/privacy",
		"https://stashbook.online/game/test-game",
		"https://stashbook.online/game/bare-game",
	} {
		if !strings.Contains(sitemap, want) {
			t.Fatalf("sitemap is missing %s\n%s", want, sitemap)
		}
	}
}
