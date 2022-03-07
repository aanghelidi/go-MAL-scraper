package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
