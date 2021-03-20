package main

import (
	"fmt"
	"log"

	"github.com/Haizza1/tronics/config"
	"github.com/Haizza1/tronics/handlers"
	"github.com/Haizza1/tronics/lib"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
	cfg        config.Properties
)

func init() {
	client, db, collection = lib.GetConnection()
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Unable to load configuration: %+v", err)
	}
}

func main() {
	e := echo.New()

	h := &handlers.ProductHandler{Col: collection}
	e.GET("/products/:id", h.GetProduct)
	e.GET("/products", h.GetProducts)
	e.POST("/products", h.CreateProducts)
	e.PUT("/products/:id", h.UpdateProduct)
	e.DELETE("/products/:id", h.DeleteProduct)

	e.Logger.Info("Listening on port %s:%s", cfg.Host, cfg.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)))
}
