package handlers

import (
	"context"
	"net/http"

	"github.com/Haizza1/tronics/lib"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product describes an electronic product e.g. phone
type Product struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"product_name" bson:"product_name" validate:"required,max=30"`
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
	var products []Product
	c.Echo().Validator = &ProductValidator{validator: v}

	//parse request body
	if err := c.Bind(&products); err != nil {
		return c.JSON(http.StatusBadRequest, "unable to parse request body")
	}

	// validate request body
	for _, product := range products {
		if err := c.Validate(product); err != nil {
			return c.JSON(http.StatusBadRequest, "Unable validate request payload")
		}
	}

	// insert product into the db
	ids, err := insertProducts(context.Background(), products, p.Col)
	if err != nil {
		return c.JSON(err.Code, err.Message)
	}

	return c.JSON(http.StatusCreated, ids)
}

// Get 1 product by id
func (p *ProductHandler) GetProduct(c echo.Context) error {
	product, err := findProduct(context.Background(), c.Param("id"), p.Col)
	if err != nil {
		return c.JSON(err.Code, err.Message)
	}

	return c.JSON(http.StatusOK, product)
}

// Get all products
func (p *ProductHandler) GetProducts(c echo.Context) error {
	products, err := findProducts(context.Background(), c.QueryParams(), p.Col)
	if err != nil {
		return c.JSON(err.Code, err.Message)
	}
	return c.JSON(http.StatusOK, products)
}

// Handle Products update
func (p *ProductHandler) UpdateProduct(c echo.Context) error {
	product, err := updateProduct(context.Background(), c.Param("id"), c.Request().Body, p.Col)
	if err != nil {
		return c.JSON(err.Code, err.Message)
	}

	return c.JSON(http.StatusCreated, product)
}

// Handles products deletion form db
func (p *ProductHandler) DeleteProduct(c echo.Context) error {

	delIDS, err := deleteProduct(context.Background(), c.Param("id"), p.Col)
	if err != nil {
		return c.JSON(err.Code, err.Message)
	}

	return c.JSON(http.StatusOK, delIDS)
}
