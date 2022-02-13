package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Configuration
const (
	DATABASE = "../minitwit.db"
)

var (
	// Create our little application :)
	r       *mux.Router = mux.NewRouter()
	db      *sql.DB     = ConnectDb()
	session map[string]string
)

// ConnectDb returns a new connection to the database
func ConnectDb() *sql.DB {
	db, _ := sql.Open("sqlite3", DATABASE)
	return db
}

// InitDb creates the database tables
func InitDb() {
	db := ConnectDb()
	defer db.Close()

	query, _ := ioutil.ReadFile("../schema.sql")

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(string(query))
	stmt.Exec()
	tx.Commit()
}

// Hack for an array of maps in golang:
// https://stackoverflow.com/questions/47130003/how-can-i-declare-list-of-maps-in-golang
type M map[string]interface{}

// Queries the database and returns a list of maps
func QueryDb(query string, one bool, args ...interface{}) []M {
	rv := make([]M, 0)
	rows, _ := db.Query(query, args...)
	cols, _ := rows.Columns()
	for rows.Next() {
		// Solution for storing results in map adapted from: https://kylewbanks.com/blog/query-result-to-map-in-golang
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))

		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		_ = rows.Scan(columnPointers...)

		row := make(M)
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			row[colName] = *val
		}

		rv = append(rv, row)
	}

	if len(rv) == 0 {
		return nil
	} else if one {
		return rv[:1]
	} else {
		return rv
	}
}

// Registers a new message for the user.
func AddMessage(w http.ResponseWriter, r *http.Request) {
	if _, found := session["user_id"]; !found {
		log.Fatalln("Abort 401")
	}

	r.ParseForm()
	if _, found := r.Form["text"]; found {
		insertMessageSQL := "INSERT INTO message (author_id, text, pub_date, flagged) VALUES (%s,%s,%s,0)"
		statement, err := db.Prepare(insertMessageSQL) // Avoid SQL injections

		if err != nil {
			log.Fatalln(err.Error())
		}
		_, err = statement.Exec(session["user_id"], r.Form["text"], time.Now)
		if err != nil {
			log.Fatalln(err.Error())
		}
		return redirect(url(timeline))
	}
}

// Convenience method to look up the id for a username.
func GetUserid(username string) (int, error) {
	var usernameResult int
	// Query for a value based on a single row.
	if err := db.QueryRow("SELECT user_id from user where id = ?",
		username).Scan(&username); err != nil {
		if err == sql.ErrNoRows {
			return -1, fmt.Errorf("GetUserId %d: unknown username", username)
		}
		return -1, fmt.Errorf("GetUserId %d: %v", username)
	}
	return usernameResult, nil
}

func YourHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

func main() {
	r.HandleFunc("/", YourHandler)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8080", r))
}
