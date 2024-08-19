package utils

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"path/filepath"
	"time"
)

const port = "30001"

var networkName = ""

func startMongoDBContainer(ctx context.Context) (mongoURI string, stopContainer func(), err error) {

	newNetwork, err := network.New(ctx)
	if err != nil {
		log.Fatalf("Error while creating network: %s", err.Error())
	}
	networkName = newNetwork.Name
	log.Printf("Network Name: %s", networkName)

	req := testcontainers.ContainerRequest{
		Image:          "mongo:7.0.8",
		Name:           "mongo",
		ExposedPorts:   []string{fmt.Sprintf("%s:%s", port, port)},
		Cmd:            []string{"mongod", "--replSet", "rs0", "--bind_ip_all", "--port", port},
		WaitingFor:     wait.ForListeningPort(port),
		Networks:       []string{networkName},
		NetworkAliases: map[string][]string{networkName: {"mongo"}},
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

func StartMigrationsContainer(parentCtx context.Context, mongoURI string) error {
	ctx, cancel := context.WithTimeout(parentCtx, time.Minute*2)
	defer cancel()

	mainDir := "migrator"

	buildContext := testcontainers.FromDockerfile{
		Context:    filepath.Join("../../"),
		Dockerfile: "Dockerfile",
		BuildArgs: map[string]*string{
			"MAIN_DIR": &mainDir,
		},
	}

	req := testcontainers.ContainerRequest{
		FromDockerfile: buildContext,
		Env: map[string]string{
			"MONGO_URL": mongoURI,
		},
		WaitingFor:     wait.ForLog("Migrations completed").WithStartupTimeout(5 * time.Minute),
		Networks:       []string{networkName},
		NetworkAliases: map[string][]string{networkName: {"migrations"}},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return fmt.Errorf("failed to start container: %s", err)
	}

	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			fmt.Printf("failed to terminate container: %s\n", err)
		}
	}(container, ctx)

	fmt.Println("Migrations ran successfully")
	return nil
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
			bson.A{bson.D{{"_id", 0}, {"host", fmt.Sprintf("mongo:%s", port)}}},
		},
	}
	replicaCommandResult := client.Database("admin").
		RunCommand(mongoCtx, bson.D{{Key: "replSetInitiate", Value: config}})

	if replicaCommandResult.Err() != nil {
		log.Fatalf("Failed to create Replica Set: %v", replicaCommandResult.Err())
	}

	// connect to the replica set instead of the node we know that exists
	client.Disconnect(mongoCtx)
	newMongoURI := "mongodb://mongo:30001/?replicaSet=rs0"
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

const TestDatabaseName = "wallet"
const MongoURI = "mongodb://mongo:30001/?replicaSet=rs0"
