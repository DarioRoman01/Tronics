package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/Haizza1/tronics/lib"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Products dao stands for products Data Access Object the goal of this functions is to
communicate the db with the handlers and separete the logic of the db from the handlers
*/

func findProducts(ctx context.Context, q url.Values, collection lib.CollectionAPI) ([]Product, *echo.HTTPError) {
	var products []Product

	filter := make(map[string]interface{})
	for k, v := range q {
		filter[k] = v
	}

	if filter["_id"] != nil {
		id, err := primitive.ObjectIDFromHex(filter["_id"].(string))
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, "Unable to convert to object id")
		}
		filter["_id"] = id
	}

	cursor, err := collection.Find(ctx, bson.M(filter))
	if err != nil {
		return products, echo.NewHTTPError(http.StatusNotFound, "product not found")
	}

	if err = cursor.All(ctx, &products); err != nil {
		return products, echo.NewHTTPError(http.StatusUnprocessableEntity, "Unable to load products")
	}

	return products, nil
}

func findProduct(ctx context.Context, id string, collection lib.CollectionAPI) (Product, *echo.HTTPError) {
	var product Product

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return product, echo.NewHTTPError(http.StatusInternalServerError, "Unable to convert to object id")
	}

	result := collection.FindOne(ctx, bson.M{"_id": docID})
	if err := result.Decode(&product); err != nil {
		return product, echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	return product, nil
}

func insertProducts(ctx context.Context, products []Product, collection lib.CollectionAPI) ([]interface{}, *echo.HTTPError) {
	var insertedIDS []interface{}

	for _, product := range products {
		product.ID = primitive.NewObjectID()
		result, err := collection.InsertOne(ctx, product)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, "Unable to insert into the db")
		}
		insertedIDS = append(insertedIDS, result.InsertedID)
	}
	return insertedIDS, nil
}

func updateProduct(ctx context.Context, id string, reqBody io.ReadCloser, collection lib.CollectionAPI) (Product, *echo.HTTPError) {
	var updatedProduct Product

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return updatedProduct, echo.NewHTTPError(http.StatusInternalServerError, "Unable to convert to object id")
	}

	filter := bson.M{"_id": docID}

	result := collection.FindOne(ctx, filter)
	if err := result.Decode(&updatedProduct); err != nil {
		return updatedProduct, echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	if err := json.NewDecoder(reqBody).Decode(&updatedProduct); err != nil {
		return updatedProduct, echo.NewHTTPError(http.StatusUnprocessableEntity, "Unable to parse request payload")
	}

	if err := v.Struct(updatedProduct); err != nil {
		return updatedProduct, echo.NewHTTPError(http.StatusBadRequest, "Bad request")
	}

	if _, err = collection.UpdateOne(ctx, filter, bson.M{"$set": updatedProduct}); err != nil {
		return updatedProduct, echo.NewHTTPError(http.StatusInternalServerError, "Unable to update the product")
	}

	return updatedProduct, nil
}

func deleteProduct(ctx context.Context, id string, collection lib.CollectionAPI) (int64, *echo.HTTPError) {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, "Unable to convert to object id")
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": docID})
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	return result.DeletedCount, nil
}
