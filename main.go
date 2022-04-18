package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/NewMirai/go-MAL-scraper/internal/malctl/scraper/structs"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

func main() {
	fName := "/tmp/data/animes.json"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()

	c := colly.NewCollector(
		// TODO: Add Redis Backend
		// Local cache to prevent multiples download
		colly.CacheDir("./my_anime_list_cache"),
		colly.MaxDepth(1),
	)

	c.Limit(&colly.LimitRule{
		// Restrict domains to myanimelist.net
		DomainGlob: "*myanimelist.*",
		// Add a random delay
		RandomDelay: 30 * time.Second,
		// Add parallelism
		Parallelism: 2,
	})

	// Generate a random user agent on every request
	extensions.RandomUserAgent(c)

	//Add a detail collector to scrape specific anime information
	detailCollector := colly.NewCollector(
		colly.Async(false),
		colly.MaxDepth(0),
	)
	detailCollector.Limit(&colly.LimitRule{
		// Restrict domains to myanimelist.net
		DomainGlob: "*myanimelist.*",
		// Add a random delay
		RandomDelay: 15 * time.Second,
	})

	// Generate a random user agent on every request
	extensions.RandomUserAgent(detailCollector)

	// Initialization
	animes := make([]structs.Anime, 0)
	count := 0
	// Parse flag
	var animeUrl string
	var nAnimes int
	flag.StringVar(&animeUrl, "animeUrl", "https://foo/bar", "The animeUrl page that you want to parse")
	flag.IntVar(&nAnimes, "nAnimes", 0, "If different of 0 limit the numbers of animes to parse, default to 0")
	flag.Parse()

	c.OnHTML("a[href][class=genre-name-link]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		log.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", string(r.Body), "\nError:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
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
		log.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	detailCollector.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", string(r.Body), "\nError:", err)
	})

	detailCollector.OnHTML("div#contentWrapper", func(e *colly.HTMLElement) {
		// Write JSON if nAnimes is reached
		if count != 0 && count == nAnimes {
			enc := json.NewEncoder(file)
			enc.SetIndent("", "  ")
			// Dump json to the standard output
			enc.Encode(animes)
			log.Println("Finished writing ", nAnimes, " animes in the JSON file")
			os.Exit(0)
		}

		//  Start parsing the page
		anime := Anime{}
		anime.JTitle = e.ChildText("h1.title-name > strong")
		anime.Title = e.ChildText("p.title-english")
		anime.Synopsis = CleanSynopsis(e.ChildText("p[itemprop=description]"))

		// Extract information and statistics
		e.ForEach("div.spaceit_pad", func(_ int, s *colly.HTMLElement) {
			switch strings.Trim(s.ChildText("span.dark_text"), ":") {
			case "Type":
				anime.Type = GetDivInfo(s.Text)
			case "Episodes":
				anime.NEpisodes = GetDivInfo(s.Text)
			case "Status":
				anime.Status = GetDivInfo(s.Text)
			case "Aired":
				anime.Aired = GetDivInfo(s.Text)
			case "Premiered":
				anime.Premiered = GetDivInfo(s.Text)
			case "Broadcast":
				anime.Broadcast = GetDivInfo(s.Text)
			case "Producers":
				anime.Producers = GetDivInfoNested(s.DOM)
			case "Licensors":
				anime.Licensors = GetDivInfoNested(s.DOM)
			case "Studios":
				anime.Studios = GetDivInfoNested(s.DOM)
			case "Source":
				anime.Source = GetDivInfo(s.Text)
			case "Genres":
				anime.Genres = GetDivInfoNested(s.DOM)
			case "Themes":
				anime.Themes = GetDivInfoNested(s.DOM)
			case "Demographic":
				anime.Demographic = GetDivInfoNested(s.DOM)
			case "Duration":
				anime.Duration = GetDivInfo(s.Text)
			case "Rating":
				anime.Rating = GetDivInfo(s.Text)
			case "Score":
				anime.Score = GetDivInfo(s.Text)
			case "Ranked":
				anime.Ranked = GetDivInfo(s.Text)
			case "Popularity":
				anime.Popularity = GetDivInfo(s.Text)
			case "Members":
				anime.Members = GetDivInfo(s.Text)
			case "Favorites":
				anime.Favorites = GetDivInfo(s.Text)
			}
		})
		animes = append(animes, anime)
		count++
		log.Println("Animes currently parsed: ", count)
	})

	detailCollector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	detailCollector.OnError(func(r *colly.Response, e error) {
		log.Println("Error: ", e)
	})

	if animeUrl != "https://foo/bar" {
		detailCollector.Visit(animeUrl)
	}

	c.Visit("https://myanimelist.net/anime.php#")
	c.Wait()

	// If all is parsed then ...
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	// Dump json to the standard output
	enc.Encode(animes)
	return
}
