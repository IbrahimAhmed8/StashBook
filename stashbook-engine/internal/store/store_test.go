package store_test

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"stashbook-engine/internal/model"
	"stashbook-engine/internal/store"
)

func openTestStore(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestOpenCreatesSchemaOnColdStart(t *testing.T) {
	s := openTestStore(t)

	// A fresh database has no games, books or titles, and no errors.
	games, err := s.Games()
	if err != nil {
		t.Fatalf("Games on fresh DB: %v", err)
	}
	if len(games) != 0 {
		t.Fatalf("fresh DB has %d games, want 0", len(games))
	}

	books, err := s.Books()
	if err != nil {
		t.Fatalf("Books on fresh DB: %v", err)
	}
	if len(books) != 0 {
		t.Fatalf("fresh DB has %d books, want 0", len(books))
	}

	titles, err := s.GameTitles()
	if err != nil {
		t.Fatalf("GameTitles on fresh DB: %v", err)
	}
	if len(titles) != 0 {
		t.Fatalf("fresh DB has %d titles, want 0", len(titles))
	}
}

func TestSaveAndListGames(t *testing.T) {
	s := openTestStore(t)

	s.SaveGame(model.GameResponse{
		ID: "epic_1", Title: "Test Game", Platform: "Epic Games",
		ImageURL: "https://example.com/img.png", OriginalPrice: 19.99,
		Description: "A test game.", EndDate: "Limited Time",
		GameURL: "https://example.com/game", SeoKeywords: "free test game",
	})

	games, err := s.Games()
	if err != nil {
		t.Fatalf("Games: %v", err)
	}
	if len(games) != 1 {
		t.Fatalf("got %d games, want 1", len(games))
	}
	g := games[0]
	if g.Title != "Test Game" || g.Platform != "Epic Games" || g.OriginalPrice != 19.99 {
		t.Fatalf("unexpected game: %+v", g)
	}
	if g.IsExpired {
		t.Fatal("newly saved game is expired, want active")
	}

	// Expiry round-trip: mark all expired, then re-saving re-activates.
	s.MarkAllGamesExpired()
	games, _ = s.Games()
	if !games[0].IsExpired {
		t.Fatal("game is active after MarkAllGamesExpired, want expired")
	}

	s.SaveGame(model.GameResponse{ID: "epic_1", Title: "Test Game", Platform: "Epic Games"})
	games, _ = s.Games()
	if games[0].IsExpired {
		t.Fatal("game is expired after re-saving, want active")
	}
}

func TestSaveAndListBooks(t *testing.T) {
	s := openTestStore(t)

	s.SaveBook(model.BookResponse{
		ID: "dbook_1", Title: "Go in Action", Author: "Ada Lovelace",
		ImageURL: "https://example.com/cover.png", Description: "A book.",
		BookURL: "https://example.com/book", SeoKeywords: "free go book",
	})

	books, err := s.Books()
	if err != nil {
		t.Fatalf("Books: %v", err)
	}
	if len(books) != 1 {
		t.Fatalf("got %d books, want 1", len(books))
	}
	if books[0].Title != "Go in Action" || books[0].Author != "Ada Lovelace" {
		t.Fatalf("unexpected book: %+v", books[0])
	}
}

func TestGameTitles(t *testing.T) {
	s := openTestStore(t)

	s.SaveGame(model.GameResponse{ID: "a", Title: "Alpha"})
	s.SaveGame(model.GameResponse{ID: "b", Title: "Beta"})
	s.SaveGame(model.GameResponse{ID: "c", Title: "Alpha"})

	titles, err := s.GameTitles()
	if err != nil {
		t.Fatalf("GameTitles: %v", err)
	}
	if len(titles) != 2 {
		t.Fatalf("got %d titles %v, want 2", len(titles), titles)
	}
}

func TestPriceHistory(t *testing.T) {
	s := openTestStore(t)

	if _, err := s.LastPrice("123"); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("LastPrice for untracked app: err = %v, want sql.ErrNoRows", err)
	}

	s.InsertPrice("123", "Some Game", 29.99)
	price, err := s.LastPrice("123")
	if err != nil {
		t.Fatalf("LastPrice: %v", err)
	}
	if price != 29.99 {
		t.Fatalf("LastPrice = %v, want 29.99", price)
	}

	s.UpdatePrice("123", 0)
	price, err = s.LastPrice("123")
	if err != nil {
		t.Fatalf("LastPrice after update: %v", err)
	}
	if price != 0 {
		t.Fatalf("LastPrice after update = %v, want 0", price)
	}
}
