package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(ctx context.Context, url string) (*mongo.Client, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	ctxPing, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := client.Ping(ctxPing, nil); err != nil {
		client.Disconnect(ctx)
		return nil, err
	}
	defer cancel()
	return client, nil
}
