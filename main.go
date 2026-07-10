package main

import (
	"encoding/json"
	//"fmt"
	"net/http"
)

func GetTodo(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	json.NewEncoder(res).Encode(tasks)
}

type todo struct {
	Title string `json:"title"`
	Id    int    `json:"id"`
}

func AddTodo(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "method not allowed", http.StatusBadRequest)
		return
	}
	var new todo
	err := json.NewDecoder(req.Body).Decode(&new)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	new.Id = len(tasks) + 1
	tasks = append(tasks, new)

	res.Header().Add("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(new)
}

var tasks = []todo{
	{Id: 1, Title: "Mosab1"},
	{Id: 2, Title: "Mosab2"},
}

func main() {
	http.HandleFunc("/", GetTodo)
	http.HandleFunc("/add", AddTodo)
	http.ListenAndServe(":9000", nil)
}
