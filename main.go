package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Data struct {
	Name string
}

func main() {
	// Load .env data put to os environment
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}

	/**
	* Connection variable definition
	**/
	// Connection string
	var mongoConnectionString = os.Getenv("MONGO_DB_CONNECTION_STRING")

	// Per variable definition
	var mongoHost = os.Getenv("MONGO_DB_HOST_NATIVE")
	var mongoUsername = os.Getenv("MONGO_DB_USERNAME")
	var mongoPassword = os.Getenv("MONGO_DB_PASSWORD")
	var mongoPort = os.Getenv("MONGO_DB_PORT")

	// If mongo db connection string is empty then create the connection from variabel
	if len(mongoConnectionString) == 0 {
		mongoConnectionString = "mongodb://" + mongoUsername + ":" + mongoPassword + "@" + mongoHost + ":" + mongoPort + "/?retryWrites=true&w=majority"
	}

	/**
	* Create mongo db connection by connection string
	**/
	// Context to pass the data
	ctx := context.TODO()
	// Connect to the mongodb
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConnectionString))
	if err != nil {
		panic(err)
	}
	// Defer function to disconnect to the mongo db
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	/**
	* Define query for mongo db. Define the db and collection
	**/
	var dbName = os.Getenv("MONGO_DB_NAME")
	var dbCollectionName = os.Getenv("MONGO_DB_COLLECTION_NAME")

	// Make connection to database and collection
	coll := client.Database(dbName).Collection(dbCollectionName)

	// CREATE OPERATION
	newData := Data{Name: "Test New Data"}
	result, err := coll.InsertOne(ctx, newData)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted Data", result)

	// READ OPERATION
	// UPDATE OPERATION
	// DELETE OPERATION

}
