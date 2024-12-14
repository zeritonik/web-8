package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

const connectionString = "host=localhost port=5432 user=postgres dbname=sandbox password=postgres"

type Handlers struct {
	db *sql.DB
}

func (h *Handlers) ServeGet(w http.ResponseWriter, r *http.Request) {
	var count int
	row := h.db.QueryRow("SELECT count FROM count_table LIMIT 1")
	err := row.Scan(&count)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	err = json.NewEncoder(w).Encode(struct {
		Count int `json:"count"`
	}{Count: count})
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) ServePost(w http.ResponseWriter, r *http.Request) {
	var dcount struct {
		Count int `json:"count"`
	}

	err := json.NewDecoder(r.Body).Decode(&dcount)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = h.db.Exec("UPDATE count_table SET count = count + $1", dcount.Count)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
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
