package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type EpicResponse struct {
	Data struct {
		Catalog struct {
			SearchStore struct {
				Elements []EpicGame `json:"elements"`
			} `json:"searchStore"`
		} `json:"Catalog"`
	} `json:"data"`
}

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

type GameResponse struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Platform      string  `json:"platform"`
	ImageURL      string  `json:"image"`
	OriginalPrice float64 `json:"original_price"`
	IsExpired     bool    `json:"is_expired"`
	Description   string  `json:"description"`
	EndDate       string  `json:"end_date"`
	GameURL       string  `json:"game_url"`
	SeoKeywords   string  `json:"seo_keywords"`
}

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

type BookResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	ImageURL    string `json:"image"`
	Description string `json:"description"`
	BookURL     string `json:"book_url"`
	SeoKeywords string `json:"seo_keywords"`
}

var db *sql.DB

func execDB(query string, args ...interface{}) {
	_, err := db.Exec(query, args...)
	if err != nil {
		log.Printf("[DB] exec error: %v | query: %s", err, query)
	}
}

func safeRunEngine(name string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ENGINE PANIC] %s recovered: %v", name, r)
		}
	}()
	fn()
}

func safeIterate(name string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ENGINE ITERATION] %s recovered: %v", name, r)
		}
	}()
	fn()
}

func main() {
	initDB()
	defer db.Close()

	execDB("UPDATE free_games SET is_expired = 1")

	fmt.Println("Starting Aggregation Engines...")

	engines := []struct {
		label string
		name  string
		fn    func()
	}{
		{"Engine 1: Epic Games", "Epic Games", fetchAndStoreEpicGames},
		{"Engine 2: Steam Live Deals", "Steam Live Deals", fetchAndStoreSteamGames},
		{"Engine 3: CheapShark Aggregator", "CheapShark Aggregator", fetchCheapSharkDeals},
		{"Engine 4: StashBook Self-Tracker", "StashBook Self-Tracker", selfTrackSteamHistory},
		{"Engine 5: Amazon Prime Gaming (GamerPower API)", "Amazon Prime Gaming", fetchGamerPowerPrimeGames},
		{"Engine 6: Global Indie & Consoles", "Global Indie & Consoles", fetchGamerPowerGlobalGames},
		{"Engine 7: Iconic Deals Hunter", "Iconic Deals Hunter", fetchOfficialLegendaryDeals},
		{"Engine 8: Tech Vault (Digital Books Tracker)", "Tech Vault", fetchFreeTechBooks},
	}

	var wg sync.WaitGroup
	for _, engine := range engines {
		wg.Add(1)
		go func(label, name string, fn func()) {
			defer wg.Done()
			fmt.Println("-->", label)
			safeRunEngine(name, fn)
		}(engine.label, engine.name, engine.fn)
	}
	wg.Wait()

	fmt.Println("All Engines executed successfully. Writing static assets...")

	if err := writeGamesJSON(); err != nil {
		log.Fatalf("failed to write games.json: %v", err)
	}
	if err := writeBooksJSON(); err != nil {
		log.Fatalf("failed to write books.json: %v", err)
	}
	if err := writeSitemapXML(); err != nil {
		log.Fatalf("failed to write sitemap.xml: %v", err)
	}

	fmt.Println("Wrote public/games.json and public/books.json")
	fmt.Println("Wrote public/sitemap.xml")
	os.Exit(0)
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./stashbook.db")
	if err != nil {
		log.Fatal(err)
	}

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
	execDB(createGamesTable)

	execDB("ALTER TABLE free_games ADD COLUMN description TEXT DEFAULT ''")
	execDB("ALTER TABLE free_games ADD COLUMN end_date TEXT DEFAULT ''")
	execDB("ALTER TABLE free_games ADD COLUMN game_url TEXT DEFAULT ''")
	execDB("ALTER TABLE free_games ADD COLUMN seo_keywords TEXT DEFAULT ''")

	execDB("CREATE INDEX IF NOT EXISTS idx_free_games_title ON free_games(title)")
	execDB("CREATE INDEX IF NOT EXISTS idx_free_games_is_expired ON free_games(is_expired)")
	execDB("CREATE INDEX IF NOT EXISTS idx_free_games_title_expired ON free_games(title, is_expired)")

	createHistoryTable := `
    CREATE TABLE IF NOT EXISTS price_history (
        app_id TEXT PRIMARY KEY,
        title TEXT,
        last_price REAL
    );`
	execDB(createHistoryTable)

	createBooksTable := `
    CREATE TABLE IF NOT EXISTS free_books (
        id TEXT PRIMARY KEY,
        title TEXT,
        author TEXT,
        image_url TEXT,
        description TEXT,
        book_url TEXT
    );`
	execDB(createBooksTable)

	execDB("ALTER TABLE free_books ADD COLUMN seo_keywords TEXT DEFAULT ''")
}

