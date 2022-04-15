# go-MAL-scraper

## Description
A go scraper using `Colly` to scrape MyAnimelist website information.
For each anime page the scraper will fetch the following informations if there are present:

- title
- jtitle , its japanese variation
- synopsis, the cleaned synopsis
- type, TV or movie for example
- n_episodes, the numbers of episodes of the show
- status, current status of the show
- aired, aired period
- premiered, premiered period
- broadcast, when it's broadcasted
- producers, list of producers
- licensors, list of licensors
- studios, list of studios
- source, what is the source of the show
- genres, list of genres
- themes, list of themes
- demographic, population concerned by the show
- duration, duration of the show
- rating, rating among users
- score, score among users
- ranked, the ranking
- popularity, its popularity
- members, how many members are on this anime page
- favorites, how many members put this anime as favourite

## How-tos

### How to run the scraper

The following assumes you have golang installed.

To run on all the website
```golang
go run main.go cleanUtils.go
```

To display help
```golang
go run main.go cleanUtils.go -h
```

To scrape a single anime using its URL
```golang
go run main.go cleanUtils.go -animeUrl https//...
```

To scrape 10 animes and save them to `animes.json` 
```golang
go run main.go cleanUtils.go -nAnimes 10
```

### How to run the scraper (docker version)

The following assume you have `Docker` installed.

First build the image
```
docker build . -t go-mal-scraper:latest
```

To run on all the website
```
docker run --name mal-scraper --mount type=bind,source="$(pwd)"/data,target=/tmp/data go-mal-scraper
```

To display help
```
docker run --name mal-scraper --mount type=bind,source="$(pwd)"/data,target=/tmp/data go-mal-scraper -h
```

To scrape a single anime using its URL
```
docker run --name mal-scraper --mount type=bind,source="$(pwd)"/data,target=/tmp/data go-mal-scraper -animeUrl https//...
```

To scrape 10 animes and save them to `animes.json` 
```
docker run --name mal-scraper --mount type=bind,source="$(pwd)"/data,target=/tmp/data go-mal-scraper -nAnimes 10
```
