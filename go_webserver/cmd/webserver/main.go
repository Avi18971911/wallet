package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"webserver/internal/app/server/router"
	repositories2 "webserver/internal/pkg/domain/repositories"
	"webserver/internal/pkg/domain/services"
	"webserver/internal/pkg/infrastructure/transactional"
)

// @title Wallet API
// @version 1.0
// @description This is a simple wallet API
// termsOfService: http://swagger.io/terms/
// contact:
//   name: API Support
//   url: http://www.swagger.io/support
//   email: support@swagger.io

// license:
//   name: Apache 2.0
//   url: http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1

func main() {
	mongoURL := os.Getenv("MONGO_URL")
	cli, cleanup, ctx := createDatabase(mongoURL)
	accountCollection := cli.Database("wallet").Collection("account")
	transactionCollection := cli.Database("wallet").Collection("transaction")
	defer cleanup()

	ar := repositories2.CreateNewAccountRepositoryMongodb(accountCollection)
	tr := repositories2.CreateNewTransactionRepositoryMongodb(transactionCollection)

	tra := transactional.NewMongoTransactional(cli)

	as := services.CreateNewAccountServiceImpl(ar, tr, tra)
	ts := services.CreateNewTransactionServiceImpl(tr, ar, tra)
	r := router.CreateRouter(as, ts, ctx)
	log.Printf("Starting webserver")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createDatabase(mongoURL string) (*mongo.Client, func(), context.Context) {
	// TODO: Do not pass this context down, figure out another way
	ctx, cancel := context.WithCancel(context.Background())
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
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
