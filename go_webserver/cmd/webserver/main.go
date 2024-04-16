package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
	"webserver/internal/app/server/router"
	repositories2 "webserver/internal/pkg/domain/repositories"
	"webserver/internal/pkg/domain/services"
	"webserver/internal/pkg/infrastructure/transactional"
)

func main() {
	cli, cleanup, ctx := createDatabase()
	accountCollection := cli.Database("bank").Collection("account")
	transactionCollection := cli.Database("bank").Collection("transaction")
	defer cleanup()

	ar := repositories2.CreateNewAccountRepositoryMongodb(accountCollection)
	tr := repositories2.CreateNewTransactionRepositoryMongodb(transactionCollection)

	tra := transactional.NewMongoTransactional(cli)

	as := services.CreateNewAccountServiceImpl(ar, tr, tra)
	ts := services.CreateNewTransactionServiceImpl(tr, ar, tra)
	r := router.CreateRouter(as, ts, ctx)
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
