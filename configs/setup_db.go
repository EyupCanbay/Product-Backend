package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var timeOut = 5 * time.Second

func ConnectDB() *mongo.Client {

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(DatabaseEnv()))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to mongoDB")
	return client
}

var DB *mongo.Client = ConnectDB()

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := DB.Database("tesodev_product").Collection(collectionName)
	return collection
}
