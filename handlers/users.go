package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Haizza1/tronics/config"
	"github.com/Haizza1/tronics/lib"
	"github.com/dgrijalva/jwt-go"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var prop config.Properties

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

func (u User) createToken() (string, error) {
	if err := cleanenv.ReadEnv(&prop); err != nil {
		log.Errorf("Configuration cannot be read: %+v", err)
	}

	// set token payload
	claims := jwt.MapClaims{}
	claims["authorized"] = u.IsAdmin
	claims["user_id"] = u.Username

	// set secret
	claims["exp"] = time.Now().Add(time.Minute * 90).UTC()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(prop.JwtTokenSecret))
	if err != nil {
		log.Errorf("Unable to create token: %+v", err)
	}

	return token, nil
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

func (u *UserHandler) LoginUser(c echo.Context) error {
	var user User

	// parse request body
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "Unable to parse request body")
	}

	// authenticate user
	logUser, err := loginUser(context.Background(), user, u.Col)
	if err != nil {
		return c.JSON(err.Code, err.Message)
	}

	// generate token
	token, jwterr := logUser.createToken()
	if jwterr != nil {
		return c.JSON(http.StatusInternalServerError, "Unable to create token")
	}

	// add token to the response
	c.Response().Header().Set("x-auth-token", "Bearer "+token)
	return c.JSON(http.StatusOK, logUser.Username)
}