func fetchAndStoreEpicGames() {
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

			insertQuery := `INSERT OR REPLACE INTO free_games (id, title, platform, image_url, original_price, is_expired, description, end_date, game_url, seo_keywords) VALUES (?, ?, ?, ?, ?, 0, ?, ?, ?, ?)`
			seoKeywords := fmt.Sprintf("Free %s, %s free game, discount, PC deals", game.Title, "Epic Games")
			execDB(insertQuery, "epic_"+game.ID, game.Title, "Epic Games", imageUrl, float64(game.Price.TotalPrice.OriginalPrice)/100.0, desc, "Limited Time (Epic Weekly)", gameUrl, seoKeywords)
			fmt.Println("Saved Active (Epic):", game.Title)
		}
	}
}

func fetchAndStoreSteamGames() {
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

				insertQuery := `INSERT OR REPLACE INTO free_games (id, title, platform, image_url, original_price, is_expired, description, end_date, game_url, seo_keywords) VALUES (?, ?, ?, ?, ?, 0, ?, ?, ?, ?)`
				seoKeywords := fmt.Sprintf("Free %s, %s free game, discount, PC deals", title, "Steam")
				execDB(insertQuery, "steam_live_"+appID, title, "Steam", headerImage, originalPrice, desc, "Limited Time", gameUrl, seoKeywords)
				fmt.Println("💾 Saved Active (Steam Live):", title)
			}
		}
	}
}

func fetchCheapSharkDeals() {
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

			insertQuery := `INSERT OR REPLACE INTO free_games (id, title, platform, image_url, original_price, is_expired, description, end_date, game_url, seo_keywords) VALUES (?, ?, ?, ?, ?, 0, ?, ?, ?, ?)`
			seoKeywords := fmt.Sprintf("Free %s, %s free game, discount, PC deals", deal.Title, platformLabel)
			execDB(insertQuery, "cs_"+deal.DealID, deal.Title, platformLabel, deal.Thumb, normalPrice, desc, "Limited Time", gameUrl, seoKeywords)
			fmt.Println("Saved Active (CheapShark):", deal.Title)
		}
	}
}

func selfTrackSteamHistory() {
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
		safeIterate("selfTrackSteamHistory", func() {
			item := itemInterface.(map[string]interface{})
			appID := fmt.Sprintf("%v", item["id"])
			title := fmt.Sprintf("%v", item["name"])

			var currentPrice float64 = 0
			originalPriceInt, hasPrice := item["original_price"].(float64)
			if hasPrice {
				currentPrice = originalPriceInt / 100.0
			}

			var lastPrice float64
			err := db.QueryRow("SELECT last_price FROM price_history WHERE app_id = ?", appID).Scan(&lastPrice)

			if err == sql.ErrNoRows {
				execDB("INSERT INTO price_history (app_id, title, last_price) VALUES (?, ?, ?)", appID, title, currentPrice)
			} else if err != nil {
				log.Printf("[DB] query error: %v | app_id: %s", err, appID)
			} else {
				if lastPrice > 0 && currentPrice == 0 {
					headerImage := fmt.Sprintf("https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/%s/header.jpg", appID)
					gameUrl := "https://store.steampowered.com/app/" + appID
					desc := "F2P Transition Tracked Nativey by StashBook Engine."

					insertQuery := `INSERT OR REPLACE INTO free_games (id, title, platform, image_url, original_price, is_expired, description, end_date, game_url, seo_keywords) VALUES (?, ?, ?, ?, ?, 0, ?, ?, ?, ?)`
					seoKeywords := fmt.Sprintf("Free %s, %s free game, discount, PC deals", title, "Steam (Tracked F2P)")
					execDB(insertQuery, "tracker_"+appID, title, "Steam (Tracked F2P)", headerImage, lastPrice, desc, "Permanent", gameUrl, seoKeywords)
					fmt.Println("F2P TRANSITION CAUGHT BY TRACKER:", title)
				}
				execDB("UPDATE price_history SET last_price = ? WHERE app_id = ?", currentPrice, appID)
			}
		})
	}
	fmt.Println("Internal Price History Updated.")
}

