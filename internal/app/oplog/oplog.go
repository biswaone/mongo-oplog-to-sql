package oplog

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OplogEntry struct {
	Timestamp time.Time              `json:"ts" bson:"ts"`
	Operation string                 `json:"op" bson:"op"`
	Namespace string                 `json:"ns" bson:"ns"`
	Doc       map[string]interface{} `json:"o" bson:"o"`
}

// go routine to fetch the oplogs
func GetOplog(ctx context.Context, client *mongo.Client, database string) (<-chan OplogEntry, error) {
	cursor, err := client.Database(database).Collection("oplog.rs").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	oplogCh := make(chan OplogEntry)

	go func() {
		defer close(oplogCh)
		for cursor.Next(ctx) {
			var entry OplogEntry
			if err := cursor.Decode(&entry); err != nil {
				log.Println("Error decoding oplog entry ", err)
				continue
			}
			oplogCh <- entry
		}
		if err := cursor.Err(); err != nil {
			log.Println("Change stream error ", err)
		}
	}()
	return oplogCh, err
}
