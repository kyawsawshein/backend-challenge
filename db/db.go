package db

import (
	"context"
	"fmt"
	"log"
	"sync"

	"backend-challenge/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Connection *mongo.Database
var dbSync sync.Once
var collection *mongo.Collection


func GetConn(ctx context.Context) *mongo.Database {
	dbSync.Do(func() { Connection = CreateConn(ctx) })
	return Connection
}

func CreateConn(ctx context.Context, dsn ...string) *mongo.Database {
	var connStr string
	if len(dsn) == 1 {
		connStr = dsn[0]
	} else if len(dsn) == 0 {
		db_config := config.Cfg.DB
		connStr = fmt.Sprintf("mongodb://%s:%s@%s:%s", db_config.User, db_config.Pwd, db_config.Host, db_config.Port)
	} else {
		fmt.Errorf("Unabel to create connecetion!")
		return nil
	}
	clientOpts := options.Client().ApplyURI(connStr)
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}

	// Access the database and collection
	conn := client.Database("mdb")
	go func() {
		<-ctx.Done()
		if err := client.Disconnect(ctx); err != nil {
			fmt.Errorf("Unable to close DB connection", err)
			return
		}
		fmt.Printf("DB connection closed")
	}()
	return conn
}

func GetCollection(index string) (*mongo.Collection) {
	collection := Connection.Collection(index)
	return collection
}
