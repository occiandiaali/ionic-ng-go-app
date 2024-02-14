package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Entry struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	About string `json:"about"`
}

func main() {
	//connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//create the table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS entries (id SERIAL PRIMARY KEY, title TEXT, about TEXT)")

	if err != nil {
		log.Fatal(err)
	}

	//create router
	router := mux.NewRouter()
	router.HandleFunc("/entries", getEntries(db)).Methods("GET")
	router.HandleFunc("/entries/{id}", getEntry(db)).Methods("GET")
	router.HandleFunc("/entries", createEntry(db)).Methods("POST")
	router.HandleFunc("/entries/{id}", updateEntry(db)).Methods("PUT")
	router.HandleFunc("/entries/{id}", deleteEntry(db)).Methods("DELETE")

	//start server
	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// get all entries
func getEntries(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM entries")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		entries := []Entry{}
		for rows.Next() {
			var e Entry
			if err := rows.Scan(&e.ID, &e.Title, &e.About); err != nil {
				log.Fatal(err)
			}
			entries = append(entries, e)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(entries)
	}
}

// get entry by id
func getEntry(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var e Entry
		err := db.QueryRow("SELECT * FROM entries WHERE id = $1", id).Scan(&e.ID, &e.Title, &e.About)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(e)
	}
}

// create entry
func createEntry(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var e Entry
		json.NewDecoder(r.Body).Decode(&e)

		err := db.QueryRow("INSERT INTO entries (title, about) VALUES ($1, $2) RETURNING id", e.Title, e.About).Scan(&e.ID)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(e)
	}
}

// update entry
func updateEntry(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var e Entry
		json.NewDecoder(r.Body).Decode(&e)

		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.Exec("UPDATE entries SET title = $1, about = $2 WHERE id = $3", e.Title, e.About, id)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(e)
	}
}

// delete entry
func deleteEntry(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var e Entry
		err := db.QueryRow("SELECT * FROM entries WHERE id = $1", id).Scan(&e.ID, &e.Title, &e.About)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("DELETE FROM entries WHERE id = $1", id)
			if err != nil {
				//todo : fix error handling
				w.WriteHeader(http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode("Entry deleted")
		}
	}
}
