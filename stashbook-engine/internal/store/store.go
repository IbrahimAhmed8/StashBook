// Package store owns the engine's SQLite database: schema setup, writes
// from the source clients, and reads for the static asset writers.
//
// The database is a working cache, not the source of truth: every run
// re-marks all games expired and re-inserts the offers that are still
// live. A missing database is recreated from scratch by Open, so a cold
// start (for example a cache miss in CI) still produces a working site;
// only accumulated history (expired offers, tracked prices) is lost.
package store

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"stashbook-engine/internal/model"
)

// Store wraps the engine database connection.
type Store struct {
	db *sql.DB
}

// Open creates the database at path if needed and ensures the schema
// exists.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	s.initSchema()
	return s, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// Exec runs a statement, logging (but not returning) errors. The engine
// treats writes as fire-and-forget: a failed insert must not abort the
// run.
func (s *Store) Exec(query string, args ...any) {
	_, err := s.db.Exec(query, args...)
	if err != nil {
		log.Printf("[DB] exec error: %v | query: %s", err, query)
	}
}

// MarkAllGamesExpired flags every stored game as expired; the source
// clients then re-activate the offers that are still live.
func (s *Store) MarkAllGamesExpired() {
	s.Exec("UPDATE free_games SET is_expired = 1")
}

// SaveGame inserts or replaces a game offer, marked as active.
func (s *Store) SaveGame(g model.GameResponse) {
	s.Exec(`INSERT OR REPLACE INTO free_games (id, title, platform, image_url, original_price, is_expired, description, end_date, game_url, seo_keywords) VALUES (?, ?, ?, ?, ?, 0, ?, ?, ?, ?)`,
		g.ID, g.Title, g.Platform, g.ImageURL, g.OriginalPrice, g.Description, g.EndDate, g.GameURL, g.SeoKeywords)
}

// SaveBook inserts or replaces a free book.
func (s *Store) SaveBook(b model.BookResponse) {
	s.Exec(`INSERT OR REPLACE INTO free_books (id, title, author, image_url, description, book_url, seo_keywords) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		b.ID, b.Title, b.Author, b.ImageURL, b.Description, b.BookURL, b.SeoKeywords)
}

// LastPrice returns the last recorded price for a Steam app, or
// sql.ErrNoRows when the app has never been tracked.
func (s *Store) LastPrice(appID string) (float64, error) {
	var lastPrice float64
	err := s.db.QueryRow("SELECT last_price FROM price_history WHERE app_id = ?", appID).Scan(&lastPrice)
	return lastPrice, err
}

// InsertPrice starts tracking a Steam app's price.
func (s *Store) InsertPrice(appID, title string, price float64) {
	s.Exec("INSERT INTO price_history (app_id, title, last_price) VALUES (?, ?, ?)", appID, title, price)
}

// UpdatePrice records a Steam app's latest price.
func (s *Store) UpdatePrice(appID string, price float64) {
	s.Exec("UPDATE price_history SET last_price = ? WHERE app_id = ?", price, appID)
}

// Games returns all stored games, one row per title, preferring the
// non-expired copy when a title exists more than once.
func (s *Store) Games() ([]model.GameResponse, error) {
	query := "SELECT id, title, platform, image_url, original_price, description, end_date, game_url, seo_keywords, MIN(is_expired) FROM free_games GROUP BY title"
	rows, err := s.db.Query(query)
	if err != nil {
		log.Printf("[DB] query error: %v | query: %s", err, query)
		return nil, err
	}
	defer rows.Close()

	var games []model.GameResponse
	for rows.Next() {
		var g model.GameResponse
		var isExpiredInt int
		if err := rows.Scan(&g.ID, &g.Title, &g.Platform, &g.ImageURL, &g.OriginalPrice, &g.Description, &g.EndDate, &g.GameURL, &g.SeoKeywords, &isExpiredInt); err != nil {
			log.Printf("[DB] scan error: %v", err)
			continue
		}
		g.IsExpired = isExpiredInt == 1
		games = append(games, g)
	}
	if err := rows.Err(); err != nil {
		log.Printf("[DB] rows error: %v", err)
		return nil, err
	}
	return games, nil
}

// Books returns all stored books.
func (s *Store) Books() ([]model.BookResponse, error) {
	query := "SELECT id, title, author, image_url, description, book_url, seo_keywords FROM free_books"
	rows, err := s.db.Query(query)
	if err != nil {
		log.Printf("[DB] query error: %v | query: %s", err, query)
		return nil, err
	}
	defer rows.Close()

	var books []model.BookResponse
	for rows.Next() {
		var b model.BookResponse
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.ImageURL, &b.Description, &b.BookURL, &b.SeoKeywords); err != nil {
			log.Printf("[DB] scan error: %v", err)
			continue
		}
		books = append(books, b)
	}
	if err := rows.Err(); err != nil {
		log.Printf("[DB] rows error: %v", err)
		return nil, err
	}
	return books, nil
}

// GameTitles returns every stored game title, for the sitemap.
func (s *Store) GameTitles() ([]string, error) {
	rows, err := s.db.Query("SELECT title FROM free_games GROUP BY title")
	if err != nil {
		return nil, fmt.Errorf("query games: %w", err)
	}
	defer rows.Close()

	var titles []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			log.Printf("[DB] scan error: %v", err)
			continue
		}
		titles = append(titles, title)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return titles, nil
}

func (s *Store) initSchema() {
	createGamesTable := `
    CREATE TABLE IF NOT EXISTS free_games (
        id TEXT PRIMARY KEY,
        title TEXT,
        platform TEXT,
        image_url TEXT,
        original_price REAL,
        is_expired INTEGER DEFAULT 0,
        description TEXT,
        end_date TEXT,
        game_url TEXT
    );`
	s.Exec(createGamesTable)

	s.Exec("ALTER TABLE free_games ADD COLUMN description TEXT DEFAULT ''")
	s.Exec("ALTER TABLE free_games ADD COLUMN end_date TEXT DEFAULT ''")
	s.Exec("ALTER TABLE free_games ADD COLUMN game_url TEXT DEFAULT ''")
	s.Exec("ALTER TABLE free_games ADD COLUMN seo_keywords TEXT DEFAULT ''")

	s.Exec("CREATE INDEX IF NOT EXISTS idx_free_games_title ON free_games(title)")
	s.Exec("CREATE INDEX IF NOT EXISTS idx_free_games_is_expired ON free_games(is_expired)")
	s.Exec("CREATE INDEX IF NOT EXISTS idx_free_games_title_expired ON free_games(title, is_expired)")

	createHistoryTable := `
    CREATE TABLE IF NOT EXISTS price_history (
        app_id TEXT PRIMARY KEY,
        title TEXT,
        last_price REAL
    );`
	s.Exec(createHistoryTable)

	createBooksTable := `
    CREATE TABLE IF NOT EXISTS free_books (
        id TEXT PRIMARY KEY,
        title TEXT,
        author TEXT,
        image_url TEXT,
        description TEXT,
        book_url TEXT
    );`
	s.Exec(createBooksTable)

	s.Exec("ALTER TABLE free_books ADD COLUMN seo_keywords TEXT DEFAULT ''")
}
