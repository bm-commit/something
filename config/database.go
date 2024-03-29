package config

import (
	"context"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectBD ...
func ConnectBD(dbUser, dbPass, dbHost string) *mongo.Client {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	connectionStr := "mongodb+srv://" + dbUser + ":" + dbPass + "@" + dbHost
	if strings.Contains(dbHost, "localhost") {
		connectionStr = "mongodb://" + dbUser + ":" + dbPass + "@" + dbHost
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionStr))
	defer cancel()
	if err != nil {
		log.Fatalf("Error connecting to DB: %s", err.Error())
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %s", err.Error())
	}
	return client
}

// CheckConnection ...
func CheckConnection(client *mongo.Client) int {
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		return 0
	}
	return 1
}
