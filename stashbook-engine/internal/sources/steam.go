package sources

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"stashbook-engine/internal/model"
	"stashbook-engine/internal/store"
)

// FetchSteamGames scans Steam's free specials search and stores the
// games that are currently 100% off.
func FetchSteamGames(s *store.Store) {
	searchURL := "https://store.steampowered.com/search/?sort_by=Price_ASC&maxprice=free&specials=1"
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	re := regexp.MustCompile(`data-ds-appid=["'\\]*(\d+)["'\\]*`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	uniqueIDs := make(map[string]bool)
	var appIDs []string
	for _, match := range matches {
		id := match[1]
		if !uniqueIDs[id] {
			uniqueIDs[id] = true
			appIDs = append(appIDs, id)
		}
	}

	for _, appID := range appIDs {
		detailsURL := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%s&cc=us&filters=basic,price_overview", appID)
		detReq, _ := http.NewRequest("GET", detailsURL, nil)
		detReq.Header.Set("User-Agent", "Mozilla/5.0")
		detResp, err := client.Do(detReq)
		if err != nil {
			continue
		}

		detBody, _ := io.ReadAll(detResp.Body)
		detResp.Body.Close()

		var result map[string]interface{}
		json.Unmarshal(detBody, &result)

		appData, ok := result[appID].(map[string]interface{})
		if !ok {
			continue
		}
		if success, _ := appData["success"].(bool); !success {
			continue
		}
		data, ok := appData["data"].(map[string]interface{})
		if !ok {
			continue
		}

		priceOverview, ok := data["price_overview"].(map[string]interface{})
		if ok {
			discountPercent, _ := priceOverview["discount_percent"].(float64)
			initialPriceCents, _ := priceOverview["initial"].(float64)

			if discountPercent == 100 && initialPriceCents > 0 {
				title := data["name"].(string)
				headerImage := data["header_image"].(string)
				originalPrice := initialPriceCents / 100.0

				desc := "Free game on Steam."
				if shortDesc, ok := data["short_description"].(string); ok {
					desc = shortDesc
				}
				gameUrl := "https://store.steampowered.com/app/" + appID

				seoKeywords := fmt.Sprintf("Free %s, %s free game, discount, PC deals", title, "Steam")
				s.SaveGame(model.GameResponse{
					ID:            "steam_live_" + appID,
					Title:         title,
					Platform:      "Steam",
					ImageURL:      headerImage,
					OriginalPrice: originalPrice,
					Description:   desc,
					EndDate:       "Limited Time",
					GameURL:       gameUrl,
					SeoKeywords:   seoKeywords,
				})
				fmt.Println("💾 Saved Active (Steam Live):", title)
			}
		}
	}
}
