package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func AllPlayerHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		res.Header().Set("Content-Type", "application/json")
		players := GetAllPlayers()
		err := json.NewEncoder(res).Encode(players)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

	case http.MethodPost:
		if req.Header.Get("Content-Type") != "application/json" {
			http.Error(res, "error", http.StatusInternalServerError)
			return
		}
		var NewPlayer PlayerData
		err := json.NewDecoder(req.Body).Decode(&NewPlayer)
		if err != nil {
			http.Error(res, "error", http.StatusBadRequest)
			return
		}
		if NewPlayer.Name == "" || NewPlayer.Age <= 0 {
			http.Error(res, "Error", http.StatusBadRequest)
			return
		}
		AddPlayer(&NewPlayer)

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
		player, err := GetPlayer(id)
		if err != nil {
			http.Error(res, "not found", http.StatusNotFound)
			return
		}
		err = json.NewEncoder(res).Encode(player)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

	case http.MethodDelete:
		err := DeletePlayer(id)
		if err != nil {
			http.Error(res, "Error", http.StatusNotFound)
			return
		}
		res.WriteHeader(http.StatusNoContent)

	case http.MethodPut:
		if req.Header.Get("Content-Type") != "application/json" {
			http.Error(res, "Error", http.StatusInternalServerError)
			return
		}
		var UpdatedPlayer PlayerData
		err := json.NewDecoder(req.Body).Decode(&UpdatedPlayer)

		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		if UpdatedPlayer.Name == "" || UpdatedPlayer.Age <= 0 {
			http.Error(res, "Error", http.StatusBadRequest)
			return
		}
		UpdatedPlayer.Id = id
		player, err := UpdatePlayer(UpdatedPlayer)
		if err != nil {
			http.Error(res, "Not found", http.StatusNotFound)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		err = json.NewEncoder(res).Encode(player)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}
}
