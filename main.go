package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type AnimeInfo struct {
	Synopsis    string
	Type        string
	startedAt   time.Time
	endedAt     time.Time
	Producers   []string
	Licensors   []string
	Studios     []string
	Source      string
	Genres      []string
	Themes      []string
	Demographic []string
	Duration    int // in minutes
	Rating      string
}

func main() {
	c := colly.NewCollector(
		// Restrict domains to myanimelist.net
		colly.AllowedDomains("myanimelist.net", "www.myanimelist.net"),
		// Local cache to prevent multiples download
		colly.CacheDir("./my_anime_list_cache"),
		// Limit to depth 2 to control parallel
		colly.MaxDepth(2),
		colly.Async(),
	)

	// Add a detail collector to scrape specific anime information
	detailCollector := c.Clone()

	c.OnHTML("a[href][class=genre-name-link]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", string(r.Body), "\nError:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnHTML("a[href][class=link-title]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		if strings.Index(link, "myanimelist.net/anime/") != -1 {
			parts := strings.Split(link, "/")
			N := len(parts)
			log.Println("Anime found: ", parts[N-1])
			detailCollector.Visit(e.Request.AbsoluteURL(link))
		}
	})

	detailCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", string(r.Body), "\nError:", err)
	})

	detailCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnHTML("a[href][class=link]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.Visit("https://myanimelist.net/anime.php#")
	// Wait until threads finishes
	c.Wait()

}
