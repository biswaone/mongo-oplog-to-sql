package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetDatabases(ctx context.Context, client *mongo.Client) ([]string, error) {
	//List all databases
	allDatabases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defaultDatabases := map[string]bool{
		"local":  true,
		"admin":  true,
		"config": true,
	}
	// only return the user defined databases
	var databases []string

	for _, database := range allDatabases {
		if !defaultDatabases[database] {
			databases = append(databases, database)
		}
	}
	return databases, nil

}

func GetCollections(database string, ctx context.Context, client *mongo.Client) ([]string, error) {
	collection, err := client.Database(database).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		log.Fatal("cannot get the collections ", err)
	}
	return collection, nil
}
