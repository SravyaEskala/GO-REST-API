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
Match string `json:"match"`
Runs int `json:"runs"`
Wickets int `json:"wickets"`
}



type Player struct {
Name string `json:"name"`
ID int `json:"id"`
Team string `json:"team"`
Scores []Score `json:"scores"`
}



type OnlyPlayer struct {
Name string `json:"name"`
ID int `json:"id"`
Team string `json:"team"`
}



type DisplayPlayer struct {
Players []OnlyPlayer `json:"players"`
}



type OnlyScores struct {
ID int `json:"id"`
Scores []Score `json:"scores"`
}



type displayScores struct {
PlayerScores []OnlyScores `json:"playerscores"`
}



var playerservice []Player



var TempScoreData []Score



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



func getPlayerScore(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
var tempScores displayScores
var playerdetails OnlyScores
for _, item := range playerservice {
playerdetails.ID = item.ID
playerdetails.Scores = append(playerdetails.Scores, item.Scores...)
tempScores.PlayerScores = append(tempScores.PlayerScores, playerdetails)
}
json.NewEncoder(w).Encode(tempScores)
tempScores.PlayerScores = nil
}



func main() {
r := mux.NewRouter()
TempScoreData = append(TempScoreData, Score{Match: "1", Wickets: 2, Runs: 150})
playerservice = append(playerservice, Player{
ID: 1,
Name: "Virat",
Team: "RCB",
Scores: TempScoreData,
})
TempScoreData = nil
TempScoreData = append(TempScoreData, Score{Match: "1", Wickets: 3, Runs: 50})
playerservice = append(playerservice, Player{
ID: 7,
Name: "Dhoni",
Team: "CSK",
Scores: TempScoreData,
})
TempScoreData = nil
r.HandleFunc("/player", postPlayer).Methods("POST")
r.HandleFunc("/player/{id}/score", postPlayerScore).Methods("POST")
r.HandleFunc("/players", getPlayers).Methods("GET")
r.HandleFunc("/players/scores", getPlayerScore).Methods("GET")
fmt.Print("starting server at port 8000\n")
log.Fatal(http.ListenAndServe(":8000", r))
}