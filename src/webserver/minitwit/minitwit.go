package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
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

	stmt, _ := db.Prepare(query)
	defer stmt.Close()

	log.Printf("Attempting query with: %v", stmt)

	rows, _ := stmt.Query(args...)
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

// Make sure that we are connected to the database each request and look up the current user so that we know they're
// there
func BeforeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db = ConnectDb()
		user = nil
		if _, found := session["user_id"]; found {
			queryString := "select * from user where user_id = ?"
			user = QueryDb(queryString, true, session["user_id"])[0]
		}

		next.ServeHTTP(w, r)
	})
}

// Closes the database again at the end of the request
func AfterRequest() {
	db.Close()
}

type messageData struct {
	Request *http.Request
	Message string
	User    interface{}
	Error   string
}

// Registers a new message for the user.
func AddMessage(w http.ResponseWriter, r *http.Request) {
	defer AfterRequest()
	userError := ""
	/* if _, found := session["user_id"]; !found {
		log.Fatalln("Abort 401")
	} */

	r.ParseForm()
	if _, found := r.Form["message"]; found {
		currentTime := int32(time.Now().Unix())
		insertMessageSQL := "INSERT INTO message (author_id, text, pub_date, flagged) VALUES (?,?,?,0)"
		statement, err := db.Prepare(insertMessageSQL) // Avoid SQL injections

		if err != nil {
			log.Fatalln(err.Error())
		}
		_, err = statement.Exec(session["user_id"], r.Form["message"][0], currentTime)
		if err != nil {
			log.Fatalln(err.Error())
		}
		log.Printf("SHOULD FLASH: Your message was recorded")
		http.Redirect(w, r, "/public", http.StatusFound)
	}

	message := ""
	if len(r.Form["message"]) != 0 {
		message = r.Form["message"][0]
	}

	data := messageData{
		Request: r,
		Message: message,
		User:    user,
		Error:   userError,
	}

	tmpl, err := initTemplate("addmessage.html").ParseFiles("templates/layout.html", "templates/addmessage.html")
	if err != nil {
		log.Printf("Failed to parse login template with err: %v", err)
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Failed to render login template with err: %v", err)
	}

}

type loginData struct {
	Request  *http.Request
	Username string
	User     interface{}
	Error    string
}

