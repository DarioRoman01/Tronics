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
	cfg config.Properties
	col *mongo.Collection
)

func init() {
	_, _, col = lib.Initialize()
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Unable to load configuration: %+v", err)
	}
}

func main() {
	e := echo.New()
	h := &handlers.ProductHandler{Col: col}
	e.POST("/products", h.CreateProducts)
	e.Logger.Info("Listening on port %s:%s", cfg.Host, cfg.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)))
}
