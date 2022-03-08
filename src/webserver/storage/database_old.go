package storage

import (
	"database/sql"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Hack for an array of maps in golang:
// https://stackoverflow.com/questions/47130003/how-can-i-declare-list-of-maps-in-golang
type M map[string]interface{}

var (
	// Configuration
	DATABASE = "../minitwitcopy.db"
	Db      *sql.DB = ConnectDb()
	UserM    M
	Session  map[string]string = make(map[string]string)
)

const (
	PER_PAGE = 30
)

// ConnectDb returns a new connection to the database
func ConnectDb() *sql.DB {
	return ConnectPsql()
}

// InitDb creates the database tables
func InitDb() {
	defer Db.Close()

	query, _ := ioutil.ReadFile("../schema.sql")

	tx, _ := Db.Begin()
	stmt, _ := tx.Prepare(string(query))
	stmt.Exec()
	tx.Commit()
}

// Queries the database and returns a list of maps
func QueryDb(query string, one bool, args ...interface{}) []M {
	rv := make([]M, 0)

	stmt, _ := Db.Prepare(query)
	defer stmt.Close()

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
		Db = ConnectDb()
		UserM = nil
		if _, found := Session["user_id"]; found {
			queryString := "select * from \"user\" where user_id = ?"
			UserM = QueryDb(queryString, true, Session["user_id"])[0]
		}

		next.ServeHTTP(w, r)
	})
}

func LoginQuery(r *http.Request) []M {
	getMessageSQL := "SELECT * FROM \"user\" WHERE username = '" + r.Form["username"][0] + "'"
	log.Println("Query in login method: " + getMessageSQL)
	queryResult := QueryDb(getMessageSQL, true)
	return queryResult
}

func CreateUserQuery(r *http.Request, hash hash.Hash) error {
	insertMessageSQL := "INSERT INTO \"user\" (username, email, pw_hash) values (?, ?, ?)"
	stmt, _ := Db.Prepare(insertMessageSQL)
	defer stmt.Close()
	_, err := stmt.Exec(r.Form["username"][0], r.Form["email"][0], fmt.Sprintf("%x", hash.Sum(nil)))
	return err
}

func AddMessageQuery(r *http.Request) error {
	currentTime := int32(time.Now().Unix())
	insertMessageSQL := "INSERT INTO message (author_id, text, pub_date, flagged) VALUES (?,?,?,0)"
	statement, err := Db.Prepare(insertMessageSQL)

	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(Session["user_id"], r.Form["message"][0], currentTime)
	return err
}

func IsUsernameTaken(username string) bool {
	log.Printf("database.go/UserNameExistsInDB: looking for %s", username)
	UsernameQuery := "SELECT username FROM \"user\" WHERE username = ?"
	UsernameMap := QueryDb(UsernameQuery, true, username)

	return !(len(UsernameMap) == 0)
}

func CreateNewFollowingQuery(r *http.Request) error {
	insertMessageSQL := "INSERT INTO follower (who_id, whom_id) VALUES (?, ?)"
	statement, err := Db.Prepare(insertMessageSQL)

	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = statement.Exec(Session["user_id"], r.Form["text"], time.Now)
	return err
}

func DeleteFollowerQuery(r *http.Request) error {
	deleteMessageSQL := "DELETE FROM follower WHERE who_id = ? AND whom_id = ?"
	statement, err := Db.Prepare(deleteMessageSQL)

	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = statement.Exec(Session["user_id"], r.Form["text"], time.Now)
	return err
}

// Convenience method to look up the id for a username.
func GetUserId(username string) *int {
	messageQuery := "SELECT user_id FROM \"user\" WHERE username = '?'"
	usernameResult := QueryDb(messageQuery, false, username)
	if len(usernameResult) == 0 {
		return nil
	}
	userID := int(usernameResult[0]["user_id"].(int64))

	return &userID
}

func GetAllMessages() []M {
	getMessageQuery := "SELECT text from message"
	return QueryDb(getMessageQuery, false)
}

func GetAllNonFlaggedMessages() []M {
	getMessageQuery := "SELECT text from message where message.flagged = 0"
	return QueryDb(getMessageQuery, false)
}

func GetAllNonFlaggedMessagesFromUser(userString string) []M {
	userID := GetUserId(userString)
	getMessageQuery := "SELECT text from message where message.flagged = 0 and author_id = ?"
	return QueryDb(getMessageQuery, false, strconv.Itoa(*userID))
}

func Get30NonFlaggedMessagesFromTimeline(r *http.Request) []M {
	_ = r.URL.Query().Get("offset")

	messageQuery := "select message.*, \"user\".* from message, \"user\" where message.flagged = 0 and message.author_id = \"user\".user_id and ( \"user\".user_id = ? or \"user\".user_id in (select whom_id from follower where who_id = ?)) order by message.pub_date desc limit ?"
	messages := QueryDb(messageQuery, false, Session["user_id"], Session["user_id"], PER_PAGE)
	return messages
}

func Get30NonFlaggedMessagesFromPublicTimeline() []M {
	messageQuery := "select message.*, \"user\".* from message, \"user\" where message.flagged = 0 and message.author_id = \"user\".user_id order by message.pub_date desc limit 30"
	messages := QueryDb(messageQuery, false)
	return messages
}

func GetCurrentUserQuery(r *http.Request) M {
	vars := mux.Vars(r)
	username := vars["username"]
	log.Printf("HELLO! User is: %v", username)

	UserQuery := "SELECT * FROM \"user\" WHERE username = ?"
	ProfileUser := QueryDb(UserQuery, true, username)[0]
	log.Println(ProfileUser)
	return ProfileUser
}

func IsUserFollowed(UserMap *interface{}) bool {
	followed := false
	FollowerEmptyQuery := "select 1 from follower where follower.who_id = ? and follower.whom_id = ?"
	FollowerMap := QueryDb(FollowerEmptyQuery, true, Session["user_id"], UserMap)
	if UserM != nil {
		followed = len(FollowerMap) == 0
	}
	return followed
}

func Get30MessagesFromLoggedInUser(UserMap *interface{}) []M {
	MessagesFromLoggedInUserQuery := "select message.*, \"user\".* from message, \"user\" where \"user\".user_id = message.author_id and \"user\".user_id = ? order by message.pub_date desc limit ?"
	MessagesFromUserMap := QueryDb(MessagesFromLoggedInUserQuery, false, UserMap, PER_PAGE)
	return MessagesFromUserMap
}

// Closes the database again at the end of the request
func AfterRequest() {
	Db.Close()
}
