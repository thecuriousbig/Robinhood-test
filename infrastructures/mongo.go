package infrastructure

import (
	"context"
	"log"
	"robinhood/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Get().Mongo.URI))
	if err != nil {
		log.Fatalf("failed to connect mongo: %s\n", err.Error())
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("failed to ping mongo: %s\n", err.Error())
	}

	return client
}
