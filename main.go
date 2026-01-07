package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	fmt.Print("hi")
	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	_, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}
