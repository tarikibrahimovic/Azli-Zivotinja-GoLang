package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		return nil, err
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB")
	return client, nil
}

var DB *mongo.Client
var err error

func init() {
	DB, err = ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Get a handle to the collection
	animalUsersCollection := GetCollection(DB, "Azil_Animal-Users")

	// Create unique index on AnimalId
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"animal_id": 1, // Index animal_id field in ascending order.
		},
		Options: options.Index().SetUnique(true),
	}

	_, err = animalUsersCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal(err)
	}
}

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("Aliz-Animal_Users").Collection(collectionName)
	return collection
}
