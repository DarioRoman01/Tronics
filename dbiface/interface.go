package dbiface

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	// CollectionAPI collection interface
	CollectionAPI interface {
		InsertOne(ctx context.Context, document interface{}, opt ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	}
)
