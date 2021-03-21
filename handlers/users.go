package handlers

import (
	"context"
	"net/http"

	"github.com/Haizza1/tronics/lib"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username" validate:"required,min=3"`
	Email    string             `json:"email" bson:"email" validate:"required,email"`
	Password string             `json:"password" bson:"password" validate:"required,min=8"`
	IsAdmin  bool               `json:"is_admin,omitempty" bson:"is_admin"`
}

type UserHandler struct {
	Col lib.CollectionAPI
}

// Handle validation of create user request
func (u *UserHandler) CreateUser(c echo.Context) error {
	var user User
	c.Echo().Validator = &UserValidator{validator: v}

	// parse and validate request
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse request body")
	}

	if err := c.Validate(user); err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to validate request body")
	}

	// insert data in the db
	result, err := insertUser(context.Background(), user, u.Col)
	if err != nil {
		return c.JSON(err.Code, err.Message)
	}

	return c.JSON(http.StatusCreated, result)
}
