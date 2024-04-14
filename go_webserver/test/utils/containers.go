package utils

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func startMongoDBContainer(ctx context.Context) (mongoURI string, stopContainer func(), err error) {
	port := "27017/tcp"
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{port},
		WaitingFor:   wait.ForListeningPort(nat.Port(port)),
		Env:          map[string]string{"MONGO_INITDB_DATABASE": "test"},
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
	p, err := mongoContainer.MappedPort(ctx, nat.Port(port))
	if err != nil {
		stopContainer()
		return "", nil, fmt.Errorf("failed to get container port: %w", err)
	}

	mongoURI = fmt.Sprintf("mongodb://%s:%s/test", host, p.Port())
	return mongoURI, stopContainer, nil
}

func CreateMongoRuntime(ctx context.Context) (*mongo.Client, func()) {
	mongoCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)

	mongoURI, stopContainer, err := startMongoDBContainer(mongoCtx)
	if err != nil {
		cancel()
		log.Fatalf("Failed to start MongoDB container: %v", err)
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(mongoCtx, clientOptions)
	if err != nil {
		stopContainer()
		cancel()
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	cleanup := func() {
		client.Disconnect(mongoCtx)
		stopContainer()
		cancel()
	}
	return client, cleanup
}
