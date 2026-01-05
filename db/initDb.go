package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	log.Println("db opening sqlite database...")
	var err error

	DB, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal("db open error:", err)
	}

	log.Println("db ping")
	if err := DB.Ping(); err != nil {
		log.Fatal("db ping error:", err)
	}

	log.Println("db enabling foreign keys...")
	if _, err := DB.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		log.Fatal("db pragma error:", err)
	}

	log.Println("db ready")
}
