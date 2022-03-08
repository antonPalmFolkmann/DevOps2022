package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func ConnectPsql() *gorm.DB {
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		"database",
		5432,
		os.Getenv("POSTGRES_DB"))
	db, err := gorm.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("psql.go/ConnectPsql(): Failed to connect to PSQL: %s", err)
	}

	return db
}
