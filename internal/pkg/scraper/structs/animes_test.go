package structs

import (
	"reflect"
	"testing"
)

func TestAnimeInit(t *testing.T) {
	want := &Anime{
		Title:       "",
		JTitle:      "",
		Synopsis:    "",
		Type:        "",
		NEpisodes:   "",
		Status:      "",
		Aired:       "",
		Premiered:   "",
		Broadcast:   "",
		Producers:   []string{},
		Licensors:   []string{},
		Studios:     []string{},
		Source:      "",
		Genres:      []string{},
		Themes:      []string{},
		Demographic: []string{},
		Duration:    "",
		Rating:      "",
		Score:       "",
		Ranked:      "",
		Popularity:  "",
		Members:     "",
		Favorites:   "",
	}

	got := &Anime{}
	got.Init()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("want %s, but got %s", want, got)
	}
}
