package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
	"webserver/internal/app/server/router"
	repositories2 "webserver/internal/pkg/domain/repositories"
	"webserver/internal/pkg/domain/services"
	"webserver/internal/pkg/infrastructure/transactional"
)

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
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
