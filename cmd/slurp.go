package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("247sports.com"),
	)
	done := false
	var players []player
	pagePlayerCount := 0
	c.OnHTML("li.rankings-page__list-item", func(e *colly.HTMLElement) {
		e.ForEach("div.wrapper", func(i int, element *colly.HTMLElement) {
			var player player
			pagePlayerCount++
			primary := e.ChildText(".primary")
			matchRank := regexp.MustCompile(`(\d+)`)
			primaryRank := matchRank.FindStringSubmatch(primary)
			rank, err := strconv.ParseInt(primaryRank[1], 10, 64)
			if err != nil {
				fmt.Println(fmt.Errorf("%w", err))
			}
			player.Rank = rank
			player.Name = e.ChildText("a.rankings-page__name-link")
			player.Position = e.ChildText("div.position")

			metricsText := e.ChildText("div.metrics")
			metricMatch := regexp.MustCompile(`(\d)-(\d+\.?\d*) / (\d+)`)
			metrics := metricMatch.FindStringSubmatch(metricsText)
			heightFeet, err := strconv.ParseInt(metrics[1], 10, 64)
			if err != nil {
				fmt.Println(fmt.Errorf("%w", err))
			}
			heightInches, err := strconv.ParseFloat(metrics[2], 64)
			if err != nil {
				fmt.Println(fmt.Errorf("%w", err))
			}
			weightPounds, err := strconv.ParseInt(metrics[3], 10, 64)
			if err != nil {
				fmt.Println(fmt.Errorf("%w", err))
			}
			player.Height = float64(heightFeet)*12.0 + heightInches
			player.Weight = weightPounds

			scoreText := e.ChildText("span.score")
			score, err := strconv.ParseFloat(scoreText, 64)
			if err != nil {
				fmt.Println(fmt.Errorf("%w", err))
			}
			player.Score = score
			players = append(players, player)
		})
	})

	c.OnScraped(func(response *colly.Response) {
		if pagePlayerCount < 50 {
			done = true
		}
		pagePlayerCount = 0
	})

	// Start scraping on 247sports.com
	for i := 1; ; i++ {
		err := c.Visit(fmt.Sprintf("https://247sports.com/Season/2023-Football/CompositeRecruitRankings/?ViewPath=~/Views/SkyNet/PlayerSportRanking/_SimpleSetForSeason.ascx&InstitutionGroup=HighSchool&Page=%d", i))
		if err != nil {
			fmt.Println(fmt.Errorf("%w", err))
		}
		if done {
			break
		}
	}
	for _, player := range players {
		fmt.Printf("%+v\n", player)
	}
}

type player struct {
	Rank     int64
	Name     string
	Position string
	Height   float64 // inches
	Weight   int64   // pounds
	Score    float64
}
