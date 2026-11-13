package sources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"stashbook-engine/internal/model"
	"stashbook-engine/internal/run"
	"stashbook-engine/internal/store"
)

// FetchLegendaryDeals hunts for deep discounts (50% or more) on acclaimed
// titles: Steam top sellers, plus Epic deals via CheapShark.
func FetchLegendaryDeals(s *store.Store) {
	var caughtCount int

	fmt.Println("Scanning Steam Global Top Sellers...")

	steamURL := "https://store.steampowered.com/search/results/?query&start=0&count=100&dynamic_data=&sort_by=Top_Sellers&filter=globaltopsellers&specials=1&infinite=1"

	req, _ := http.NewRequest("GET", steamURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err == nil {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var searchResult map[string]interface{}
		json.Unmarshal(body, &searchResult)

		if htmlResult, ok := searchResult["results_html"].(string); ok {
			re := regexp.MustCompile(`data-ds-appid="(\d+)"`)
			matches := re.FindAllStringSubmatch(htmlResult, -1)

			appIDs := make(map[string]bool)
			for _, match := range matches {
				if len(match) > 1 {
					appIDs[match[1]] = true
				}
			}

			for appID := range appIDs {
				run.SafeIterate("fetchOfficialLegendaryDeals/steam", func() {
					detailsURL := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%s&cc=us&filters=basic,price_overview", appID)
					detReq, _ := http.NewRequest("GET", detailsURL, nil)
					detReq.Header.Set("User-Agent", "Mozilla/5.0")
					detResp, err := client.Do(detReq)
					if err != nil {
						return
					}

					detBody, _ := io.ReadAll(detResp.Body)
					detResp.Body.Close()

					var result map[string]interface{}
					json.Unmarshal(detBody, &result)

					appData, ok := result[appID].(map[string]interface{})
					if !ok {
						return
					}
					if success, _ := appData["success"].(bool); !success {
						return
					}
					data, ok := appData["data"].(map[string]interface{})
					if !ok {
						return
					}

					priceOverview, ok := data["price_overview"].(map[string]interface{})
					if ok {
						discountPercent, _ := priceOverview["discount_percent"].(float64)
						finalPriceCents, _ := priceOverview["final"].(float64)

						if discountPercent >= 50 {
							title := data["name"].(string)
							headerImage := data["header_image"].(string)
							finalPrice := finalPriceCents / 100.0
							gameUrl := "https://store.steampowered.com/app/" + appID

							shortDesc := getSteamDescription(appID)
							desc := fmt.Sprintf("[ %.0f%% OFF ] %s", discountPercent, shortDesc)

							seoKeywords := fmt.Sprintf("Free %s, %s free game, discount, PC deals", title, "Steam (Iconic)")
							s.SaveGame(model.GameResponse{
								ID:            "legendary_steam_" + appID,
								Title:         title,
								Platform:      "Steam (Iconic)",
								ImageURL:      headerImage,
								OriginalPrice: finalPrice,
								Description:   desc,
								EndDate:       "Unknown",
								GameURL:       gameUrl,
								SeoKeywords:   seoKeywords,
							})
							fmt.Printf("STEAM ICONIC SECURED: %s (%.0f%% OFF)\n", title, discountPercent)
							caughtCount++
						}
					}
					time.Sleep(200 * time.Millisecond)
				})
			}
		}
	}

	fmt.Println("Scanning Epic Games Iconic Deals...")
	epicURL := "https://www.cheapshark.com/api/1.0/deals?storeID=25&sortBy=Savings&pageNumber=0&pageSize=100"
	epicResp, epicErr := http.Get(epicURL)
	if epicErr == nil {
		var deals []CheapSharkDeal
		json.NewDecoder(epicResp.Body).Decode(&deals)
		epicResp.Body.Close()

		for _, deal := range deals {
			savings, _ := strconv.ParseFloat(deal.Savings, 64)
			ratingPercent, _ := strconv.Atoi(deal.SteamRatingPercent)

			if savings >= 50.0 && ratingPercent >= 70 {
				salePrice, _ := strconv.ParseFloat(deal.SalePrice, 64)
				gameUrl := "https://www.cheapshark.com/redirect?dealID=" + deal.DealID

				realDesc := getEpicDescription(deal.Title)
				if realDesc == "" {
					realDesc = "Official Epic Games Deal."
				}
				desc := fmt.Sprintf("[ %.0f%% OFF ] %s", savings, realDesc)

				seoKeywords := fmt.Sprintf("Free %s, %s free game, discount, PC deals", deal.Title, "Epic Games (Iconic)")
				s.SaveGame(model.GameResponse{
					ID:            "legendary_epic_" + deal.DealID,
					Title:         deal.Title,
					Platform:      "Epic Games (Iconic)",
					ImageURL:      deal.Thumb,
					OriginalPrice: salePrice,
					Description:   desc,
					EndDate:       "Unknown",
					GameURL:       gameUrl,
					SeoKeywords:   seoKeywords,
				})
				fmt.Printf("EPIC ICONIC SECURED: %s (%.0f%% OFF)\n", deal.Title, savings)
				caughtCount++
			}
		}
	}

	fmt.Printf("Engine 7: Iconic Hunt finished. Secured %d Drops.\n", caughtCount)
}

func getSteamDescription(appID string) string {
	if appID == "" || appID == "0" {
		return ""
	}
	detailsURL := "https://store.steampowered.com/api/appdetails?appids=" + appID + "&cc=us&filters=basic"
	req, _ := http.NewRequest("GET", detailsURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{Timeout: 5 * time.Second}
	if resp, err := client.Do(req); err == nil {
		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		resp.Body.Close()

		if appData, ok := res[appID].(map[string]interface{}); ok {
			if success, _ := appData["success"].(bool); success {
				if data, ok := appData["data"].(map[string]interface{}); ok {
					if shortDesc, ok := data["short_description"].(string); ok {
						return shortDesc
					}
				}
			}
		}
	}
	return ""
}

func getEpicDescription(title string) string {
	graphqlQuery := map[string]interface{}{
		"query": `query searchStoreQuery($keywords: String!) { Catalog { searchStore(keywords: $keywords, country: "US", locale: "en-US", count: 1) { elements { description } } } }`,
		"variables": map[string]interface{}{
			"keywords": title,
		},
	}

	jsonData, _ := json.Marshal(graphqlQuery)
	req, _ := http.NewRequest("POST", "https://graphql.epicgames.com/graphql", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return ""
	}
	catalog, ok := data["Catalog"].(map[string]interface{})
	if !ok {
		return ""
	}
	searchStore, ok := catalog["searchStore"].(map[string]interface{})
	if !ok {
		return ""
	}
	elements, ok := searchStore["elements"].([]interface{})
	if !ok || len(elements) == 0 {
		return ""
	}
	firstElement, ok := elements[0].(map[string]interface{})
	if !ok {
		return ""
	}
	description, ok := firstElement["description"].(string)
	if !ok {
		return ""
	}

	return description
}
