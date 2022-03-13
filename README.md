# go-MAL-scraper
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
