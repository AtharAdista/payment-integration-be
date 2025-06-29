package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(host, port, user, password, dbname string) {

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error

	DB, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal("Failed to open DB: ", err)
	}

	err = DB.Ping()

	if err != nil {
		log.Fatal("Failed to ping DB:", err)
	}

	fmt.Println("Database connected!")
}

func GetDB() *sql.DB {
	if DB == nil {
		log.Fatal("Database connection is nil!")
	}

	return DB
}
