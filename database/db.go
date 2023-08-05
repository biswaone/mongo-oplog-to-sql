package database

import (
	"context"
	"fmt"

	"github.com/biswaone/mongo-oplog-to-sql/config"
	"github.com/jackc/pgx/v5"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(ctx context.Context, cfg config.Config) *mongo.Client {
	clientOptions := options.Client().ApplyURI(cfg.MongoUri).SetDirect(true)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return nil
	}
	return client
}

func NewPostgresConn(ctx context.Context, cfg config.Config) *pgx.Conn {
	conn, err := pgx.Connect(ctx, cfg.PostgresUri)
	if err != nil {
		fmt.Println("Error connecting to Postgres:", err)
		return nil
	}
	return conn
}

func DisconnectPostgresConn(ctx context.Context, conn *pgx.Conn) {
	conn.Close(ctx)
}

func DisconnectMongoClient(ctx context.Context, client *mongo.Client) {
	client.Disconnect(ctx)
}
