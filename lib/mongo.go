package lib

import (
	"context"
	"fmt"
	"log"

	"github.com/Haizza1/tronics/config"
	"github.com/ilyakaznacheev/cleanenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (

	//CollectionAPI collection interface
	CollectionAPI interface {
		InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	}
)

// initialize connection to mongoDB
func Initialize() (*mongo.Client, *mongo.Database, *mongo.Collection) {
	var cfg config.Properties
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Configuration cannot be read: %+v", err)
	}

	connectURI := fmt.Sprintf("mongodb://%s:%s", cfg.DBHost, cfg.DBPort)
	c, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectURI))
	if err != nil {
		log.Fatalf("Unable to connect to mongoDB: %+v", err)
	}

	db := c.Database(cfg.DBName)
	col := db.Collection(cfg.CollectionName)
	return c, db, col
}
