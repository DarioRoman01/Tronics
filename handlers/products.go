package handlers

import (
	"context"
	"net/http"

	"github.com/Haizza1/tronics/lib"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product describes an electronic product e.g. phone
type Product struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"product_name" bson:"product_name" validate:"required,max=10"`
	Price       int                `json:"price" bson:"price" validate:"required,max=2000"`
	Currency    string             `json:"currency" bson:"currency" validate:"required,len=3"`
	Discount    int                `json:"discount,omitempty" bson:"discount,omitempty"`
	Vendor      string             `json:"vendor" bson:"vendor" valida:"required"`
	Accessories []string           `json:"accessories,omitempty" bson:"accessories,omitempty"`
	IsEssential bool               `json:"is_essential" bson:"is_essential"`
}

type ProductHandler struct {
	Col lib.CollectionAPI
}

//Create products on mongodb database
func (p *ProductHandler) CreateProducts(c echo.Context) error {
	var product Product
	c.Echo().Validator = &ProductValidator{validator: v}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, "unable to parse request body")
	}

	if err := c.Validate(product); err != nil {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}

	result, err := p.Col.InsertOne(context.Background(), product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Unable to insert in db")
	}
	return c.JSON(http.StatusCreated, result)
}

// Get 1 product by id
func (p *ProductHandler) GetProduct(c echo.Context) error {
	var product Product
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Unable to convert to objectID")
	}
	result := p.Col.FindOne(context.Background(), bson.M{"_id": id})
	if err := result.Decode(&product); err != nil {
		return c.JSON(http.StatusNotFound, "Product not found")
	}

	return c.JSON(http.StatusOK, product)
}

// Get all products
func (p *ProductHandler) GetProducts(c echo.Context) error {
	var products []Product

	cursor, err := p.Col.Find(context.Background(), bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Unable to load the data")
	}

	if err := cursor.All(context.Background(), &products); err != nil {
		return c.JSON(http.StatusInternalServerError, "Unable to load the data")
	}

	return c.JSON(http.StatusOK, products)

}

// Handle Prodcuts update
func (p *ProductHandler) UpdateProduct(c echo.Context) error {
	var product Product
	c.Echo().Validator = &ProductValidator{validator: v}

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Unable to convert to objectID")
	}

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, "unable to parse request body")
	}

	if err := c.Validate(product); err != nil {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}

	result, err := p.Col.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": product})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Unable to insert in db")
	}

	return c.JSON(http.StatusCreated, result)
}

// Handles products deletion form db
func (p *ProductHandler) DeleteProduct(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Unable to convert to objectID")
	}

	result, err := p.Col.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Unable to delete product")
	}

	return c.JSON(http.StatusOK, result)
}
