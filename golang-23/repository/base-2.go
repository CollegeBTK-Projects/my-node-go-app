package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetVersion(ctx context.Context, client *mongo.Client) (string, error) {
	var vrs bson.M
	err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "buildInfo", Value: 1}}).Decode(&vrs)
	if err != nil {
		return "", err
	}

	v, ok := vrs["version"].(string)
	if !ok {
		return "", err
	}
	return v, nil
}

func GetListDatabases(ctx context.Context, client *mongo.Client) ([]string, error) {
	dbs, err := client.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	return dbs, err
}

func GetListCollections(ctx context.Context, client *mongo.Client) ([]string, error) {
	colls, err := client.Database("group478").ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	return colls, nil
}
