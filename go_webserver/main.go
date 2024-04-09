package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
	"webserver/repositories"
)

func main() {
	db, cleanup, ctx := createDatabase()
	accountCollection := db.Database("bank").Collection("account")
	transactionCollection := db.Database("bank").Collection("bank")
	defer cleanup()

	ar := repositories.CreateNewAccountRepositoryMongodb(ctx, accountCollection)
	r := createRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createDatabase() (*mongo.Client, func(), context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	cleanup := func() {
		cancel()
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}
	return client, cleanup, ctx
}
