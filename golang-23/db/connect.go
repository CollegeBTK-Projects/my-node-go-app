package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context, url string) (*mongo.Client, error) {
	heart := &event.ServerMonitor{
		ServerHeartbeatSucceeded: func(e *event.ServerHeartbeatSucceededEvent) {
			log.Printf("Соединение в порядке: %s %v", e.ConnectionID, e.Duration)
		},
		ServerHeartbeatFailed: func(f *event.ServerHeartbeatFailedEvent) {
			log.Printf("Ошибка в соединении: %s %v", f.ConnectionID, f.Failure)
		},
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url).SetServerMonitor(heart).SetHeartbeatInterval(2*time.Second).SetServerSelectionTimeout(5*time.Second))
	if err != nil {
		panic(err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		client.Disconnect(ctx)
		return nil, err
	}

	return client, nil
}
