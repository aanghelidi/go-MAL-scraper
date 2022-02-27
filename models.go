package main

import "time"

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
