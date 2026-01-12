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

type MClient struct {
	mc *mongo.Client
}

func Client() MClient {
	uri := access_creds()
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to bible_app!")
	return MClient{mc: client}
}

func (client *MClient) Close(ctx context.Context) error {
	return client.mc.Disconnect(ctx)
}
