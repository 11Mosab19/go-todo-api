package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func AllPlayerHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		res.Header().Set("Content-Type", "application/json")
		erro, players := GetAllPlayers()
		if erro != nil {
			fmt.Println(erro)
			return
		}
		err := json.NewEncoder(res).Encode(players)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

	case http.MethodPost:
		claims, err := VerifyToken(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusUnauthorized)
			return
		}
		if claims.Role != "admin" {
			http.Error(res, "Not allowed", http.StatusForbidden)
			return
		}
		if req.Header.Get("Content-Type") != "application/json" {
			http.Error(res, "error", http.StatusInternalServerError)
			return
		}
		var NewPlayer PlayerData
		err = json.NewDecoder(req.Body).Decode(&NewPlayer)
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
		claims, err := VerifyToken(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusUnauthorized)
			return
		}
		if claims.Role != "admin" {
			http.Error(res, "Not allowed", http.StatusForbidden)
			return
		}
		err = DeletePlayer(id)
		if err != nil {
			http.Error(res, "Error", http.StatusNotFound)
			return
		}
		res.WriteHeader(http.StatusNoContent)

	case http.MethodPut:
		claims, err := VerifyToken(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusUnauthorized)
			return
		}
		if claims.Role != "admin" {
			http.Error(res, "Not allowed", http.StatusForbidden)
			return
		}
		if req.Header.Get("Content-Type") != "application/json" {
			http.Error(res, "Error", http.StatusInternalServerError)
			return
		}
		var UpdatedPlayer PlayerData
		err = json.NewDecoder(req.Body).Decode(&UpdatedPlayer)

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

func HandleRegister(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(res, "Error", http.StatusInternalServerError)
	}
	var NewUser UserData
	err := json.NewDecoder(req.Body).Decode(&NewUser)
	if err != nil {
		http.Error(res, "Error", http.StatusBadRequest)
		return
	}
	err = AddUser(&NewUser)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
}

func HandleLogin(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(res, "Error", http.StatusInternalServerError)
		return
	}
	var user LoginData
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(res, "Error", http.StatusInternalServerError)
		return
	}
	LoggedUser, err := login(user)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}
	token, err := GenerateToken(LoggedUser)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(LoginResponse{
		Token: token,
	})
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
