package sources

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestFetchEpicFreeGames parses a real sample of the free-games feed
// (testdata/epic_freegames.json, trimmed from the live endpoint) and
// checks that only the currently free game is returned.
func TestFetchEpicFreeGames(t *testing.T) {
	fixture, err := os.ReadFile("testdata/epic_freegames.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture)
	}))
	defer server.Close()

	games, err := fetchEpicFreeGames(server.Client(), server.URL)
	if err != nil {
		t.Fatalf("fetchEpicFreeGames: %v", err)
	}
	if len(games) != 1 {
		t.Fatalf("got %d free games, want 1 (only Knightica is free in the fixture)", len(games))
	}

	g := games[0]
	if g.Title != "Knightica" {
		t.Errorf("Title = %q, want %q", g.Title, "Knightica")
	}
	if g.Price.TotalPrice.DiscountPrice != 0 || g.Price.TotalPrice.OriginalPrice != 699 {
		t.Errorf("price = discount %d / original %d, want 0 / 699",
			g.Price.TotalPrice.DiscountPrice, g.Price.TotalPrice.OriginalPrice)
	}
	// Knightica has a null productSlug in the feed, so callers must fall
	// back to urlSlug; make sure both fields parse as expected.
	if g.ProductSlug != "" {
		t.Errorf("ProductSlug = %q, want empty", g.ProductSlug)
	}
	if g.UrlSlug != "72b46f5d95624fd5b7e7d81b7dcaead1" {
		t.Errorf("UrlSlug = %q, want %q", g.UrlSlug, "72b46f5d95624fd5b7e7d81b7dcaead1")
	}
}

// TestFetchEpicFreeGamesNonOK covers the failure mode that broke the Epic
// engine: the old GraphQL endpoint answers server-side requests with a
// 403 Cloudflare challenge page. The fetch must surface that as an error
// instead of silently returning zero games.
func TestFetchEpicFreeGamesNonOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("cloudflare challenge"))
	}))
	defer server.Close()

	games, err := fetchEpicFreeGames(server.Client(), server.URL)
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
	if len(games) != 0 {
		t.Errorf("got %d games, want 0", len(games))
	}
}

// TestFetchEpicFreeGamesFilter checks the free filter on edge cases: a
// title that is permanently free (original price 0) is not a deal, and a
// discounted-but-not-free title is not free either.
func TestFetchEpicFreeGamesFilter(t *testing.T) {
	body := `{"data":{"Catalog":{"searchStore":{"elements":[
		{"title":"Always Free","id":"1","price":{"totalPrice":{"discountPrice":0,"originalPrice":0}}},
		{"title":"Half Off","id":"2","price":{"totalPrice":{"discountPrice":500,"originalPrice":1000}}},
		{"title":"Free This Week","id":"3","price":{"totalPrice":{"discountPrice":0,"originalPrice":2999}}}
	]}}}}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
	}))
	defer server.Close()

	games, err := fetchEpicFreeGames(server.Client(), server.URL)
	if err != nil {
		t.Fatalf("fetchEpicFreeGames: %v", err)
	}
	if len(games) != 1 || games[0].Title != "Free This Week" {
		t.Fatalf("got %v, want only [Free This Week]", games)
	}
}
