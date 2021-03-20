package handlers

import (
	"context"
	"net/http"

	"github.com/Haizza1/tronics/lib"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductHandler struct {
	Col lib.CollectionAPI
}

func (p *ProductHandler) InsertProduct(ctx context.Context, product Product, collection lib.CollectionAPI) (interface{}, error) {
	product.ID = primitive.NewObjectID()
	insertedID, err := collection.InsertOne(ctx, product)

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "unable to insert to database")
	}
	return insertedID, nil
}
