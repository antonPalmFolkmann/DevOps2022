package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Configuration
const (
	DATABASE = "../minitwit.db"
	PER_PAGE = 30
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
		http.Redirect(w, r, "http:localhost:8080/timeline", http.StatusFound)
	}
}

// Convenience method to look up the id for a username.
func GetUserId(username string) (*int, error) {
	var usernameResult int
	// Query for a value based on a single row.
	if err := db.QueryRow("SELECT user_id from user where id = ?", username).Scan(&username); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("GetUserId %s: unknown username", username)
		}
		return nil, fmt.Errorf("GetUserId %s failed", username)
	}
	return &usernameResult, nil
}

type timelineData struct {
	Title       string
	Request     *http.Request
	Messages    []M
	UserId      string
	User        M
	Followed    bool
	ProfileUser M
	PerPage     int
}

func timeline(w http.ResponseWriter, r *http.Request) {
	data := timelineData{
		Title:    "Public Timeline",
		Request:  r,
		Messages: QueryDb("select * from message", false),
		UserId:   "123123",
		PerPage:  30,
	}

	tmpl := parseTemplate("templates/timeline.html")
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Failed to render the template with err: %v", err)
	}
}

func userTimeline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	data := timelineData{
		Title:       "User Timeline",
		Request:     r,
		Messages:    QueryDb("select * from message limit 50", false),
		ProfileUser: QueryDb("select * from user where username = %s", true, username)[0],
		UserId:      "123123",
		PerPage:     30,
	}
	tmpl := parseTemplate("templates/timeline.html")
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Failed to render the template with err: %v", err)
	}
}

func PublicTimeline(w http.ResponseWriter, r *http.Request) {
	messageQuery := "select message.*, user.* from message, user where message.flagged = 0 and message.author_id = user.user_id order by message.pub_date desc limit 30"

	data := timelineData{
		Title:    "Public Timeline",
		Request:  r,
		Messages: QueryDb(messageQuery, false, PER_PAGE),
		PerPage:  PER_PAGE,
	}

	tmpl := parseTemplate("templates/timeline.html")
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Failed to render the template with err: %v", err)
	}
}

func parseTemplate(file string) *template.Template {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Failed to read the template contents: %v", err)
	}

	tmpl, err := template.New("timeline").Funcs(template.FuncMap{
		"gravatar": func(size int, email string) string { return GravatarUrl(email, size) },
	}).Parse(string(contents))
	if err != nil {
		log.Printf("Failed to parse the template: %v", err)
	}
	return tmpl
}

func YourHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

func main() {
	r.HandleFunc("/", YourHandler)
	r.HandleFunc("/timeline", timeline)
	r.HandleFunc("/public_timeline", PublicTimeline)
	// r.HandleFunc("/{username}", userTimeline)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Return the gravatar image for the given email address.
// Converting string to bytes: https://stackoverflow.com/questions/42541297/equivalent-of-pythons-encodeutf8-in-golang
// Converting bytes to hexadecimal string: https://pkg.go.dev/encoding/hex#EncodeToString
func GravatarUrl(email string, size int) string {
	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d",
		hex.EncodeToString([]byte(strings.ToLower(strings.TrimSpace(email)))), size)
}