func fetchGamerPowerPrimeGames() {
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

			insertQuery := `INSERT OR REPLACE INTO free_games (id, title, platform, image_url, original_price, is_expired, description, end_date, game_url, seo_keywords) VALUES (?, ?, ?, ?, ?, 0, ?, ?, ?, ?)`
			seoKeywords := fmt.Sprintf("Free %s, %s free game, discount, PC deals", deal.Title, "Amazon Prime Gaming")
			execDB(insertQuery, "prime_"+strconv.Itoa(deal.ID), deal.Title, "Amazon Prime Gaming", deal.Thumbnail, originalPrice, deal.Description, endDate, deal.OpenGiveawayUrl, seoKeywords)
			fmt.Println("Saved Active (Prime Gaming):", deal.Title)
		}
	}

	fmt.Printf("Checked %d global offers. Found %d active on Amazon Prime Gaming.\n", len(deals), primeCount)
}

func fetchGamerPowerGlobalGames() {
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

		insertQuery := `INSERT OR REPLACE INTO free_games (id, title, platform, image_url, original_price, is_expired, description, end_date, game_url, seo_keywords) VALUES (?, ?, ?, ?, ?, 0, ?, ?, ?, ?)`
		seoKeywords := fmt.Sprintf("Free %s, %s free game, discount, PC deals", deal.Title, displayPlatform)
		execDB(insertQuery, "gp_global_"+strconv.Itoa(deal.ID), deal.Title, displayPlatform, deal.Thumbnail, originalPrice, deal.Description, endDate, deal.OpenGiveawayUrl, seoKeywords)

		addedCount++
	}

	fmt.Printf("Engine 6: Added %d games from GOG, itch.io, and other platforms.\n", addedCount)
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

func fetchOfficialLegendaryDeals() {
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
				safeIterate("fetchOfficialLegendaryDeals/steam", func() {
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

							insertQuery := `INSERT OR REPLACE INTO free_games (id, title, platform, image_url, original_price, is_expired, description, end_date, game_url, seo_keywords) VALUES (?, ?, ?, ?, ?, 0, ?, ?, ?, ?)`
							seoKeywords := fmt.Sprintf("Free %s, %s free game, discount, PC deals", title, "Steam (Iconic)")
							execDB(insertQuery, "legendary_steam_"+appID, title, "Steam (Iconic)", headerImage, finalPrice, desc, "Unknown", gameUrl, seoKeywords)
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

				insertQuery := `INSERT OR REPLACE INTO free_games (id, title, platform, image_url, original_price, is_expired, description, end_date, game_url, seo_keywords) VALUES (?, ?, ?, ?, ?, 0, ?, ?, ?, ?)`
				seoKeywords := fmt.Sprintf("Free %s, %s free game, discount, PC deals", deal.Title, "Epic Games (Iconic)")
				execDB(insertQuery, "legendary_epic_"+deal.DealID, deal.Title, "Epic Games (Iconic)", deal.Thumb, salePrice, desc, "Unknown", gameUrl, seoKeywords)
				fmt.Printf("EPIC ICONIC SECURED: %s (%.0f%% OFF)\n", deal.Title, savings)
				caughtCount++
			}
		}
	}

	fmt.Printf("Engine 7: Iconic Hunt finished. Secured %d Drops.\n", caughtCount)
}

func fetchFreeTechBooks() {
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

		insertQuery := `INSERT OR REPLACE INTO free_books (id, title, author, image_url, description, book_url, seo_keywords) VALUES (?, ?, ?, ?, ?, ?, ?)`
		seoKeywords := fmt.Sprintf("Free %s, %s free book, programming, tech books", b.Title, b.Authors)
		execDB(insertQuery, "dbook_"+b.ID, b.Title, b.Authors, b.Image, desc, b.URL, seoKeywords)
		count++
	}

	fmt.Printf("Engine 8: Tech Vault Hunt finished. Secured %d Free Books.\n", count)
}

func publicAssetPath(name string) string {
	return filepath.Join("..", "public", name)
}

