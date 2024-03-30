package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"io"
	"log"
	"net/http"
	"os"
)

const LeaderboardURL = "https://api.pubg.com/shards/xbox-na/leaderboards/division.bro.official.console-28/squad"

var API_KEY = os.Getenv("PUBG_TOKEN")
var REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")

type LeaderboardEntry struct {
	AccountId string      `json:"accountId"`
	Stats     LeaderStats `json:"stats"`
}

type LeaderStats struct {
	Rank  int `json:"rank"`
	Wins  int `json:"wins"`
	Games int `json:"games"`
}

type LeaderboardResponse struct {
	Data struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes struct {
			ShardID  string `json:"shardId"`
			GameMode string `json:"gameMode"`
			SeasonID string `json:"seasonId"`
		} `json:"attributes"`
		Relationships struct {
			Players struct {
				Data []struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				} `json:"data"`
			} `json:"players"`
		} `json:"relationships"`
	} `json:"data"`
	Included []struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes struct {
			Name  string `json:"name"`
			Rank  int    `json:"rank"`
			Stats struct {
				RankPoints     int     `json:"rankPoints"`
				Wins           int     `json:"wins"`
				Games          int     `json:"games"`
				WinRatio       int     `json:"winRatio"`
				AverageDamage  int     `json:"averageDamage"`
				Kills          int     `json:"kills"`
				KillDeathRatio int     `json:"killDeathRatio"`
				Kda            float64 `json:"kda"`
				AverageRank    float64 `json:"averageRank"`
				Tier           string  `json:"tier"`
				SubTier        string  `json:"subTier"`
			} `json:"stats"`
		} `json:"attributes"`
	} `json:"included"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
	Meta struct {
	} `json:"meta"`
}

var ctx = context.Background()

func UpdateRedis(leaderboardEntries []LeaderboardEntry) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-c82bfc2f:6379",
		Password: REDIS_PASSWORD,
		DB:       0,
	})

	for _, entry := range leaderboardEntries {
		fmt.Println(entry)
		err := rdb.HSet(ctx, entry.AccountId, map[string]interface{}{
			"rank":  entry.Stats.Rank,
			"wins":  entry.Stats.Wins,
			"games": entry.Stats.Games,
		}).Err()
		if err != nil {
			panic(err)
		}
	}
}

func GetLeaderboard() LeaderboardResponse {
	client := &http.Client{}

	req, err := http.NewRequest("GET", LeaderboardURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+API_KEY)
	req.Header.Set("Accept", "application/vnd.api+json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		panic(err)
	}

	var leaderboardResponse LeaderboardResponse
	err = json.Unmarshal(body, &leaderboardResponse)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		panic(err)
	}

	return leaderboardResponse
}

func PrepareData(response LeaderboardResponse) []LeaderboardEntry {
	var statsList []LeaderboardEntry
	for _, entry := range response.Included {
		stats := LeaderStats{
			Rank:  entry.Attributes.Rank,
			Wins:  entry.Attributes.Stats.Wins,
			Games: entry.Attributes.Stats.Games,
		}
		statsList = append(statsList, LeaderboardEntry{AccountId: entry.ID, Stats: stats})
	}
	return statsList
}

func main() {

	log.Println("Fetching new leaderboard data")
	leaderboardResponse := GetLeaderboard()
	parsedData := PrepareData(leaderboardResponse)
	UpdateRedis(parsedData)
}
