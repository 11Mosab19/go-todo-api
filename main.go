package main

import (
	"encoding/json"
	"net/http"
)

func GetData(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(Players)
}

func AddPlayer(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "method not allowed", http.StatusBadRequest)
		return
	}
	var NewPlayer PlayerData
	err := json.NewDecoder(req.Body).Decode(&NewPlayer)
	if err != nil {
		http.Error(res, "error", http.StatusBadRequest)
		return
	}
	NewPlayer.Id = len(Players)
	Players = append(Players, NewPlayer)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(NewPlayer)
}

type PlayerData struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
	Age  int    `json:"age"`
}

var Players = []PlayerData{
	{Name: "mosab", Age: 20, Id: 1},
	{Name: "messi", Age: 38, Id: 2},
}

func main() {
	http.HandleFunc("/players", GetData)
	http.HandleFunc("/players/add", AddPlayer)
	http.ListenAndServe(":8080", nil)
}
