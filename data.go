package main

import (
	"errors"
	"fmt"
	"sync"
)

var mx sync.Mutex

func AddPlayer(NewPlayer *PlayerData) {
	mx.Lock()
	defer mx.Unlock()
	InsertQuery := "INSERT INTO 'players' ('name','age') VALUES (?,?)"
	db.Exec(InsertQuery, NewPlayer.Name, NewPlayer.Age)
}

func DeletePlayer(id int) error {
	DeleteQuery := "DELETE FROM 'players' WHERE 'id' = ?"
	result, err := db.Exec(DeleteQuery, id)
	if err != nil {
		return err
	}
	if x, _ := result.RowsAffected(); x == 0 {
		return errors.New("not found")
	}
	return nil
}

func GetPlayer(id int) (PlayerData, error) {
	var player PlayerData
	query := "SELECT * FROM 'players' WHERE 'id'= ?;"
	err := db.QueryRow(query, id).Scan(&player.Id, &player.Name, &player.Age)
	if err != nil {
		return PlayerData{}, err
	}
	return player, nil
}

func UpdatePlayer(Update PlayerData) (PlayerData, error) {
	var player PlayerData
	mx.Lock()
	defer mx.Unlock()
	UpdateQuery := "Update 'players' set 'name'= ? , 'age'= ? WHERE 'id' = ?;"
	result, err := db.Exec(UpdateQuery, Update.Name, Update.Age, Update.Id)
	if x, _ := result.RowsAffected(); x == 0 && err == nil {
		return PlayerData{}, errors.New("Not found")
	}
	if err != nil {
		return PlayerData{}, err
	}
	player, err = GetPlayer(Update.Id)
	if err != nil {
		return PlayerData{}, err
	}
	return player, nil
}

func GetAllPlayers() ([]PlayerData, error) {
	var Players = []PlayerData{}

	rows, err := db.Query("SELECT * FROM 'players';")
	if err != nil {
		fmt.Println(err)
		return []PlayerData{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var player PlayerData
		err := rows.Scan(&player.Id, &player.Name, &player.Age)
		if err != nil {
			return []PlayerData{}, err
		}
		Players = append(Players, player)
	}
	if err := rows.Err(); err != nil {
		return []PlayerData{}, err
	}
	return Players, nil
}
