package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

func GoogleScrape() {
	logrus.Info("GoogleScrape")
	// Request the HTML page.
	res, err := http.Get("https://www.google.com/search?q=latest+news+ai&gl=us&hl=en")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	logrus.Info("GoogleScrape query")
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	c := 0
	doc.Find("#search div").Each(func(i int, result *goquery.Selection) {
		title := result.Find("cite").()
		// link, _ := result.Find("a").First().Attr("href")
		// snippet := result.Find(".VwiC3b").First().Text()

		fmt.Printf("Title: %s\n", title)
		// fmt.Printf("Link: %s\n", link)
		// fmt.Printf("Snippet: %s\n", snippet)
		// fmt.Printf("Position: %d\n", c+1)
		fmt.Println()

		c++
	})
}

func ExampleScrape() {
	// Request the HTML page.
	res, err := http.Get("http://www.dronasys.com")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("#mainContentBlock section").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find("h2").Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})
}

func main() {
	ExampleScrape()
	GoogleScrape()
}
