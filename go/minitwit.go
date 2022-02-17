package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

const (
	PER_PAGE = 30
)

var (
	// Configuration
	DATABASE = "../minitwit.db"

	// Create our little application :)
	r       *mux.Router = mux.NewRouter()
	db      *sql.DB     = ConnectDb()
	user    M
	session map[string]string = make(map[string]string)
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
	log.Printf("Attempting query with: "+query, args...)
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

// Make sure that we are connected to teh database each request and look up the current user to that we know they're
// there
func BeforeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db = ConnectDb()
		user = nil
		if _, found := session["user_id"]; found {
			user = QueryDb("select * from user where user id = %s", true, session["user_id"])[0]
		}

		next.ServeHTTP(w, r)
	})
}

// Closes the database again at the end of the request
func AfterRequest() {
	db.Close()
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

type loginData struct {
	Request  *http.Request
	Username string
	User     interface{}
	Error    string
}

func Login(w http.ResponseWriter, r *http.Request) {
	if _, found := session["user_id"]; found {
		http.Redirect(w, r, "http:localhost:8080/timeline", http.StatusMultipleChoices)
		return
	}

	userError := ""
	if r.Method == "POST" {
		r.ParseForm()

		if _, found := r.Form["username"]; found {
			//We concatenate like this because variable assignment with % doesn't seem to work here
			getMessageSQL := "SELECT * FROM user WHERE username = '" + r.Form["username"][0] + "'"
			queryResult := QueryDb(getMessageSQL, true)

			hash := md5.New()
			io.WriteString(hash, r.Form["password"][0])
			formPwHash := fmt.Sprintf("%x", hash.Sum(nil))

			if queryResult == nil {
				userError = "Invalid username"
			} else if queryResult[0]["pw_hash"].(string) != formPwHash {
				userError = "Invalid password"
			} else {
				queryUserID := queryResult[0]["user_id"].(int64)
				session["user_id"] = strconv.Itoa(int(queryUserID))
				user = queryResult[0]
				http.Redirect(w, r, "http:localhost:8080/timeline", http.StatusFound)
				return
			}
		}
	}

	username := ""
	if len(r.Form["username"]) != 0 {
		username = r.Form["username"][0]
	}

	data := loginData{
		Request:  r,
		Username: username,
		User:     user,
		Error:    userError,
	}

	tmpl, err := initTemplate("login.html").ParseFiles("templates/layout.html", "templates/login.html")
	if err != nil {
		log.Printf("Failed to parse login template with err: %v", err)
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Failed to render login template with err: %v", err)
	}
}

type registerData struct {
	Request  *http.Request
	Username string
	Email    string
	User     interface{}
	Error    string
}

func Register(w http.ResponseWriter, r *http.Request) {
	registerError := ""
	if _, found := session["user_id"]; found {
		http.Redirect(w, r, "http:localhost:8080/timeline", http.StatusFound)
	}

	if r.Method == "POST" {
		r.ParseForm()
		log.Printf("%v", r.Form)
		if _, found := r.Form["username"]; !found {
			registerError = "Please enter a username"
		} else if _, found := r.Form["email"]; !found {
			registerError = "Please enter a valid e-mail address"
		} else if !strings.Contains(r.Form["email"][0], "@") {
			registerError = "Please enter a valid e-mail address"
		} else if _, found := r.Form["password"]; !found {
			registerError = "Please enter a password"
		} else if _, err := GetUserId(r.Form["username"][0]); err == nil {
			registerError = "Username already taken"
		} else {
			hash := md5.New()
			io.WriteString(hash, r.Form["password"][0])

			insertMessageSQL := "INSERT INTO user (username, email, pw_hash) values ('%s', '%s', '%x')"

			insertMessageSQL = fmt.Sprintf(insertMessageSQL, r.Form["username"][0], r.Form["email"][0], hash.Sum(nil))

			_, err = db.Exec(insertMessageSQL, r.Form["username"][0], r.Form["email"][0], hash.Sum(nil))
			if err != nil {
				log.Fatalln(err)
			}

			http.Redirect(w, r, "http:localhost:8080/timeline", http.StatusFound)
			return
		}
	}

	username := ""
	if len(r.Form["username"]) != 0 {
		username = r.Form["username"][0]
	}

	email := ""
	if len(r.Form["email"]) != 0 {
		email = r.Form["email"][0]
	}

	data := registerData{
		Request:  r,
		User:     user,
		Username: username,
		Email:    email,
		Error:    registerError,
	}

	tmpl, err := initTemplate("register.html").ParseFiles("templates/layout.html", "templates/register.html")
	if err != nil {
		log.Printf("Failed to parse the templates with err: %v", err)
	}

	err = tmpl.ExecuteTemplate(w, "register.html", data)
	if err != nil {
		log.Printf("Failed to render the template with err: %v", err)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	delete(session, "user_id")
	http.Redirect(w, r, "http:localhost:8080/public_timeline", http.StatusOK)
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	//TO-DO: This check needs to be changed in all methods using it.
	if _, found := session["user_id"]; !found {
		log.Fatalln("Abort 401")
	}

	r.ParseForm()
	if _, found := r.Form["text"]; found {
		//TO-DO: Again, from where are these variables piped
		insertMessageSQL := "INSERT INTO follower (who_id, whom_id) VALUES (%s, %s)"
		statement, err := db.Prepare(insertMessageSQL) // Avoid SQL injections

		if err != nil {
			log.Fatalln(err.Error())
		}

		_, err = statement.Exec(session["user_id"], r.Form["text"], time.Now)
		if err != nil {
			log.Fatalln(err.Error())
		}
		//TO-DO: I am imagnining the following url redirects to the followed users timeline
		http.Redirect(w, r, "http:localhost:8080/user_timeline/%s", http.StatusFound)
	}
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	if _, found := session["user_id"]; !found {

		log.Fatalln("Abort 401")
	}

	r.ParseForm()
	if _, found := r.Form["text"]; found {
		//TO-DO: Again, from where are these variables piped
		deleteMessageSQL := "DELETE FROM follower WHERE who_id = %s AND whom_id = %s"
		statement, err := db.Prepare(deleteMessageSQL) // Avoid SQL injections

		if err != nil {
			log.Fatalln(err.Error())
		}

		_, err = statement.Exec(session["user_id"], r.Form["text"], time.Now)
		if err != nil {
			log.Fatalln(err.Error())
		}
		//TO-DO: I am imagnining the following url redirects to the un-followed users timeline
		http.Redirect(w, r, "http:localhost:8080/user_timeline/%s", http.StatusFound)
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

// Shows a users timeline or if no user is logged in it will
// redirect to the public timeline.  This timeline shows the user's
// messages as well as all the messages of followed users.
func Timeline(w http.ResponseWriter, r *http.Request) {
	log.Printf("We got a vistor from %s", r.RemoteAddr)

	if user == nil {
		http.Redirect(w, r, "/public", http.StatusMultipleChoices)
		return
	}

	_ = r.URL.Query().Get("offset")

	messageQuery := "select message.*, user.* from message, user where message.flagged = 0 and message.author_id = user.user_id and ( user.user_id = %s or user.user_id in (select whom_id from follower where who_id = %s)) order by message.pub_date desc limit %s"

	data := timelineData{
		Title:    "Public Timeline",
		Request:  r,
		Messages: QueryDb(messageQuery, false, session["user_id"], session["user_id"], PER_PAGE),
		UserId:   "123123",
		PerPage:  PER_PAGE,
	}

	tmpl, err := initTemplate("timeline.html").ParseFiles("templates/layout.html", "templates/timeline.html")
	if err != nil {
		log.Printf("Failed to parse the templates with err: %v", err)
	}

	err = tmpl.ExecuteTemplate(w, "timeline.html", data)
	if err != nil {
		log.Printf("Failed to render the template with err: %v", err)
	}
}

// Displays the latest messages of all users.
func PublicTimeline(w http.ResponseWriter, r *http.Request) {
	messageQuery := "select message.*, user.* from message, user where message.flagged = 0 and message.author_id = user.user_id order by message.pub_date desc limit 30"

	data := timelineData{
		Title:    "Public Timeline",
		Request:  r,
		Messages: QueryDb(messageQuery, false),
		PerPage:  PER_PAGE,
	}

	tmpl, err := initTemplate("timeline.html").ParseFiles("templates/layout.html", "templates/timeline.html")
	if err != nil {
		log.Printf("Failed to parse the templates with err: %v", err)
	}

	err = tmpl.ExecuteTemplate(w, "timeline.html", data)
	if err != nil {
		log.Printf("Failed to render the template with err: %v", err)
	}
}

// Displays a user's tweets
func UserTimeline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	ProfileUser := QueryDb("select * from user where username = %s", true, username)[0]
	if ProfileUser == nil {
		w.Write([]byte("404 Not Found"))
	}

	followed := false
	if user != nil {
		followed = QueryDb("select 1 from follower where follower.who_id = %s and follower.whom_id = %s", true, session["user_id"], ProfileUser["user_id"])[0] != nil
	}

	data := timelineData{
		Title:       "User Timeline",
		Request:     r,
		Messages:    QueryDb("select * from message limit 50", false),
		ProfileUser: ProfileUser,
		Followed:    followed,
		PerPage:     PER_PAGE,
		User:        user,
	}

	tmpl, err := initTemplate("timeline.html").ParseFiles("templates/layout.html", "templates/timeline.html")
	if err != nil {
		log.Printf("Failed to parse the templates with err: %v", err)
	}

	err = tmpl.ExecuteTemplate(w, "timeline.html", data)
	if err != nil {
		log.Printf("Failed to render the template with err: %v", err)
	}
}

func initTemplate(name string) *template.Template {
	return template.New(name).Funcs(template.FuncMap{
		"gravatar": func(size int, email string) string { return GravatarUrl(email, size) },
	})
}

func ServeCSS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/style.css")
}

func YourHandler(w http.ResponseWriter, r *http.Request) {
	defer AfterRequest()
	w.Write([]byte("Gorilla!\n"))
}

func main() {
	r.Use(BeforeRequest)

	r.HandleFunc("/static/style.css", ServeCSS)

	r.HandleFunc("/public", PublicTimeline)
	r.HandleFunc("/user/{username}", UserTimeline)

	r.HandleFunc("/user/{username}/follow", FollowUser)
	r.HandleFunc("/user/{username}/unfollow", UnfollowUser)

	r.HandleFunc("/login", Login)
	r.HandleFunc("/register", Register)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Return the gravatar image for the given email address.
// Converting string to bytes: https://stackoverflow.com/questions/42541297/equivalent-of-pythons-encodeutf8-in-golang
// Converting bytes to hexadecimal s%}tring: https://pkg.go.dev/encoding/hex#EncodeToString
func GravatarUrl(email string, size int) string {
	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d",
		hex.EncodeToString([]byte(strings.ToLower(strings.TrimSpace(email)))), size)
}
