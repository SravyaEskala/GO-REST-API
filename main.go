package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Score struct {
	Match   string `json:"match"`
	Runs    int    `json:"runs"`
	Wickets int    `json:"wickets"`
}

type Player struct {
	Name   string  `json:"name"`
	ID     int     `json:"id"`
	Team   string  `json:"team"`
	Scores []Score `json:"scores"`
}

type OnlyPlayer struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	Team string `json:"team"`
}

type DisplayPlayer struct {
	Players []OnlyPlayer `json:"players"`
}

type OnlyScores struct {
	ID     int     `json:"id"`
	Scores []Score `json:"scores"`
}

type displayScores struct {
	PlayerScores []OnlyScores `json:"playerscores"`
}

type FantasyScore struct {
	Name   string `json:"name"`
	FScore int    `json:"fantasyscore"`
}

type CapHolder struct {
	PurpleCap string `json:"purpleCap"`
	OrangeCap string `json:"orangeCap"`
}

var playerservice []Player

var TempScoreData []Score

var FantasyScores []FantasyScore

// used to add the player to the slice
func postPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newplayerdata Player
	_ = json.NewDecoder(r.Body).Decode(&newplayerdata)
	// check for empty values and eliminate them in the records
	if newplayerdata.ID != 0 && newplayerdata.Name != "" {
		newplayerdata.Scores = nil
		playerservice = append(playerservice, newplayerdata)
	}
}

// used to add the player score
func postPlayerScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var newMatchScore Score
	_ = json.NewDecoder(r.Body).Decode(&newMatchScore)
	for index, item := range playerservice {
		id, _ := strconv.Atoi(params["id"])
		if item.ID == id {
			playerservice = append(playerservice[:index], playerservice[index+1:]...)
			item.Scores = append(item.Scores, newMatchScore)
			playerservice = append(playerservice, item)
			break
		}
	}
}

// Used to get player details (not completed)
func getPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var playerdetails DisplayPlayer
	var playerdetail OnlyPlayer
	for _, item := range playerservice {
		playerdetail.ID = item.ID
		playerdetail.Name = item.Name
		playerdetail.Team = item.Team
		playerdetails.Players = append(playerdetails.Players, playerdetail)
	}
	json.NewEncoder(w).Encode(playerdetails)
}

// Used to get player details along with there scores
func getPlayerScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tempScores displayScores
	var playerdetails OnlyScores
	for _, item := range playerservice {
		playerdetails.ID = item.ID
		playerdetails.Scores = item.Scores
		tempScores.PlayerScores = append(tempScores.PlayerScores, playerdetails)
	}
	json.NewEncoder(w).Encode(tempScores)
	tempScores.PlayerScores = nil
}

// used to calculate fantasy score
func fantasyScoreCal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	for _, item := range playerservice {
		var singlePFS FantasyScore
		singlePFS.Name = item.Name
		singlePFS.FScore = 0
		for _, temp := range item.Scores {
			if temp.Wickets > 0 {
				singlePFS.FScore = singlePFS.FScore + 10*temp.Wickets
			}
			if temp.Wickets > 5 {
				singlePFS.FScore = singlePFS.FScore + 50
			}
			if temp.Runs >= 30 {
				singlePFS.FScore = singlePFS.FScore + 20
			}
			if temp.Runs >= 50 {
				singlePFS.FScore = singlePFS.FScore + 50
			}
			if temp.Runs >= 100 {
				singlePFS.FScore = singlePFS.FScore + 100
			}
		}
		FantasyScores = append(FantasyScores, singlePFS)
	}
	json.NewEncoder(w).Encode(FantasyScores)
	FantasyScores = nil
}

// used to calculate capholders
func capHolders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var capholderdetails CapHolder
	for _, item := range playerservice {
		var PerformanceCal Score
		var totalwickets int = 0
		var totalruns int = 0
		for _, temp := range item.Scores {
			totalwickets = totalwickets + temp.Wickets
			totalruns = totalruns + temp.Runs
		}
		PerformanceCal.Match = item.Name
		PerformanceCal.Wickets = totalwickets
		PerformanceCal.Runs = totalruns
		TempScoreData = append(TempScoreData, PerformanceCal)
	}
	var maxruns int = 0
	var maxwickets int = 0
	for _, item := range TempScoreData {
		if maxwickets < item.Wickets {
			maxwickets = item.Wickets
			capholderdetails.PurpleCap = item.Match
		}
		if maxruns < item.Runs {
			maxruns = item.Runs
			capholderdetails.OrangeCap = item.Match
		}
	}
	json.NewEncoder(w).Encode(capholderdetails)
	TempScoreData = nil
}

func main() {
	r := mux.NewRouter()
	TempScoreData = append(TempScoreData, Score{Match: "1", Wickets: 2, Runs: 150})
	playerservice = append(playerservice, Player{
		ID:     1,
		Name:   "SURYA KUMAR",
		Team:   "MI",
		Scores: TempScoreData,
	})
	TempScoreData = nil
	TempScoreData = append(TempScoreData, Score{Match: "1", Wickets: 3, Runs: 50})
	playerservice = append(playerservice, Player{
		ID:     7,
		Name:   "K L RAHUL",
		Team:   "GT",
		Scores: TempScoreData,
	})
	TempScoreData = nil
	r.HandleFunc("/player", postPlayer).Methods("POST")
	r.HandleFunc("/player/{id}/score", postPlayerScore).Methods("POST")
	r.HandleFunc("/players", getPlayers).Methods("GET")
	r.HandleFunc("/players/scores", getPlayerScore).Methods("GET")
	r.HandleFunc("/fantasy-scores", fantasyScoreCal).Methods("GET")
	r.HandleFunc("/cap-holders", capHolders).Methods("GET")

	fmt.Print("starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
