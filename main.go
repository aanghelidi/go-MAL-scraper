package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
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
		// TODO: Add Redis Backend
		// Local cache to prevent multiples download
		colly.CacheDir("./my_anime_list_cache"),
		// Uncomment to add a debugger
		// colly.Debugger(&debug.LogDebugger{}),
		colly.Async(true),
		colly.MaxDepth(2),
	)

	c.Limit(&colly.LimitRule{
		// Restrict domains to myanimelist.net
		DomainGlob: "*myanimelist.*",
		// Add a random delay
		RandomDelay: 5 * time.Second,
		// Add parallelism
		Parallelism: 2,
	})

	// Generate a random user agent on every request
	extensions.RandomUserAgent(c)

	//Add a detail collector to scrape specific anime information
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
