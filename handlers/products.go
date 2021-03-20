package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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

//Create products on mongodb database
func (h *ProductHandler) CreateProducts(c echo.Context) error {
	var product Product
	c.Echo().Validator = &ProductValidator{validator: v}

	if err := c.Bind(&product); err != nil {
		log.Errorf("Unable to bind : %v", err)
		return c.JSON(http.StatusUnprocessableEntity, "unable to parse request payload")
	}

	if err := c.Validate(product); err != nil {
		log.Errorf("Unable to validate the product %+v %v", product, err)
		return c.JSON(http.StatusBadRequest, "unable to validate request payload")
	}

	id, httpError := h.InsertProduct(context.Background(), product, h.Col)
	if httpError != nil {
		return c.JSON(http.StatusInternalServerError, "unable to insert to database")
	}
	return c.JSON(http.StatusCreated, id)
}

// Get all products
func GetProducts(c echo.Context) error {
	return c.JSON(http.StatusOK, "listed products")
}

// Get 1 product by id
func GetProduct(c echo.Context) error {
	return c.JSON(http.StatusOK, "product retrieved")
}

// Handle Prodcuts update
func UpdateProduct(c echo.Context) error {
	return c.JSON(http.StatusOK, "product updated")
}

// Handles products deletion form db
func DeleteProduct(c echo.Context) error {
	return c.JSON(http.StatusOK, "Product deleted")
}
