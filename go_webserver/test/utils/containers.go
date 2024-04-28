package utils

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const port = "30001"

func startMongoDBContainer(ctx context.Context) (mongoURI string, stopContainer func(), err error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:7.0.8",
		Name:         "mongo",
		ExposedPorts: []string{fmt.Sprintf("%s:%s", port, port)},
		Cmd:          []string{"mongod", "--replSet", "rs0", "--bind_ip_all", "--port", port},
		WaitingFor:   wait.ForListeningPort(port),
		// Env:          map[string]string{"MONGO_INITDB_DATABASE": "test"},
	}

	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to start container: %w", err)
	}

	stopContainer = func() {
		mongoContainer.Terminate(ctx)
	}

	// Get the container IP
	host, err := mongoContainer.Host(ctx)
	if err != nil {
		stopContainer()
		return "", nil, fmt.Errorf("failed to get container host: %w", err)
	}

	// Get the mapped port
	p, err := mongoContainer.MappedPort(ctx, port)
	if err != nil {
		stopContainer()
		return "", nil, fmt.Errorf("failed to get container port: %w", err)
	}

	mongoURI = fmt.Sprintf("mongodb://%s:%s", host, p.Port())
	log.Printf(mongoURI)
	return mongoURI, stopContainer, nil
}

func CreateMongoRuntime(ctx context.Context) (*mongo.Client, func()) {
	mongoCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)

	mongoURI, stopContainer, err := startMongoDBContainer(mongoCtx)
	if err != nil {
		cancel()
		log.Fatalf("Failed to start MongoDB container: %v", err)
	}

	clientOptions := options.Client().ApplyURI(mongoURI).SetDirect(true)
	client, err := mongo.Connect(mongoCtx, clientOptions)
	if err != nil {
		stopContainer()
		cancel()
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	config := bson.D{
		{"_id", "rs0"},
		{"members",
			bson.A{bson.D{{"_id", 0}, {"host", fmt.Sprintf("localhost:%s", port)}}},
		},
	}
	replicaCommandResult := client.Database("admin").
		RunCommand(mongoCtx, bson.D{{Key: "replSetInitiate", Value: config}})

	if replicaCommandResult.Err() != nil {
		log.Fatalf("Failed to create Replica Set: %v", replicaCommandResult.Err())
	}

	// connect to the replica set instead of the node we know that exists
	client.Disconnect(mongoCtx)
	newMongoURI := mongoURI + "/?replicaSet=rs0"
	newClientOptions := options.Client().ApplyURI(newMongoURI)
	newClient, newErr := mongo.Connect(mongoCtx, newClientOptions)
	if newErr != nil {
		stopContainer()
		cancel()
		log.Fatalf("Failed to connect to MongoDB: %v", newErr)
	}

	cleanup := func() {
		newClient.Disconnect(mongoCtx)
		stopContainer()
		cancel()
	}
	return newClient, cleanup
}

func CleanupDatabase(client *mongo.Client, ctx context.Context) error {
	collections, err := client.Database(TestDatabaseName).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return err
	}
	for _, collection := range collections {
		err := client.Database(TestDatabaseName).Collection(collection).Drop(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

const TestDatabaseName = "test"
