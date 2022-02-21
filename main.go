package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector(
		// Restrict domains to myanimelist.net
		colly.AllowedDomains("myanimelist.net"),
		// Local cache to prevent multiples download
		colly.CacheDir("./my_anime_list_cache"),
		// Limit to depth 0 for debugging
		colly.MaxDepth(0),
	)

	c.OnHTML("div[id=horiznav_nav] > ul > li > a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		switch e.Text {
		case "Upcoming", "Just Added":
			return
		}
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://myanimelist.net/anime.php#")

}
