package main

import (
	"net/http"
)

func main() {
	ConnectDB()
	InitTables()
	http.HandleFunc("/register", HandleRegister)
	http.HandleFunc("/players", AllPlayerHandler)
	http.HandleFunc("/players/{PlayerID}", PlayerHandle)
	http.ListenAndServe(":8080", nil)

}
