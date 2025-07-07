package googlesearch

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/lib/pq"
	"googlemaps.github.io/maps"
)

// findEmail searches the given HTML body for an email address.
func findEmail(body string) string {
	re := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	return re.FindString(body)
}

// detectSiteType analyzes the HTML body to identify the underlying website technology.
func detectSiteType(body string) string {
	if strings.Contains(body, "wp-content") || strings.Contains(body, "WordPress") {
		return "WordPress"
	}
	if strings.Contains(body, "sites/default/files") || strings.Contains(body, "Drupal") {
		return "Drupal"
	}
	if strings.Contains(body, `id="root"`) {
		return "React"
	}
	if strings.Contains(body, `id="app"`) {
		return "Vue.js"
	}
	if strings.Contains(body, "/_next/") {
		return "Next.js"
	}
	if strings.Contains(body, "/_nuxt/") {
		return "Nuxt.js"
	}
	if strings.Contains(body, "cdn.shopify.com") {
		return "Shopify"
	}
	if strings.Contains(body, "wix.com") {
		return "Wix"
	}
	if strings.Contains(body, "joomla") {
		return "Joomla"
	}
	if strings.Contains(body, "squarespace.com") {
		return "Squarespace"
	}
	if strings.Contains(body, "ghost/content/") {
		return "Ghost"
	}
	return "Unknown"
}

// getPageBody fetches the HTML content of a given URL.
func getPageBody(pageURL string) (string, error) {
	resp, err := http.Get(pageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// getWebsiteData scrapes a website for an email and its technology type.
func getWebsiteData(websiteURL string) (string, string, error) {
	if websiteURL == "" {
		return "", "N/A", nil
	}

	// 1. Get homepage body
	homepageBody, err := getPageBody(websiteURL)
	if err != nil {
		return "", "", err
	}

	// 2. Detect site type from homepage
	siteType := detectSiteType(homepageBody)

	// 3. Try to find a "contact" link
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(homepageBody))
	if err != nil {
		return "", siteType, err
	}

	contactLink, found := "", false
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			if strings.Contains(strings.ToLower(href), "contact") || strings.Contains(strings.ToLower(s.Text()), "contact") {
				contactLink = href
				found = true
				return
			}
		}
	})

	// 4. Scrape for email
	email := ""
	if found {
		// Build absolute URL for the contact page
		contactURL, err := url.Parse(contactLink)
		if err == nil {
			base, _ := url.Parse(websiteURL)
			absoluteContactURL := base.ResolveReference(contactURL)

			contactPageBody, err := getPageBody(absoluteContactURL.String())
			if err == nil {
				email = findEmail(contactPageBody)
			}
		}
	}

	// 5. If no email on contact page, fall back to homepage
	if email == "" {
		email = findEmail(homepageBody)
	}

	return email, siteType, nil
}

func Search(apiKey, username, password, query string, pages int) {
	// Connect to the database
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require",
		"places-db.cwju0uoiwc8m.us-east-1.rds.amazonaws.com",
		username,
		password,
		"places",
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the businesses table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS businesses (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE,
			address TEXT,
			phone_number TEXT,
			website TEXT,
			email TEXT,
			site_type TEXT
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create businesses table: %v", err)
	}

	// Create a new Google Maps client
	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	// Perform a text search for businesses in Manchester
	req := &maps.TextSearchRequest{
		Query: query,
	}

	for i := 0; i < pages; i++ {
		log.Printf("Fetching page %d...", i+1)
		resp, err := c.TextSearch(context.Background(), req)
		if err != nil {
			log.Fatalf("fatal error on page %d: %s", i+1, err)
		}

		// Insert the results into the database
		for _, result := range resp.Results {
			detailReq := &maps.PlaceDetailsRequest{
				PlaceID: result.PlaceID,
				Fields:  []maps.PlaceDetailsFieldMask{"website", "formatted_phone_number"},
			}
			detailResp, err := c.PlaceDetails(context.Background(), detailReq)
			if err != nil {
				log.Printf("failed to get place details for %s: %s", result.Name, err)
				continue
			}

			email, siteType, err := getWebsiteData(detailResp.Website)
			if err != nil {
				log.Printf("could not scrape website %s: %v", detailResp.Website, err)
			}

			_, err = db.Exec(
				"INSERT INTO businesses (name, address, phone_number, website, email, site_type) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (name) DO NOTHING",
				result.Name,
				result.FormattedAddress,
				detailResp.FormattedPhoneNumber,
				detailResp.Website,
				email,
				siteType,
			)

			if err != nil {
				log.Printf("failed to insert business %s: %s", result.Name, err)
			} else {
				fmt.Printf("Processed: %s (Email: %s, Site Type: %s)\n", result.Name, email, siteType)
			}
		}

		// Check if there is a next page
		if resp.NextPageToken == "" {
			log.Println("No more pages to fetch.")
			break
		}

		// Prepare the request for the next page
		req.PageToken = resp.NextPageToken
		// It's important to wait before making the next request
		log.Println("Waiting before fetching next page...")
		time.Sleep(2 * time.Second)
	}
}

func Dump(username, password, format string) {
	// Connect to the database
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require",
		"places-db.cwju0uoiwc8m.us-east-1.rds.amazonaws.com",
		username,
		password,
		"places",
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query the database
	rows, err := db.Query("SELECT name, address, phone_number, website, email, site_type FROM businesses")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if format == "csv" {
		// Create a new CSV file
		file, err := os.Create("dump.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Create a new CSV writer
		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Write the header row
		writer.Write([]string{"Name", "Address", "Phone Number", "Website", "Email", "Site Type"})

		// Write the data rows
		for rows.Next() {
			var name, address, phoneNumber, website, email, siteType string
			err := rows.Scan(&name, &address, &phoneNumber, &website, &email, &siteType)
			if err != nil {
				log.Fatal(err)
			}
			writer.Write([]string{name, address, phoneNumber, website, email, siteType})
		}

		fmt.Println("Database dumped to dump.csv")
	} else if format == "sql" {
		// Create a new SQL file
		file, err := os.Create("dump.sql")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Write the data rows
		for rows.Next() {
			var name, address, phoneNumber, website, email, siteType string
			err := rows.Scan(&name, &address, &phoneNumber, &website, &email, &siteType)
			if err != nil {
				log.Fatal(err)
			}
			file.WriteString(fmt.Sprintf("INSERT INTO businesses (name, address, phone_number, website, email, site_type) VALUES ('%s', '%s', '%s', '%s', '%s', '%s');\n", name, address, phoneNumber, website, email, siteType))
		}

		fmt.Println("Database dumped to dump.sql")
	} else {
		log.Fatalf("Invalid format: %s. Please use 'csv' or 'sql'.", format)
	}
}