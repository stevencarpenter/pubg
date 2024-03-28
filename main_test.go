package main

import (
	"reflect"
	"testing"
)

func TestPrepareData(t *testing.T) {
	type args struct {
		response LeaderboardResponse
	}
	tests := []struct {
		name string
		args args
		want []LeaderboardEntry
	}{
		{
			name: "parses correctly",
			args: args{LeaderboardResponse{
				Data: TopData{
					Included: Included{
						{Id: "player1", Attributes: Attributes{Rank: 6, Stats: Stats{Games: 10, Wins: 7}}},
						{Id: "player2", Attributes: Attributes{Rank: 5, Stats: Stats{Games: 10, Wins: 5}}},
						{Id: "player3", Attributes: Attributes{Rank: 4, Stats: Stats{Games: 10, Wins: 4}}},
					},
				},
			}},
			want: []LeaderboardEntry{
				{AccountId: "player1", Stats: LeaderStats{Rank: 6, Wins: 7, Games: 10}},
				{AccountId: "player2", Stats: LeaderStats{Rank: 5, Wins: 5, Games: 10}},
				{AccountId: "player3", Stats: LeaderStats{Rank: 4, Wins: 4, Games: 10}},
				//{AccountId: "player4", Stats: LeaderStats{Rank: 3, Wins: 3, Games: 10}},
				//{AccountId: "player5", Stats: LeaderStats{Rank: 2, Wins: 2, Games: 10}},
				//{AccountId: "player6", Stats: LeaderStats{Rank: 1, Wins: 1, Games: 10}},
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrepareData(tt.args.response); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrepareData() = %v, want %v", got, tt.want)
			}
		})
	}
}
