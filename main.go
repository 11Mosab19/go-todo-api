package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/players", AllPlayerHandler)
	http.HandleFunc("/players/{PlayerID}", PlayerHandle)
	http.ListenAndServe(":8080", nil)

}
