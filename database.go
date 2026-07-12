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
	"name" TEXT,
	"age" INTEGER,
	PRIMARY KEY ("id")
	);`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
