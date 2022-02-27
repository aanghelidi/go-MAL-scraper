package main

import (
	"regexp"
	"strings"
	"time"

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

func ExtractFromToDates(s string) (from time.Time, to time.Time, err error) {
	shortLayout := "Jan 2, 2006"

	lines := strings.Split(strings.TrimSpace(s), "\n")
	dates := strings.TrimSpace(lines[1])
	parts := strings.Split(dates, " to ")
	if len(parts) == 1 {
		is := parts[0]
		if matched, _ := regexp.MatchString(`[0-9]{4}`, is); matched && len(is) == 4 {
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
