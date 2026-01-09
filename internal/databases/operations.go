package databases

// import (
// 	"database/sql"
// 	"fmt"
// 	"github.com/joho/godotenv"
// 	_ "github.com/lib/pq"
// 	"log"
// 	"net/url"
// 	"os"
// )
//
// type DB_conn struct {
// 	db *sql.DB
// }
//
// func Connect() *DB_conn {
// 	err := godotenv.Load("configs/.env")
// 	if err != nil {
// 		log.Fatalf("Error loading .env file")
// 	}
//
// 	db_user := os.Getenv("DB_USER")
// 	db_password := url.QueryEscape(os.Getenv("DB_PASSWORD"))
// 	db_host := os.Getenv("DB_HOST")
// 	db_name := os.Getenv("DB_NAME")
// 	db_SSLMode := os.Getenv("DB_SSLMODE")
//
// 	conn_str := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
// 		db_user, db_password, db_host, db_name, db_SSLMode)
//
// 	db, err := sql.Open("postgres", conn_str)
// 	if err != nil {
// 		log.Fatalf("Error connecting to bible_app: %v", err)
// 	}
//
// 	if err = db.Ping(); err != nil {
// 		log.Fatalf("Error pinging bible_app: %v", err)
// 	}
//
// 	return &DB_conn{db: db}
// }
//
// func (db *DB_conn) ReadChapter(trnsl_abbr string, trnsl_name string, chp_num int) {}
//
// func (db *DB_conn) ReadVerses(trnsl_abbr string, trnsl_name string, chp_num int) {}
