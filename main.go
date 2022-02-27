package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

var animeUrl string

func init() {
	flag.StringVar(&animeUrl, "animeUrl", "https://foo/bar", "The animeUrl page that you want to parse")
}

type AnimeInfo struct {
	Synopsis    string
	Type        string
	From        time.Time
	To          time.Time
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

func CleanSynopsis(s string) (cleanedSynopsis string) {
	cleanedSynopsis = strings.TrimSpace(strings.Split(s, "\n")[0])
	return
}

func GetDivInfo(s string) (divInfo string) {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	divInfo = strings.TrimSpace(lines[1])
	if strings.Contains(divInfo, "None found") {
		divInfo = "null"
	}
	return
}

func GetDivInfoNested(s *goquery.Selection) (divInfoNested []string) {
	divInfoNested = make([]string, 0)
	s.Find("a[href]").Each(func(_ int, is *goquery.Selection) {
		divInfoNested = append(divInfoNested, is.Text())
	})
	return
}

func ExtractFromToDates(s string) (from time.Time, to time.Time, err error) {
	shortLayout := "Jan 2, 2006"

	lines := strings.Split(strings.TrimSpace(s), "\n")
	dates := strings.TrimSpace(lines[1])
	parts := strings.Split(dates, " to ")
	if len(parts) == 1 {
		is := parts[0]
		if matched, _ := regexp.MatchString(`[0-9]{4}`, is); matched {
			shortLayout = "2006"
			year, err := time.Parse(shortLayout, is)
			if err != nil {
				return time.Time{}, time.Time{}, err
			}
			return year, time.Time{}, err
		}
		from, err = time.Parse(shortLayout, is)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return from, time.Time{}, err
	}
	fromString := parts[0]
	if fromString == "?" {
		from = time.Time{}
	} else {
		from, err = time.Parse(shortLayout, fromString)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}
	toString := parts[1]
	if toString == "?" {
		to = time.Time{}
	} else {
		to, err = time.Parse(shortLayout, parts[1])
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}
	return from, to, nil
}

func main() {
	c := colly.NewCollector(
		// TODO: Add Redis Backend
		// Local cache to prevent multiples download
		colly.CacheDir("./my_anime_list_cache"),
		// Uncomment to add a debugger
		//colly.Debugger(&debug.LogDebugger{}),
		colly.Async(true),
		colly.MaxDepth(2),
	)

	c.Limit(&colly.LimitRule{
		// Restrict domains to myanimelist.net
		DomainGlob: "*myanimelist.*",
		// Add a random delay
		RandomDelay: 15 * time.Second,
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

	c.OnHTML("a[href][class=link]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	detailCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", string(r.Body), "\nError:", err)
	})

	// Extract japanese title
	detailCollector.OnHTML("h1.title-name > strong", func(e *colly.HTMLElement) {
		fmt.Println("Starting parsing an anime page")
		fmt.Println("")
		fmt.Println("Title:", e.Text)
	})

	// Extract english title
	detailCollector.OnHTML("p.title-english", func(e *colly.HTMLElement) {
		fmt.Println("English title: ", e.Text)
	})

	// Extract synopsis
	detailCollector.OnHTML("p[itemprop=description]", func(e *colly.HTMLElement) {
		cleanedSynopsis := CleanSynopsis(e.Text)
		fmt.Println("")
		fmt.Println("Synopsis: ", cleanedSynopsis)
		fmt.Println("")
	})

	// Extract information & statistics
	detailCollector.OnHTML("div.spaceit_pad", func(e *colly.HTMLElement) {
		doc := e.DOM
		doc.Each(func(_ int, s *goquery.Selection) {
			switch strings.Trim(s.Find("span.dark_text").Text(), ":") {
			case "Type":
				fmt.Println("Type: ", GetDivInfo(s.Text()))
			case "Episodes":
				fmt.Println("Episodes: ", GetDivInfo(s.Text()))
			case "Status":
				fmt.Println("Status: ", GetDivInfo(s.Text()))
			case "Aired":
				from, to, err := ExtractFromToDates(s.Text())
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("From: ", from)
				fmt.Println("To: ", to)
			case "Premiered":
				fmt.Println("Premiered: ", GetDivInfo(s.Text()))
			case "Broadcast":
				fmt.Println("Broadcast: ", GetDivInfo(s.Text()))
			case "Producers":
				fmt.Println("Producers: ", GetDivInfoNested(s))
			case "Licensors":
				fmt.Println("Licensors: ", GetDivInfoNested(s))
			case "Studios":
				fmt.Println("Studios: ", GetDivInfoNested(s))
			case "Source":
				fmt.Println("Source: ", GetDivInfo(s.Text()))
			case "Genres":
				fmt.Println("Genres: ", GetDivInfoNested(s))
			case "Themes":
				fmt.Println("Themes: ", GetDivInfoNested(s))
			case "Demographic":
				fmt.Println("Demographic: ", GetDivInfoNested(s))
			case "Duration":
				fmt.Println("Duration: ", GetDivInfo(s.Text()))
			case "Rating":
				fmt.Println("Rating: ", GetDivInfo(s.Text()))
			}
		})
	})

	detailCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Parse flag
	flag.Parse()
	if animeUrl != "https://foo/bar" {
		detailCollector.Visit(animeUrl)
		detailCollector.Wait()
	}
	c.Visit("https://myanimelist.net/anime.php#")
	c.Wait()

}
