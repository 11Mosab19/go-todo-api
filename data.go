package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var SecretKey = []byte("mosab1212")
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

func AddUser(New *UserData) error {
	var name string
	CheckUsernameQuery := "SELECT username FROM 'users' WHERE 'username'= ?;"
	err := db.QueryRow(CheckUsernameQuery, New.UserName).Scan(&name)
	if err == nil {
		return errors.New("already exists")
	}
	if err == sql.ErrNoRows {
		HashedPassword, err := bcrypt.GenerateFromPassword([]byte(New.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		AddQuery := "INSERT INTO 'users' ('username','role','hashed_password') VALUES (?,?,?);"
		_, err = db.Exec(AddQuery, New.UserName, New.Role, string(HashedPassword))
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func login(user LoginData) (UserData, error) {
	var role, name, hashed_password string
	var id int
	LookForUserQuery := "SELECT username,hashed_password FROM 'users' WHERE 'username' = ?;"
	err := db.QueryRow(LookForUserQuery, user.Username).Scan(&name, &hashed_password)
	if err == nil {
		err = bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(user.Password))
		if err != nil {
			return UserData{}, err
		}
		GetIdAndRole := "SELECT id,role FROM 'users' WHERE 'username' = ?;"
		err1 := db.QueryRow(GetIdAndRole, user.Username).Scan(&id, &role)
		if err1 != nil {
			return UserData{}, err1
		}
		return UserData{
			Id:       id,
			UserName: user.Username,
			Role:     role,
			Password: hashed_password,
		}, nil
	}
	if err == sql.ErrNoRows {
		return UserData{}, errors.New("invalid data")
	}
	return UserData{}, err
}

func GenerateToken(user UserData) (string, error) {
	claim := Token{
		UserId: user.Id,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	obj := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	key, err := obj.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	return key, nil
}

func VerifyToken(req *http.Request) (*Token, error) {
	if req.Header.Get("Authorization") == "" {
		return &Token{}, errors.New("Unauthorized")
	}
	if !strings.HasPrefix(req.Header.Get("Authorization"), "Bearer ") {
		return &Token{}, errors.New("Unauthorized")
	}
	claims := &Token{}
	TokenString := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	ReturnToken, err := jwt.ParseWithClaims(TokenString, claims, func(t *jwt.Token) (any, error) { return SecretKey, nil })
	if err != nil || !ReturnToken.Valid {
		return &Token{}, errors.New("Unauthorized")
	}
	return claims, nil
}
