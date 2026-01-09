package databases

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"os"
)

func access_creds() string {
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")

	return fmt.Sprintf("mongodb://%s:%s", db_host, db_port)
}

var mongo_uri = access_creds()
var MongoClient *mongo.Client

func Client() {
	client, err := mongo.Connect(options.Client().ApplyURI(mongo_uri))
	if err != nil {
		panic(err)
	}

	MongoClient = client

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Successfully connected to bible_app!")
}
