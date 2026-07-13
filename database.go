package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var db *sql.DB // عشان افضل كونيكت مع الداتا بيز

func ConnectDB() error {
	db, err := sql.Open("sqlite", "players.db") // كونيكت هنا
	if err != nil {
		return err
	}
	err = db.Ping() // بتاكد ان كونيكت تم
	if err != nil {
		return err
	}
	return nil
}

func InitTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS players (
	"id" INTEGER,
	"name" TEXT NOT NULL,
	"age" INTEGER,
	PRIMARY KEY ("id")
	);`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	query2 := `
	CREATE TABLE IF NOT EXISTS users (
	"id" INTEGER,
	"username" TEXT UNIQUE NOT NULL,
	"role" TEXT NOT NULL DEFAULT "user",
	"hashed_password" TEXT NOT NULL,
	PRIMARY KEY ("id")
	);`
	_, err = db.Exec(query2)
	if err != nil {
		return err
	}
	query3 := `
	CREATE TABLE IF NOT EXISTS user_player (
	"user_id" INTEGER,
	"player_id" INTEGER,
	PRIMARY KEY ("user_id","player_id"),
	FOREIGN KEY ("user_id") REFERENCES "users"("id"),
    FOREIGN KEY("player_id") REFERENCES "players"("id")
	);`
	_, err = db.Exec(query3)
	if err != nil {
		return err
	}
	return nil
}
