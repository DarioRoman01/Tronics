package handlers

import (
	"net/http"

	"github.com/Haizza1/tronics/lib"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type user struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username" validate:"required,min=3"`
	Email    string             `json:"email" bson:"email" validate:"required,min=7"`
	Password string             `json:"password" bson:"password" validate:"required,min=8"`
	IsAdmin  bool               `json:"is_admin,omitempty" bson:"is_admin"`
}

type UserHandler struct {
	Col lib.CollectionAPI
}

func (u *UserHandler) CreateUser(c echo.Context) error {
	return c.JSON(http.StatusCreated, "user created")
}
