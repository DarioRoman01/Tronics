package handlers

import (
	"context"
	"net/http"

	"github.com/Haizza1/tronics/lib"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// compare passwords to verify credentials
func isValidCredential(givenpwd, storepwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(storepwd), []byte(givenpwd)); err != nil {
		return false
	}

	return true
}

// insertUser handle users creation in the db and validation of unique fields
func insertUser(ctx context.Context, user User, collection lib.CollectionAPI) (interface{}, *echo.HTTPError) {
	var newUser User

	// check if username or email is already in use
	result := collection.FindOne(ctx, bson.M{"username": user.Username, "email": user.Email})
	err := result.Decode(&newUser)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Unable to decode new user")
	}

	if newUser.Email != "" || newUser.Username != "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "This username or email is already in use")
	}

	// hash the password for security
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Unable to hash the password")
	}

	// set hashed password and new object id to user
	user.Password = string(hashedPassword)
	id := primitive.NewObjectID()
	user.ID = id

	// insert user into the db if err means that a unique index already exist
	res, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "username or email already in use")
	}

	return res, nil
}

// Handle users authentication
func loginUser(ctx context.Context, reqUser User, collection lib.CollectionAPI) (User, *echo.HTTPError) {
	var user User

	// check if user exist
	result := collection.FindOne(ctx, bson.M{"username": reqUser.Username})
	err := result.Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments {
		return user, echo.NewHTTPError(http.StatusBadRequest, "Unable to parse request user")
	}

	if err == mongo.ErrNoDocuments {
		return user, echo.NewHTTPError(http.StatusBadRequest, "User does not exist")
	}

	// validate credentials
	if !isValidCredential(reqUser.Password, user.Password) {
		return user, echo.NewHTTPError(http.StatusBadRequest, "Invalid credentials")
	}

	return user, nil
}
