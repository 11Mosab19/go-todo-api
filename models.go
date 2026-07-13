package main

import "github.com/golang-jwt/jwt/v5"

type PlayerData struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
	Age  int    `json:"age"`
}

type UserData struct {
	Id       int    `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type LoginResponse struct {
	Token string `json:"token"`
}