func Login(w http.ResponseWriter, r *http.Request) {
	defer AfterRequest()
	if _, found := session["user_id"]; found {
		log.Printf("Session is: %v", session)
		http.Redirect(w, r, "/", http.StatusMultipleChoices)
		return
	}

	userError := ""
	if r.Method == "POST" {
		r.ParseForm()

		if _, found := r.Form["username"]; found {
			//We concatenate like this because variable assignment with % doesn't seem to work here
			getMessageSQL := "SELECT * FROM user WHERE username = '" + r.Form["username"][0] + "'"
			log.Println("Query in login method: " + getMessageSQL)
			queryResult := QueryDb(getMessageSQL, true)
			log.Println(queryResult)
			log.Printf("Query result: %v", queryResult)
			user = queryResult[0]

			hash := md5.New()
			io.WriteString(hash, r.Form["password"][0])
			formPwHash := fmt.Sprintf("%x", hash.Sum(nil))

			if queryResult == nil {
				userError = "Invalid username"
			} else if queryResult[0]["pw_hash"].(string) != formPwHash {
				userError = "Invalid password"
			} else {
				log.Printf("SHOULD FLASH: You were logged in")
				queryUserID := queryResult[0]["user_id"].(int64)
				session["user_id"] = strconv.Itoa(int(queryUserID))
				http.Redirect(w, r, "/", http.StatusFound)
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
	defer AfterRequest()
	if _, found := session["user_id"]; found {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	registerError := ""
	if r.Method == "POST" {
		r.ParseForm()
		if _, found := r.Form["username"]; !found {
			registerError = "Please enter a username"
		} else if _, found := r.Form["email"]; !found {
			registerError = "Please enter a valid e-mail address"
		} else if !strings.Contains(r.Form["email"][0], "@") {
			registerError = "Please enter a valid e-mail address"
		} else if _, found := r.Form["password"]; !found {
			registerError = "Please enter a password"
		} else if _, err := UserNameExistsInDB(r.Form["username"][0]); err != nil {
			registerError = "Username already taken"
		} else {
			hash := md5.New()
			io.WriteString(hash, r.Form["password"][0])

			insertMessageSQL := "INSERT INTO user (username, email, pw_hash) values (?, ?, ?)"
			stmt, _ := db.Prepare(insertMessageSQL)
			defer stmt.Close()
			_, err = stmt.Exec(r.Form["username"][0], r.Form["email"][0], fmt.Sprintf("%x", hash.Sum(nil)))

			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("SHOULD FLASH: You were successfully registered and can login now")
			http.Redirect(w, r, "/", http.StatusFound)
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

func UserNameExistsInDB(username string) (ok string, err error) {
	UsernameQuery := "SELECT username FROM user WHERE username = ?"
	UsernameMap := QueryDb(UsernameQuery, true, username)

	if len(UsernameMap) == 0 {
		return "okay", nil
	} else {
		return "error", errors.New("exists already")
	}

}

func Logout(w http.ResponseWriter, r *http.Request) {
	defer AfterRequest()
	log.Printf("SHOULD FLASH: You were logged out")
	delete(session, "user_id")
	http.Redirect(w, r, "/public", http.StatusOK)
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	defer AfterRequest()
	vars := mux.Vars(r)
	username := vars["username"]

	if _, found := session["user_id"]; !found {
		log.Fatalln("Abort 401")
	}

	r.ParseForm()
	if _, found := r.Form["text"]; found {
		insertMessageSQL := "INSERT INTO follower (who_id, whom_id) VALUES (?, ?)"
		statement, err := db.Prepare(insertMessageSQL)

		if err != nil {
			log.Fatalln(err.Error())
		}

		_, err = statement.Exec(session["user_id"], r.Form["text"], time.Now)
		if err != nil {
			log.Fatalln(err.Error())
		}

		redirectTo := fmt.Sprintf("/user/%s", username)
		http.Redirect(w, r, redirectTo, http.StatusFound)
	}
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	defer AfterRequest()
	vars := mux.Vars(r)
	username := vars["username"]

	if _, found := session["user_id"]; !found {
		log.Fatalln("Abort 401")
	}

	r.ParseForm()
	if _, found := r.Form["text"]; found {
		deleteMessageSQL := "DELETE FROM follower WHERE who_id = ? AND whom_id = ?"
		statement, err := db.Prepare(deleteMessageSQL) // Avoid SQL injections

		if err != nil {
			log.Fatalln(err.Error())
		}

		_, err = statement.Exec(session["user_id"], r.Form["text"], time.Now)
		if err != nil {
			log.Fatalln(err.Error())
		}

		redirectTo := fmt.Sprintf("/user/%s", username)
		http.Redirect(w, r, redirectTo, http.StatusFound)
	}
}

// Convenience method to look up the id for a username.
func GetUserId(username string) *int {
	messageQuery := fmt.Sprintf("SELECT user_id FROM user WHERE username = '%s'", username)
	usernameResult := QueryDb(messageQuery, false)
	if len(usernameResult) == 0 {
		return nil
	}
	userID := int(usernameResult[0]["user_id"].(int64))

	return &userID
}

func GetMessagesFromURL(url string) []M {
	var getMessageQuery string
	var resultMap []M
	split := strings.Split(url, "/")

	if split[1] == "public" {
		getMessageQuery = "SELECT text from message where message.flagged = 0"
		resultMap = QueryDb(getMessageQuery, false)
	} else if split[1] == "" {
		getMessageQuery = "SELECT text from message"
		resultMap = QueryDb(getMessageQuery, false)
	} else if split[1] == "user_timeline" {
		userID := GetUserId(split[2])
		getMessageQuery = "SELECT text from message where message.flagged = 0 and author_id = " + strconv.Itoa(*userID)
		resultMap = QueryDb(getMessageQuery, false)
	}

	/*

	 */
	return resultMap
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
	defer AfterRequest()
	log.Printf("We got a vistor from %s", r.RemoteAddr)
	log.Printf("User is: %v", user)

	if user == nil {
		http.Redirect(w, r, "/public", http.StatusMultipleChoices)
		return
	}

	log.Printf("User is: %v", user)

	_ = r.URL.Query().Get("offset")

	messageQuery := "select message.*, user.* from message, user where message.flagged = 0 and message.author_id = user.user_id and ( user.user_id = ? or user.user_id in (select whom_id from follower where who_id = ?)) order by message.pub_date desc limit ?"

	data := timelineData{
		Title:    "Public  Timeline",
		Request:  r,
		Messages: QueryDb(messageQuery, false, session["user_id"], session["user_id"], PER_PAGE),
		UserId:   session["user_id"],
		User:     user,
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
	defer AfterRequest()
	messageQuery := "select message.*, user.* from message, user where message.flagged = 0 and message.author_id = user.user_id order by message.pub_date desc limit 30"

	data := timelineData{
		Title:    "Public Timeline",
		Request:  r,
		Messages: QueryDb(messageQuery, false),
		User:     user,
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
	log.Println("User is:", username)

	ProfileUser := QueryDb("select * from user where username = ?", true, username)[0]
	log.Println(ProfileUser)
	if ProfileUser == nil {
		w.Write([]byte("404 Not Found"))
	}

	followed := false
	if user != nil {
		followed = len(QueryDb("select 1 from follower where follower.who_id = ? and follower.whom_id = ?", true, session["user_id"], ProfileUser["user_id"])) == 0
	}

	messages := QueryDb("select * from message limit 50", false)
	log.Println("messages: ", messages)

	data := timelineData{
		Title:       "User Timeline",
		Request:     r,
		Messages:    QueryDb("select message.*, user.* from message, user where user.user_id = message.author_id and user.user_id = ? order by message.pub_date desc limit ?", false, ProfileUser["user_id"], PER_PAGE),
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
		"gravatar":       func(size int, email interface{}) string { return GravatarUrl(email, size) },
		"datetimeformat": FormatDatetime,
	})
}

func ServeCSS(w http.ResponseWriter, r *http.Request) {
	defer AfterRequest()
	http.ServeFile(w, r, "static/style.css")
}

func YourHandler(w http.ResponseWriter, r *http.Request) {
	defer AfterRequest()
	w.Write([]byte("Gorilla!\n"))
}

func FormatDatetime(timestamp int64) string {
	timeUnix := time.Unix(timestamp, 0)
	return timeUnix.Format("2006-01-02 15:04")
}

func main() {
	HandleFuncRoutesMain()

	go ApiMain()

	log.Fatal(http.ListenAndServe(":8080", r))
}

func HandleFuncRoutesMain() {
	r.Use(BeforeRequest)

	r.HandleFunc("/static/style.css", ServeCSS)

	r.HandleFunc("/", Timeline)
	r.HandleFunc("/public", PublicTimeline)
	r.HandleFunc("/user/{username}", UserTimeline)

	r.HandleFunc("/user/{username}/follow", FollowUser)
	r.HandleFunc("/user/{username}/unfollow", UnfollowUser)
	r.HandleFunc("/addmessage", AddMessage)

	r.HandleFunc("/login", Login)
	r.HandleFunc("/logout", Logout)
	r.HandleFunc("/register", Register)
}

// Return the gravatar image for the given email address.
// Converting string to bytes: https://stackoverflow.com/questions/42541297/equivalent-of-pythons-encodeutf8-in-golang
// Converting bytes to hexadecimal s%}tring: https://pkg.go.dev/encoding/hex#EncodeToString
func GravatarUrl(email interface{}, size int) string {
	strEmail := email.(string)
	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d",
		hex.EncodeToString([]byte(strings.ToLower(strings.TrimSpace(strEmail)))), size)
}
