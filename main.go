package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
)

func AllPlayerHandler(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(res, "error", http.StatusInternalServerError)
		return
	}
	switch req.Method {
	case http.MethodGet:
		res.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(res).Encode(Players)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
		}

	case http.MethodPost:
		var NewPlayer PlayerData
		err := json.NewDecoder(req.Body).Decode(&NewPlayer)
		if err != nil {
			http.Error(res, "error", http.StatusBadRequest)
			return
		}
		NewPlayer.Id = len(Players) + 1
		Players = append(Players, NewPlayer)

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(NewPlayer)

	default:
		http.Error(res, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func PlayerHandle(res http.ResponseWriter, req *http.Request) {
	strId := req.PathValue("PlayerID")
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	switch req.Method {
	case http.MethodGet:
		res.Header().Set("Content-Type", "application/json")
		for _, player := range Players {
			if player.Id == id {
				err := json.NewEncoder(res).Encode(player)
				if err != nil {
					http.Error(res, err.Error(), http.StatusBadRequest)
					return
				}
			}
		}
		http.Error(res, "not found", http.StatusNotFound)
		return

	case http.MethodDelete:
		for i, player := range Players {
			if player.Id == id {
				Players = slices.Delete(Players, i, i+1)
				res.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.Error(res, "not found", http.StatusNotFound)
		return
	}
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
	http.HandleFunc("/players", AllPlayerHandler)
	http.HandleFunc("/players/{PlayerID}", PlayerHandle)
	http.ListenAndServe(":8080", nil)

}
