package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/url"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db_user := os.Getenv("DB_USER")
	db_password := url.QueryEscape(os.Getenv("DB_PASSWORD"))
	db_host := os.Getenv("DB_HOST")
	db_name := os.Getenv("DB_NAME")
	db_SSLMode := os.Getenv("DB_SSLMODE")

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		db_user, db_password, db_host, db_name, db_SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to bible_app: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging bible_app: %v", err)
	}

	fmt.Println("Successfully connected to bible_app!")
}
