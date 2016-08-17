package giantbomb

import (
	"fmt"
	"log"
	"net/url"

	"github.com/corybuecker/jsonfetcher"
	"github.com/corybuecker/steam-stats-fetcher/database"
	"github.com/corybuecker/steam-stats-fetcher/ratelimiters"
)

type Search struct {
	Results []SearchResult `json:"results"`
}

type SearchResult struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"site_detail_url"`
}

type Fetcher struct {
	GiantBombAPIKey string `bson:"giantbombApiKey"`
	SearchResults   Search
	RateLimiter     *ratelimiters.GiantBombRateLimiter
	Jsonfetcher     jsonfetcher.Fetcher
}

func (fetcher *Fetcher) generateFetchURL(id int) string {
	return fmt.Sprintf("http://www.giantbomb.com/api/games/?filter=id:%d&api_key=%s&format=json&field_list=id,name,site_detail_url",
		id,
		fetcher.GiantBombAPIKey)
}

func (fetcher *Fetcher) generateSearchURL(name string) string {
	return fmt.Sprintf("http://www.giantbomb.com/api/games/?api_key=%s&format=json&filter=name:%s&field_list=id,name,site_detail_url",
		fetcher.GiantBombAPIKey,
		url.QueryEscape(name))
}

func (fetcher *Fetcher) FindGameByID(id int) error {
	if fetcher.RateLimiter == nil {
		fetcher.RateLimiter = &ratelimiters.GiantBombRateLimiter{}
	}

	log.Printf("fetching %d in the GiantBomb API", id)

	err := fetcher.Jsonfetcher.Fetch(fetcher.generateFetchURL(id), &fetcher.SearchResults)

	if err := fetcher.RateLimiter.ObeyRateLimit(); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (fetcher *Fetcher) FindOwnedGame(gameName string) error {
	if fetcher.RateLimiter == nil {
		fetcher.RateLimiter = &ratelimiters.GiantBombRateLimiter{}
	}

	log.Printf("searching for %s in the GiantBomb API", gameName)

	if err := fetcher.Jsonfetcher.Fetch(fetcher.generateSearchURL(gameName), &fetcher.SearchResults); err != nil {
		return err
	}

	if err := fetcher.RateLimiter.ObeyRateLimit(); err != nil {
		return err
	}

	return nil
}

func (fetcher *Fetcher) UpdateFoundGames(id int, database database.Interface) error {
	if len(fetcher.SearchResults.Results) > 1 {
		log.Printf("more than one result was found for this game, skipping")
		return nil
	}

	for _, foundGame := range fetcher.SearchResults.Results {
		foundGameMap := map[string]interface{}{
			"giantbombId": foundGame.ID,
			"url":         foundGame.URL,
		}

		if err := database.Upsert(map[string]interface{}{"id": id}, foundGameMap); err != nil {
			return err
		}
	}
	return nil
}
