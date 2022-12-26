package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Data struct {
	ID   primitive.ObjectID `bson:"_id"`
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
	// Context to controll the process
	ctx := context.TODO()
	// ctx, cancel := context.WithTimeout(ctx, 10*time.Second) / If you are using context with timeout
	// Connect to the mongodb
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConnectionString))
	if err != nil {
		panic(err)
	}

	// Defer function to disconnect to the mongo db
	defer func() {
		// cancel() // If you are using context with timeout
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
	newData := Data{ID: primitive.NewObjectID(), Name: "Test New Data"}
	createResult, err := coll.InsertOne(ctx, newData)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted Data :", createResult)

	// READ OPERATION
	findData := Data{}
	filter := bson.D{{Key: "_id", Value: createResult.InsertedID}}
	err = coll.FindOne(context.TODO(), filter).Decode(&findData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			fmt.Println("Find Data :", "Not Found")
		} else {
			panic(err)
		}
	} else {
		fmt.Println("Find Data :", findData)
	}

	// UPDATE OPERATION
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "updated_field", Value: true}}}}
	updateResult, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Data Updated :", updateResult.ModifiedCount)
	}

	// DELETE OPERATION
	deleteResult, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Documents deleted: %d\n", deleteResult.DeletedCount)
	}

}
