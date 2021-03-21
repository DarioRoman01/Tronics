package lib

import (
	"context"
	"fmt"
	"log"

	"github.com/Haizza1/tronics/config"
	"github.com/ilyakaznacheev/cleanenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// collection interface
type CollectionAPI interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

// GetConnection connect to the db and retrieve all the needed data
func GetConnection() (*mongo.Client, *mongo.Database, *mongo.Collection, *mongo.Collection) {
	// read env variables
	var cfg config.Properties
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Configuration cannot be read: %+v", err)
	}

	// set connection
	connectURI := fmt.Sprintf("mongodb://%s:%s", cfg.DBHost, cfg.DBPort)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectURI))
	if err != nil {
		log.Fatalf("Unable to connect to mongoDB: %+v", err)
	}

	// set db and collections that will be used
	db := client.Database(cfg.DBName)
	UsersCollection := db.Collection(cfg.UsersCollection)
	ProductsCollection := db.Collection(cfg.ProductsCollectoin)

	// create indexes for unique fields
	isUsernameIndexUnique := true
	usernameindexModel := mongo.IndexModel{
		Keys: bson.M{"username": 1},
		Options: &options.IndexOptions{
			Unique: &isUsernameIndexUnique,
		},
	}

	isEamilUnique := true
	emailIndexModel := mongo.IndexModel{
		Keys: bson.M{"email": 1},
		Options: &options.IndexOptions{
			Unique: &isEamilUnique,
		},
	}

	// create the indexes for users collection
	_, err = UsersCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{usernameindexModel, emailIndexModel})
	if err != nil {
		log.Fatalf("Unable to create index: %+v", err)
	}

	return client, db, UsersCollection, ProductsCollection
}
