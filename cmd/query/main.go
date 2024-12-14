package main

import (
	"encoding/json"
	"fmt"

	"net/http"

	"database/sql"

	_ "github.com/lib/pq"
)

const connectionString = "host=localhost port=5432 user=postgres dbname=sandbox password=postgres"

type Handlers struct {
	db *sql.DB
}

func (h *Handlers) ServePost(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Name string `json:name`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = h.db.Exec("UPDATE query_table SET NAME = $1", data.Name)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) ServeGet(w http.ResponseWriter, r *http.Request) {
	var name string
	row := h.db.QueryRow("SELECT name FROM query_table LIMIT 1")
	err := row.Scan(&name)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	name = "Hello, " + name
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(name))
}

func main() {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	handlers := Handlers{db: db}

	http.HandleFunc("/get", handlers.ServeGet)
	http.HandleFunc("/post", handlers.ServePost)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
