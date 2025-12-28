package database

import (
	"context"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Service struct {
	Database *mongo.Database
}

var (
	dbURI  = os.Getenv("MONGODB_URI")
	dbName = os.Getenv("MONGODB_DATABASE")
)

func New() *Service {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(dbURI)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	return &Service{
		Database: client.Database(dbName),
	}
}
