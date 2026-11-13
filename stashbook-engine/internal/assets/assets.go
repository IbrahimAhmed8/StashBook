// Package assets writes the static files the frontend consumes:
// games.json, books.json and sitemap.xml.
package assets

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"stashbook-engine/internal/model"
	"stashbook-engine/internal/sanitize"
	"stashbook-engine/internal/slug"
	"stashbook-engine/internal/store"
)

// WriteGamesJSON publishes the stored games as games.json into outDir.
func WriteGamesJSON(s *store.Store, outDir string) error {
	games, err := s.Games()
	if err != nil {
		return err
	}

	for i := range games {
		g := &games[i]
		if strings.TrimSpace(g.SeoKeywords) == "" {
			g.SeoKeywords = fmt.Sprintf("Free %s, %s free game, discount, PC deals", g.Title, g.Platform)
		}
		if g.Description == "" {
			g.Description = fmt.Sprintf("Grab %s for free or at a huge discount on %s.", g.Title, g.Platform)
		}
		sanitize.SanitizeGame(g)
	}

	if games == nil {
		games = []model.GameResponse{}
	}

	return writeJSONFile(filepath.Join(outDir, "games.json"), games)
}

// WriteBooksJSON publishes the stored books as books.json into outDir.
func WriteBooksJSON(s *store.Store, outDir string) error {
	books, err := s.Books()
	if err != nil {
		return err
	}

	for i := range books {
		b := &books[i]
		if b.SeoKeywords == "" {
			b.SeoKeywords = fmt.Sprintf("Free %s, %s free book, programming, tech books", b.Title, b.Author)
		}
		if b.Description == "" {
			b.Description = fmt.Sprintf("Download the official digital edition of %s by %s.", b.Title, b.Author)
		}
		sanitize.SanitizeBook(b)
	}

	if books == nil {
		books = []model.BookResponse{}
	}

	return writeJSONFile(filepath.Join(outDir, "books.json"), books)
}

// WriteSitemapXML publishes sitemap.xml into outDir: the static routes
// plus one entry per stored game.
func WriteSitemapXML(s *store.Store, outDir string) error {
	path := filepath.Join(outDir, "sitemap.xml")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create sitemap: %w", err)
	}
	defer file.Close()

	fmt.Fprintln(file, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintln(file, `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

	staticRoutes := []struct {
		loc        string
		changefreq string
	}{
		{"https://stashbook.online/", "daily"},
		{"https://stashbook.online/books", "weekly"},
		{"https://stashbook.online/iconic", "weekly"},
		{"https://stashbook.online/about", "weekly"},
		{"https://stashbook.online/contact", "weekly"},
		{"https://stashbook.online/privacy", "weekly"},
	}

	for _, route := range staticRoutes {
		fmt.Fprintf(file, "  <url>\n    <loc>%s</loc>\n    <changefreq>%s</changefreq>\n  </url>\n", route.loc, route.changefreq)
	}

	titles, err := s.GameTitles()
	if err != nil {
		return err
	}

	for _, title := range titles {
		slug := slug.Slugify(title)
		if slug != "" {
			fmt.Fprintf(file, "  <url><loc>https://stashbook.online/game/%s</loc><changefreq>weekly</changefreq></url>\n", slug)
		}
	}

	fmt.Fprintln(file, "</urlset>")
	return nil
}

func writeJSONFile(path string, data any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}

	if err := os.WriteFile(path, payload, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
