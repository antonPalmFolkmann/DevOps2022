package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectPsql() *sql.DB {
	// connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=localhost sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PASSWORD"))
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		"database",
		5432,
		os.Getenv("POSTGRES_DB"))
	log.Println("connStr: ", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("psql.go/ConnectPsql(): Failed to connect to PSQL: %s", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("psql.go/ConnectedPsql(): Failed to ping the database: %s", err)
	}

	log.Println("Successfully connected to database!")
	return db
}
