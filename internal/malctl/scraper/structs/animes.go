package structs

type Anime struct {
	Title       string   `json:"title"`
	JTitle      string   `json:"j_title"`
	Synopsis    string   `json:"synopsis"`
	Type        string   `json:"type"`
	NEpisodes   string   `json:"n_episodes"`
	Status      string   `json:"status"`
	Aired       string   `json:"aired"`
	Premiered   string   `json:"premiered"`
	Broadcast   string   `json:"broadcast"`
	Producers   []string `json:"producers"`
	Licensors   []string `json:"licensors"`
	Studios     []string `json:"studios"`
	Source      string   `json:"source"`
	Genres      []string `json:"genres"`
	Themes      []string `json:"themes"`
	Demographic []string `json:"demographic"`
	Duration    string   `json:"duration"`
	Rating      string   `json:"rating"`
	Score       string   `json:"score"`
	Ranked      string   `json:"ranked"`
	Popularity  string   `json:"popularity"`
	Members     string   `json:"members"`
	Favorites   string   `json:"favourites"`
}

type AnimeField func(*Anime)

func (a *Anime) Init() {
	a.Title = ""
	a.JTitle = ""
	a.Synopsis = ""
	a.Type = ""
	a.NEpisodes = ""
	a.Status = ""
	a.Aired = ""
	a.Premiered = ""
	a.Broadcast = ""
	a.Producers = []string{}
	a.Licensors = []string{}
	a.Studios = []string{}
	a.Source = ""
	a.Genres = []string{}
	a.Themes = []string{}
	a.Demographic = []string{}
	a.Duration = ""
	a.Rating = ""
	a.Score = ""
	a.Ranked = ""
	a.Popularity = ""
	a.Members = ""
	a.Favorites = ""
}

func Title(title string) AnimeField {
	return func(a *Anime) {
		a.Title = title
	}
}

func JTitle(jtitle string) AnimeField {
	return func(a *Anime) {
		a.JTitle = jtitle
	}
}
