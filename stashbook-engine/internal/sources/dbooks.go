package sources

import (
	"encoding/json"
	"fmt"
	"net/http"

	"stashbook-engine/internal/model"
	"stashbook-engine/internal/store"
)

// DBooksResponse is the dBooks "recent books" API response.
type DBooksResponse struct {
	Status string `json:"status"`
	Books  []struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
		Authors  string `json:"authors"`
		Image    string `json:"image"`
		URL      string `json:"url"`
	} `json:"books"`
}

// FetchTechBooks stores the latest free technical books from dBooks.
func FetchTechBooks(s *store.Store) {
	fmt.Println("Scanning for Free Technical Books...")

	resp, err := http.Get("https://www.dbooks.org/api/recent")
	if err != nil {
		fmt.Println("DBooks API Error:", err)
		return
	}
	defer resp.Body.Close()

	var result DBooksResponse
	json.NewDecoder(resp.Body).Decode(&result)

	var count int
	for _, b := range result.Books {
		desc := b.Subtitle
		if desc == "" {
			desc = "Official digital edition by " + b.Authors
		}

		seoKeywords := fmt.Sprintf("Free %s, %s free book, programming, tech books", b.Title, b.Authors)
		s.SaveBook(model.BookResponse{
			ID:          "dbook_" + b.ID,
			Title:       b.Title,
			Author:      b.Authors,
			ImageURL:    b.Image,
			Description: desc,
			BookURL:     b.URL,
			SeoKeywords: seoKeywords,
		})
		count++
	}

	fmt.Printf("Engine 8: Tech Vault Hunt finished. Secured %d Free Books.\n", count)
}
