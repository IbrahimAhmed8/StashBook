// Package sources implements one aggregation client per storefront. Each
// client fetches that source's current free offers and stores them via
// the store package; all clients run concurrently from main.
package sources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"stashbook-engine/internal/model"
	"stashbook-engine/internal/store"
)

// EpicResponse is the Epic Games Store GraphQL search response.
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
	url := "https://store.epicgames.com/graphql"
	payload := map[string]interface{}{
		"operationName": "searchStoreQuery",
		"variables": map[string]interface{}{
			"category": "freegames", "count": 50, "country": "EG", "locale": "en-US", "withPrice": true, "withPromotions": true,
		},
		"query": `query searchStoreQuery($category: String, $count: Int, $country: String!, $locale: String, $withPrice: Boolean = false, $withPromotions: Boolean = false) { Catalog { searchStore(category: $category, count: $count, country: $country, locale: $locale) { elements { title id description productSlug urlSlug keyImages { type url } price(country: $country) @include(if: $withPrice) { totalPrice { discountPrice originalPrice } } promotions(category: $category) @include(if: $withPromotions) { promotionalOffers { promotionalOffers { discountSetting { discountPercentage } } } } } } } }`,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var epicData EpicResponse
	json.Unmarshal(body, &epicData)

	for _, game := range epicData.Data.Catalog.SearchStore.Elements {
		if game.Price.TotalPrice.DiscountPrice == 0 && game.Price.TotalPrice.OriginalPrice > 0 {
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
}
