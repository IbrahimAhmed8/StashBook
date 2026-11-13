// Package sources implements one aggregation client per storefront. Each
// client fetches that source's current free offers and stores them via
// the store package; all clients run concurrently from main.
package sources

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"stashbook-engine/internal/model"
	"stashbook-engine/internal/store"
)

// epicFreeGamesURL is Epic's public free-games-promotions feed. The
// storefront GraphQL endpoint (store.epicgames.com/graphql) is behind
// Cloudflare bot protection and answers server-side requests with a 403
// challenge page, so the feed below is used instead. It returns the same
// data.Catalog.searchStore.elements envelope, so the parsing below is
// unchanged.
const epicFreeGamesURL = "https://store-site-backend-static-ipv4.ak.epicgames.com/freeGamesPromotions?locale=en-US&country=EG&allowCountries=EG"

// EpicResponse is the Epic Games Store free-games feed response.
type EpicResponse struct {
	Data struct {
		Catalog struct {
			SearchStore struct {
				Elements []EpicGame `json:"elements"`
			} `json:"searchStore"`
		} `json:"Catalog"`
	} `json:"data"`
}

// EpicGame is a single element of the Epic free-games catalog.
type EpicGame struct {
	Title       string `json:"title"`
	ID          string `json:"id"`
	Description string `json:"description"`
	ProductSlug string `json:"productSlug"`
	UrlSlug     string `json:"urlSlug"`
	KeyImages   []struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"keyImages"`
	Price struct {
		TotalPrice struct {
			DiscountPrice int `json:"discountPrice"`
			OriginalPrice int `json:"originalPrice"`
		} `json:"totalPrice"`
	} `json:"price"`
}

// FetchEpicGames stores the Epic Games Store's current weekly free games.
func FetchEpicGames(s *store.Store) {
	games, err := fetchEpicFreeGames(http.DefaultClient, epicFreeGamesURL)
	if err != nil {
		log.Printf("[Epic] fetch failed: %v", err)
		return
	}

	for _, game := range games {
		var imageUrl string
		for _, img := range game.KeyImages {
			if img.Type == "OfferImageWide" || img.Type == "Thumbnail" {
				imageUrl = img.URL
				break
			}
		}

		slug := game.ProductSlug
		if slug == "" {
			slug = game.UrlSlug
		}
		gameUrl := "https://store.epicgames.com/p/" + slug

		desc := game.Description
		if desc == "" {
			desc = "Free game on the Epic Games Store."
		}

		seoKeywords := fmt.Sprintf("Free %s, %s free game, discount, PC deals", game.Title, "Epic Games")
		s.SaveGame(model.GameResponse{
			ID:            "epic_" + game.ID,
			Title:         game.Title,
			Platform:      "Epic Games",
			ImageURL:      imageUrl,
			OriginalPrice: float64(game.Price.TotalPrice.OriginalPrice) / 100.0,
			Description:   desc,
			EndDate:       "Limited Time (Epic Weekly)",
			GameURL:       gameUrl,
			SeoKeywords:   seoKeywords,
		})
		fmt.Println("Saved Active (Epic):", game.Title)
	}
}

// fetchEpicFreeGames fetches the free-games feed and returns the games
// that are currently free: discounted to 0 with a non-zero original price.
// Upcoming promotions and paid titles in the feed are filtered out.
func fetchEpicFreeGames(client *http.Client, url string) ([]EpicGame, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("epic free-games feed returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var epicData EpicResponse
	if err := json.Unmarshal(body, &epicData); err != nil {
		return nil, err
	}

	var free []EpicGame
	for _, game := range epicData.Data.Catalog.SearchStore.Elements {
		if game.Price.TotalPrice.DiscountPrice == 0 && game.Price.TotalPrice.OriginalPrice > 0 {
			free = append(free, game)
		}
	}
	return free, nil
}
