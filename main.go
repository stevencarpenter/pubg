package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"io"
	"net/http"
)

const (
	APIKey         = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJqdGkiOiJlMDcwOTI1MC1jYWNjLTAxM2MtOGMwYS02YTM5YWYyZTFjYWQiLCJpc3MiOiJnYW1lbG9ja2VyIiwiaWF0IjoxNzExMTQ3OTg5LCJwdWIiOiJibHVlaG9sZSIsInRpdGxlIjoicHViZyIsImFwcCI6Ii1jNWIwMjQ4MS1lYzg3LTRhZGMtOWZmZi01N2IxMzhhNmNkZGEifQ.UcRqxmr0qSIx_SBM2H0i-XbKniCEycyUGjV_JFvj9NI" // Replace this with your PUBG API key
	LeaderboardURL = "https://api.pubg.com/shards/pc-na/leaderboards/division.bro.official.pc-2018-28/duo"                                                                                                                                                                                                                     // Replace {platform} and {gameMode} as needed
)

type LeaderboardEntry struct {
	AccountId string      `json:"accountId"`
	Stats     LeaderStats `json:"stats"`
}

type LeaderStats struct {
	Rank  int `json:"rank"`
	Wins  int `json:"wins"`
	Games int `json:"games"`
}

type Included []struct {
	Id         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	Name  string `json:"name"`
	Rank  int    `json:"rank"`
	Stats Stats  `json:"stats"`
}

type Stats struct {
	RankPoints     int     `json:"rankPoints"`
	Wins           int     `json:"wins"`
	Games          int     `json:"games"`
	WinRatio       float64 `json:"winRatio"`
	AverageDamage  float64 `json:"averageDamage"`
	Kills          int     `json:"kills"`
	KillDeathRatio float64 `json:"killDeathRatio"`
	KDA            float64 `json:"kda"`
	AverageRank    float64 `json:"averageRank"`
	Tier           string  `json:"tier"`
	SubTier        string  `json:"subTier"`
}

type Data []struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type Players struct {
	Data Data `json:"data"`
}

type Relationships struct {
	Players Players `json:"players"`
}

type SeasonAttributes struct {
	GameMode string `json:"gameMode"`
	ShardId  string `json:"shardId"`
	SeasonId string `json:"seasonId"`
}

type TopData struct {
	Attributes    SeasonAttributes `json:"attributes"`
	Relationships Relationships    `json:"relationships"`
	Included      Included         `json:"included"`
	Type          string           `json:"type"`
	Id            string           `json:"id"`
}

type Links struct {
	Self string `json:"self"`
}

type Meta struct{}

type LeaderboardResponse struct {
	Meta  Meta    `json:"meta"`
	Links Links   `json:"links"`
	Data  TopData `json:"data"`
}

var ctx = context.Background()

func UpdateRedis(leaderboardEntries []LeaderboardEntry) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-c82bfc2f:6379",
		Password: "password", // no password set
		DB:       0,          // use default DB
	})

	for _, entry := range leaderboardEntries {
		fmt.Println(entry)
		err := rdb.Set(ctx, entry.AccountId, entry.Stats, 0).Err()
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

	req.Header.Set("Authorization", "Bearer "+APIKey)
	req.Header.Set("Accept", "application/vnd.api+json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		panic(err)
	}
	defer resp.Body.Close()

	println(resp.Body)

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

	// Output some information to verify that the call was successful.
	fmt.Println("Game Mode:", leaderboardResponse.Data.Attributes.GameMode)
	for _, player := range leaderboardResponse.Data.Included {
		fmt.Println("Player ID:", player.Id)
		fmt.Println("Player Rank:", player.Attributes.Rank)
		fmt.Println("Player Wins:", player.Attributes.Stats.Wins)
		fmt.Println("Player Games:", player.Attributes.Stats.Games)
	}

	return leaderboardResponse
}

func PrepareData(response LeaderboardResponse) []LeaderboardEntry {
	var statsList []LeaderboardEntry
	for _, entry := range response.Data.Included {
		stats := LeaderStats{
			Rank:  entry.Attributes.Rank,
			Wins:  entry.Attributes.Stats.Wins,
			Games: entry.Attributes.Stats.Games,
		}
		statsList = append(statsList, LeaderboardEntry{AccountId: entry.Id, Stats: stats})
	}
	return statsList
}

func main() {
	//
	//leaderboardResponse := GetLeaderboard()
	//parsedData := PrepareData(leaderboardResponse)
	//UpdateRedis(parsedData)

	leaderboard := []LeaderboardEntry{
		{AccountId: "player1", Stats: LeaderStats{Rank: 6, Wins: 7, Games: 10}},
		{AccountId: "player2", Stats: LeaderStats{Rank: 5, Wins: 5, Games: 10}},
		{AccountId: "player3", Stats: LeaderStats{Rank: 4, Wins: 4, Games: 10}},
		{AccountId: "player4", Stats: LeaderStats{Rank: 3, Wins: 3, Games: 10}},
		{AccountId: "player5", Stats: LeaderStats{Rank: 2, Wins: 2, Games: 10}},
		{AccountId: "player6", Stats: LeaderStats{Rank: 1, Wins: 1, Games: 10}},
	}

	UpdateRedis(leaderboard)
	//UpdateRedis(parsedData)
}
