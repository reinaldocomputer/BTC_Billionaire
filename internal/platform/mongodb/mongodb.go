package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// Connection URI
const (
	uri        = "mongodb://admin:admin@db:27017"
	database   = "btc"
	collection = "transactions"
)

type MongoDB struct{}

var client *mongo.Client

func Connect() (err error) {
	// Create a new client and connect to the server
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		return err
	}
	// Ping the primary
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("MongoDB Successfully connected and pinged.")
	return nil
}

func (*MongoDB) Store(data interface{}) (err error) {
	coll := client.Database(database).Collection(collection)
	result, err := coll.InsertOne(context.TODO(), data)
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	return err
}

func (*MongoDB) Find(start time.Time, end time.Time) ([]bson.M, error) {
	coll := client.Database(database).Collection(collection)
	filterCursor, err := coll.Find(context.TODO(), bson.D{
		{"datetime",
			bson.D{
				{"$gt", start},
				{"$lt", end},
			},
		},
	},
		options.Find().SetProjection(bson.D{{"_id", 0}}))
	if err != nil {
		return nil, err
	}
	var result []bson.M
	if err = filterCursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}
	fmt.Println(result)
	return result, nil
}
