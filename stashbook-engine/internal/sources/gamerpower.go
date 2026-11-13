package sources

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"stashbook-engine/internal/model"
	"stashbook-engine/internal/store"
)

// GamerPowerDeal is a single giveaway from the GamerPower API.
type GamerPowerDeal struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Worth           string `json:"worth"`
	Thumbnail       string `json:"thumbnail"`
	Platforms       string `json:"platforms"`
	Type            string `json:"type"`
	Description     string `json:"description"`
	EndDate         string `json:"end_date"`
	OpenGiveawayUrl string `json:"open_giveaway_url"`
}

// FetchPrimeGames stores the Amazon Prime Gaming giveaways listed by
// GamerPower.
func FetchPrimeGames(s *store.Store) {
	url := "https://www.gamerpower.com/api/giveaways"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("GamerPower Network Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("GamerPower API Error. Status Code: %d\n", resp.StatusCode)
		return
	}

	var deals []GamerPowerDeal
	json.NewDecoder(resp.Body).Decode(&deals)

	var primeCount int

	for _, deal := range deals {
		titleLower := strings.ToLower(deal.Title)

		if strings.Contains(titleLower, "prime") || strings.Contains(titleLower, "amazon") {
			primeCount++

			cleanPrice := strings.ReplaceAll(deal.Worth, "$", "")
			originalPrice, _ := strconv.ParseFloat(cleanPrice, 64)

			if deal.Worth == "N/A" {
				originalPrice = 0.0
			}

			endDate := deal.EndDate
			if endDate == "N/A" {
				endDate = "Unknown"
			}

			seoKeywords := fmt.Sprintf("Free %s, %s free game, discount, PC deals", deal.Title, "Amazon Prime Gaming")
			s.SaveGame(model.GameResponse{
				ID:            "prime_" + strconv.Itoa(deal.ID),
				Title:         deal.Title,
				Platform:      "Amazon Prime Gaming",
				ImageURL:      deal.Thumbnail,
				OriginalPrice: originalPrice,
				Description:   deal.Description,
				EndDate:       endDate,
				GameURL:       deal.OpenGiveawayUrl,
				SeoKeywords:   seoKeywords,
			})
			fmt.Println("Saved Active (Prime Gaming):", deal.Title)
		}
	}

	fmt.Printf("Checked %d global offers. Found %d active on Amazon Prime Gaming.\n", len(deals), primeCount)
}

// FetchGlobalGames stores the non-Steam, non-Epic game giveaways listed
// by GamerPower (GOG, itch.io and other platforms).
func FetchGlobalGames(s *store.Store) {
	url := "https://www.gamerpower.com/api/giveaways?type=game"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != 200 {
		return
	}
	defer resp.Body.Close()

	var deals []GamerPowerDeal
	json.NewDecoder(resp.Body).Decode(&deals)

	var addedCount int

	for _, deal := range deals {
		platformLower := strings.ToLower(deal.Platforms)

		if strings.Contains(platformLower, "steam") ||
			strings.Contains(platformLower, "epic") ||
			strings.Contains(strings.ToLower(deal.Title), "prime") {
			continue
		}

		cleanPrice := strings.ReplaceAll(deal.Worth, "$", "")
		originalPrice, _ := strconv.ParseFloat(cleanPrice, 64)

		if deal.Worth == "N/A" {
			originalPrice = 0.0
		}

		displayPlatform := deal.Platforms
		if strings.Contains(platformLower, "itchio") || strings.Contains(platformLower, "itch.io") {
			displayPlatform = "itch.io"
		} else if strings.Contains(platformLower, "gog") {
			displayPlatform = "GOG"
		}

		endDate := deal.EndDate
		if endDate == "N/A" {
			endDate = "Unknown"
		}

		seoKeywords := fmt.Sprintf("Free %s, %s free game, discount, PC deals", deal.Title, displayPlatform)
		s.SaveGame(model.GameResponse{
			ID:            "gp_global_" + strconv.Itoa(deal.ID),
			Title:         deal.Title,
			Platform:      displayPlatform,
			ImageURL:      deal.Thumbnail,
			OriginalPrice: originalPrice,
			Description:   deal.Description,
			EndDate:       endDate,
			GameURL:       deal.OpenGiveawayUrl,
			SeoKeywords:   seoKeywords,
		})

		addedCount++
	}

	fmt.Printf("Engine 6: Added %d games from GOG, itch.io, and other platforms.\n", addedCount)
}