func writeJSONFile(path string, data interface{}) error {
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

func writeGamesJSON() error {
	query := "SELECT id, title, platform, image_url, original_price, description, end_date, game_url, seo_keywords, MIN(is_expired) FROM free_games GROUP BY title"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("[DB] query error: %v | query: %s", err, query)
		return err
	}
	defer rows.Close()

	var games []GameResponse
	for rows.Next() {
		var g GameResponse
		var isExpiredInt int

		if err := rows.Scan(&g.ID, &g.Title, &g.Platform, &g.ImageURL, &g.OriginalPrice, &g.Description, &g.EndDate, &g.GameURL, &g.SeoKeywords, &isExpiredInt); err != nil {
			log.Printf("[DB] scan error: %v", err)
			continue
		}
		g.IsExpired = isExpiredInt == 1
		if strings.TrimSpace(g.SeoKeywords) == "" {
			g.SeoKeywords = fmt.Sprintf("Free %s, %s free game, discount, PC deals", g.Title, g.Platform)
		}
		if g.Description == "" {
			g.Description = fmt.Sprintf("Grab %s for free or at a huge discount on %s.", g.Title, g.Platform)
		}
		sanitizeGameResponse(&g)
		games = append(games, g)
	}
	if err := rows.Err(); err != nil {
		log.Printf("[DB] rows error: %v", err)
		return err
	}

	if games == nil {
		games = []GameResponse{}
	}

	return writeJSONFile(publicAssetPath("games.json"), games)
}

func writeBooksJSON() error {
	query := "SELECT id, title, author, image_url, description, book_url, seo_keywords FROM free_books"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("[DB] query error: %v | query: %s", err, query)
		return err
	}
	defer rows.Close()

	var books []BookResponse
	for rows.Next() {
		var b BookResponse
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.ImageURL, &b.Description, &b.BookURL, &b.SeoKeywords); err != nil {
			log.Printf("[DB] scan error: %v", err)
			continue
		}
		if b.SeoKeywords == "" {
			b.SeoKeywords = fmt.Sprintf("Free %s, %s free book, programming, tech books", b.Title, b.Author)
		}
		if b.Description == "" {
			b.Description = fmt.Sprintf("Download the official digital edition of %s by %s.", b.Title, b.Author)
		}
		sanitizeBookResponse(&b)
		books = append(books, b)
	}
	if err := rows.Err(); err != nil {
		log.Printf("[DB] rows error: %v", err)
		return err
	}

	if books == nil {
		books = []BookResponse{}
	}

	return writeJSONFile(publicAssetPath("books.json"), books)
}

func sanitizeText(value string) string {
	return strings.TrimSpace(html.UnescapeString(value))
}

func sanitizeGameResponse(game *GameResponse) {
	game.Title = sanitizeText(game.Title)
	game.Platform = sanitizeText(game.Platform)
	game.ImageURL = sanitizeText(game.ImageURL)
	game.Description = sanitizeText(game.Description)
	game.EndDate = sanitizeText(game.EndDate)
	game.GameURL = sanitizeText(game.GameURL)
	game.SeoKeywords = sanitizeText(game.SeoKeywords)
}

func sanitizeBookResponse(book *BookResponse) {
	book.Title = sanitizeText(book.Title)
	book.Author = sanitizeText(book.Author)
	book.ImageURL = sanitizeText(book.ImageURL)
	book.Description = sanitizeText(book.Description)
	book.BookURL = sanitizeText(book.BookURL)
	book.SeoKeywords = sanitizeText(book.SeoKeywords)
}

func writeSitemapXML() error {
	path := publicAssetPath("sitemap.xml")
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
	}

	for _, route := range staticRoutes {
		fmt.Fprintf(file, "  <url>\n    <loc>%s</loc>\n    <changefreq>%s</changefreq>\n  </url>\n", route.loc, route.changefreq)
	}

	rows, err := db.Query("SELECT title FROM free_games GROUP BY title")
	if err != nil {
		return fmt.Errorf("query games: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			log.Printf("[DB] scan error: %v", err)
			continue
		}
		slug := slugify(title)
		if slug != "" {
			fmt.Fprintf(file, "  <url><loc>https://stashbook.online/game/%s</loc><changefreq>weekly</changefreq></url>\n", slug)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error: %w", err)
	}

	fmt.Fprintln(file, "</urlset>")
	return nil
}

var (
	slugNonAlphaNum = regexp.MustCompile(`[^a-z0-9\s-]`)
	slugSpaces      = regexp.MustCompile(`[\s_]+`)
	slugHyphens     = regexp.MustCompile(`-+`)
)

func slugify(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))
	text = slugNonAlphaNum.ReplaceAllString(text, "")
	text = slugSpaces.ReplaceAllString(text, "-")
	text = slugHyphens.ReplaceAllString(text, "-")
	return strings.Trim(text, "-")
}
