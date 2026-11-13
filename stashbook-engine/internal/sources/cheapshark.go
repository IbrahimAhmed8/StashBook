package sources

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"stashbook-engine/internal/model"
	"stashbook-engine/internal/store"
)

// CheapSharkDeal is a single deal from the CheapShark API.
type CheapSharkDeal struct {
	Title              string `json:"title"`
	DealID             string `json:"dealID"`
	StoreID            string `json:"storeID"`
	GameID             string `json:"gameID"`
	SalePrice          string `json:"salePrice"`
	NormalPrice        string `json:"normalPrice"`
	Thumb              string `json:"thumb"`
	Savings            string `json:"savings"`
	SteamRatingCount   string `json:"steamRatingCount"`
	SteamRatingPercent string `json:"steamRatingPercent"`
}

// FetchCheapSharkDeals stores the 100% off deals currently listed by
// CheapShark.
func FetchCheapSharkDeals(s *store.Store) {
	url := "https://www.cheapshark.com/api/1.0/deals?upperPrice=0"

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return
	}
	defer resp.Body.Close()

	var deals []CheapSharkDeal
	json.NewDecoder(resp.Body).Decode(&deals)

	storeNames := map[string]string{
		"1":  "Steam",
		"7":  "GOG",
		"11": "Humble Store",
		"25": "Epic Games",
	}

	for _, deal := range deals {
		salePrice, _ := strconv.ParseFloat(deal.SalePrice, 64)
		normalPrice, _ := strconv.ParseFloat(deal.NormalPrice, 64)

		if salePrice == 0 && normalPrice > 0 {
			storeName, exists := storeNames[deal.StoreID]
			if !exists {
				storeName = "Store ID: " + deal.StoreID
			}
			platformLabel := fmt.Sprintf("%s (Via CheapShark)", storeName)
			gameUrl := "https://www.cheapshark.com/redirect?dealID=" + deal.DealID
			desc := "This 100% off deal was found via CheapShark API."

			seoKeywords := fmt.Sprintf("Free %s, %s free game, discount, PC deals", deal.Title, platformLabel)
			s.SaveGame(model.GameResponse{
				ID:            "cs_" + deal.DealID,
				Title:         deal.Title,
				Platform:      platformLabel,
				ImageURL:      deal.Thumb,
				OriginalPrice: normalPrice,
				Description:   desc,
				EndDate:       "Limited Time",
				GameURL:       gameUrl,
				SeoKeywords:   seoKeywords,
			})
			fmt.Println("Saved Active (CheapShark):", deal.Title)
		}
	}
}
