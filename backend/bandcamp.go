package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strings"
// 	"time"
// )

// type BandcampResponse struct {
// 	Results []BandcampArtist `json:"results"`
// }

// type BandcampArtist struct {
// 	BandName     string `json:"band_name"`
// 	BandURL      string `json:"band_url"`
// 	BandLocation string `json:"band_location"`
// 	Title        string `json:"title"`
// 	ReleaseDate  string `json:"release_date"`
// }

// func main() {
// 	fmt.Println("Testing Bandcamp endpoint...")

// 	url := "https://bandcamp.com/api/discover/1/discover_web"

// 	// Test with larger size
// 	jsonData := `{
// 		"category_id": 0,
// 		"tag_norm_names": [],
// 		"geoname_id": 4930956,
// 		"slice": "new",
// 		"time_facet_id": null,
// 		"cursor": "*",
// 		"size": 60,
// 		"include_result_types": ["a", "s"]
// 	}`

// 	client := &http.Client{Timeout: 10 * time.Second}

// 	req, err := http.NewRequest("POST", url, strings.NewReader(jsonData))
// 	if err != nil {
// 		fmt.Printf("Error creating request: %v\n", err)
// 		return
// 	}

// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
// 	req.Header.Set("Accept", "application/json")
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Referer", "https://bandcamp.com/discover")
// 	req.Header.Set("X-Requested-With", "XMLHttpRequest")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Printf("Request failed: %v\n", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	fmt.Printf("Status: %d\n", resp.StatusCode)

// 	if resp.StatusCode != 200 {
// 		fmt.Printf("Got non-200 status: %d\n", resp.StatusCode)
// 		return
// 	}

// 	// First, let's see what we actually got
// 	var rawResponse map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
// 		fmt.Printf("Failed to parse JSON: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("Raw response keys: %v\n", getKeys(rawResponse))

// 	// Check if results exist
// 	if results, ok := rawResponse["results"].([]interface{}); ok {
// 		fmt.Printf("Found %d results\n", len(results))

// 		// Check result_count vs actual results
// 		if resultCount, ok := rawResponse["result_count"].(float64); ok {
// 			fmt.Printf("Total available results: %.0f\n", resultCount)
// 		}

// 		// Print first few results to see structure
// 		for i, result := range results {
// 			if i >= 3 {
// 				break
// 			}
// 			if resultMap, ok := result.(map[string]interface{}); ok {
// 				fmt.Printf("Result %d: %s (%s)\n", i, resultMap["band_name"], resultMap["band_location"])
// 			}
// 		}
// 	} else {
// 		fmt.Println("No 'results' field found")
// 		fmt.Printf("Full response: %v\n", rawResponse)
// 	}

// 	fmt.Println("Success! Endpoint is accessible.")
// }

// func getKeys(m map[string]interface{}) []string {
// 	keys := make([]string, 0, len(m))
// 	for k := range m {
// 		keys = append(keys, k)
// 	}
// 	return keys
// }
