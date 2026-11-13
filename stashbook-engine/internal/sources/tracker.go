package sources

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"stashbook-engine/internal/model"
	"stashbook-engine/internal/run"
	"stashbook-engine/internal/store"
)

// TrackSteamHistory watches Steam's top sellers for games whose price
// drops to zero between runs, catching free-to-play transitions. It is
// the only client that depends on history carried over from previous
// runs (the price_history table); on a cold start it only records
// current prices and catches transitions from the next run on.
func TrackSteamHistory(s *store.Store) {
	url := "https://store.steampowered.com/api/featuredcategories/?cc=us"

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	topSellers, ok := result["top_sellers"].(map[string]interface{})
	if !ok {
		return
	}
	items, ok := topSellers["items"].([]interface{})
	if !ok {
		return
	}

	for _, itemInterface := range items {
		run.SafeIterate("selfTrackSteamHistory", func() {
			item := itemInterface.(map[string]interface{})
			appID := fmt.Sprintf("%v", item["id"])
			title := fmt.Sprintf("%v", item["name"])

			var currentPrice float64 = 0
			originalPriceInt, hasPrice := item["original_price"].(float64)
			if hasPrice {
				currentPrice = originalPriceInt / 100.0
			}

			lastPrice, err := s.LastPrice(appID)

			if err == sql.ErrNoRows {
				s.InsertPrice(appID, title, currentPrice)
			} else if err != nil {
				log.Printf("[DB] query error: %v | app_id: %s", err, appID)
			} else {
				if lastPrice > 0 && currentPrice == 0 {
					headerImage := fmt.Sprintf("https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/%s/header.jpg", appID)
					gameUrl := "https://store.steampowered.com/app/" + appID
					desc := "F2P Transition Tracked Nativey by StashBook Engine."

					seoKeywords := fmt.Sprintf("Free %s, %s free game, discount, PC deals", title, "Steam (Tracked F2P)")
					s.SaveGame(model.GameResponse{
						ID:            "tracker_" + appID,
						Title:         title,
						Platform:      "Steam (Tracked F2P)",
						ImageURL:      headerImage,
						OriginalPrice: lastPrice,
						Description:   desc,
						EndDate:       "Permanent",
						GameURL:       gameUrl,
						SeoKeywords:   seoKeywords,
					})
					fmt.Println("F2P TRANSITION CAUGHT BY TRACKER:", title)
				}
				s.UpdatePrice(appID, currentPrice)
			}
		})
	}
	fmt.Println("Internal Price History Updated.")
}
