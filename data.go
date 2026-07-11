package main

import (
	"errors"
	"slices"
	"sync"
)

var mx sync.Mutex

var Players = []PlayerData{
	{Name: "mosab", Age: 20, Id: 1},
	{Name: "messi", Age: 38, Id: 2},
}

func AddPlayer(NewPlayer *PlayerData) {
	mx.Lock()
	defer mx.Unlock()
	NewPlayer.Id = len(Players) + 1
	Players = append(Players, *NewPlayer)
}

func DeletePlayer(id int) error {
	for i, player := range Players {
		if player.Id == id {
			Players = slices.Delete(Players, i, i+1)
			return nil
		}
	}
	return errors.New("not found")
}

func GetPlayer(id int) (PlayerData, error) {
	for _, player := range Players {
		if player.Id == id {
			return player, nil
		}
	}
	return PlayerData{}, errors.New("Not found")
}

func UpdatePlayer(Update PlayerData) (PlayerData, error) {
	mx.Lock()
	defer mx.Unlock()
	for i, player := range Players {
		if player.Id == Update.Id {
			Players[i].Age = Update.Age
			Players[i].Name = Update.Name
			return Players[i], nil
		}
	}
	return PlayerData{}, errors.New("Not found")
}

func GetAllPlayers() []PlayerData {
	if len(Players) == 0 {
		return []PlayerData{}
	}
	return Players
}
