package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// RSS struct matches the Google News RSS XML schema
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
	Source      string `xml:"source"`
}

func main() {
	searchQuery := "Artificial Intelligence"

	// Format the Google News RSS URL safely
	feedURL := fmt.Sprintf("https://www.google.com?q=%s", url.QueryEscape(searchQuery))

	// Make the HTTP Request
	resp, err := http.Get(feedURL)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Status code %d\n", resp.StatusCode)
		return
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
		return
	}

	// Parse the XML data
	var rss RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		fmt.Printf("Error parsing XML: %v\n", err)
		return
	}

	// Display the results
	fmt.Printf("--- Google News Results for: %s ---\n\n", rss.Channel.Title)
	for i, item := range rss.Channel.Items {
		if i >= 10 { // Limit to top 10 articles
			break
		}
		fmt.Printf("[%d] %s\n", i+1, item.Title)
		fmt.Printf("    Link: %s\n", item.Link)
		fmt.Printf("    Published: %s\n\n", item.PubDate)
	}
}
