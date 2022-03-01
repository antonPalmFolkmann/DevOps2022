package storage

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
)

// Hack for an array of maps in golang:
// https://stackoverflow.com/questions/47130003/how-can-i-declare-list-of-maps-in-golang
type M map[string]interface{}

var (
	// Configuration
	DATABASE = "../minitwit.db"
	Db      *sql.DB = ConnectDb()
	UserM    M
	Session map[string]string = make(map[string]string)
)

// ConnectDb returns a new connection to the database
func ConnectDb() *sql.DB {
	db, _ := sql.Open("sqlite3", DATABASE)
	return db
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
			queryString := "select * from user where user_id = ?"
			UserM = QueryDb(queryString, true, Session["user_id"])[0]
		}

		next.ServeHTTP(w, r)
	})
}

func LoginQuery(r *http.Request) []M {
	getMessageSQL := "SELECT * FROM user WHERE username = '" + r.Form["username"][0] + "'"
	log.Println("Query in login method: " + getMessageSQL)
	queryResult := QueryDb(getMessageSQL, true)
	return queryResult
}

// Closes the database again at the end of the request
func AfterRequest() {
	Db.Close()
}