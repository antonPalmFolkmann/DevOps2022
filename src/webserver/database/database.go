package database

import (
	"fmt"
	"log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "Semester_4"
	dbname   = "postgres"
)

func GetConnection() *gorm.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open("postgres", psqlInfo)

	if err != nil {
		panic("failed to connect database")
	}

	log.Println("DB Connection established...")
	return db
}