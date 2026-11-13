// StashBook engine: aggregates free game offers and tech books from
// several storefronts into SQLite, then publishes static JSON and a
// sitemap for the frontend.
//
// Run from this directory: go run main.go
// It writes ../public/games.json, ../public/books.json and
// ../public/sitemap.xml, and keeps its working database in
// ./stashbook.db (not committed; see .gitignore).
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"stashbook-engine/internal/assets"
	"stashbook-engine/internal/run"
	"stashbook-engine/internal/sources"
	"stashbook-engine/internal/store"
)

func main() {
	s, err := store.Open("./stashbook.db")
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	s.MarkAllGamesExpired()

	fmt.Println("Starting Aggregation Engines...")

	engines := []struct {
		label string
		name  string
		fn    func()
	}{
		{"Engine 1: Epic Games", "Epic Games", func() { sources.FetchEpicGames(s) }},
		{"Engine 2: Steam Live Deals", "Steam Live Deals", func() { sources.FetchSteamGames(s) }},
		{"Engine 3: CheapShark Aggregator", "CheapShark Aggregator", func() { sources.FetchCheapSharkDeals(s) }},
		{"Engine 4: StashBook Self-Tracker", "StashBook Self-Tracker", func() { sources.TrackSteamHistory(s) }},
		{"Engine 5: Amazon Prime Gaming (GamerPower API)", "Amazon Prime Gaming", func() { sources.FetchPrimeGames(s) }},
		{"Engine 6: Global Indie & Consoles", "Global Indie & Consoles", func() { sources.FetchGlobalGames(s) }},
		{"Engine 7: Iconic Deals Hunter", "Iconic Deals Hunter", func() { sources.FetchLegendaryDeals(s) }},
		{"Engine 8: Tech Vault (Digital Books Tracker)", "Tech Vault", func() { sources.FetchTechBooks(s) }},
	}

	var wg sync.WaitGroup
	for _, engine := range engines {
		wg.Add(1)
		go func(label, name string, fn func()) {
			defer wg.Done()
			fmt.Println("-->", label)
			run.Safe(name, fn)
		}(engine.label, engine.name, engine.fn)
	}
	wg.Wait()

	fmt.Println("All Engines executed successfully. Writing static assets...")

	publicDir := filepath.Join("..", "public")
	if err := assets.WriteGamesJSON(s, publicDir); err != nil {
		log.Fatalf("failed to write games.json: %v", err)
	}
	if err := assets.WriteBooksJSON(s, publicDir); err != nil {
		log.Fatalf("failed to write books.json: %v", err)
	}
	if err := assets.WriteSitemapXML(s, publicDir); err != nil {
		log.Fatalf("failed to write sitemap.xml: %v", err)
	}

	fmt.Println("Wrote public/games.json and public/books.json")
	fmt.Println("Wrote public/sitemap.xml")
	os.Exit(0)
}
