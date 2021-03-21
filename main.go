package main

import (
	"fmt"
	"log"

	"github.com/Haizza1/tronics/config"
	"github.com/Haizza1/tronics/handlers"
	"github.com/Haizza1/tronics/lib"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	client   *mongo.Client
	db       *mongo.Database
	userscol *mongo.Collection
	prodscol *mongo.Collection
	cfg      config.Properties
)

func init() {
	client, db, userscol, prodscol = lib.GetConnection()
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Unable to load configuration: %+v", err)
	}
}

func main() {
	// instance a new echo server and general middlewares
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339_nano} ${remote_ip} ${host} ${method} ${uri} ${user_agent}` +
			`${status} ${error} ${latency_human}` + "\n",
	}))

	// instance handlers ph(products handlers) uh(users handlers)
	ph := &handlers.ProductHandler{Col: prodscol}
	uh := &handlers.UserHandler{Col: userscol}

	e.GET("/products/:id", ph.GetProduct)
	e.GET("/products", ph.GetProducts)
	e.POST("/products", ph.CreateProducts, middleware.BodyLimit("1M"))
	e.PUT("/products/:id", ph.UpdateProduct, middleware.BodyLimit("1M"))
	e.DELETE("/products/:id", ph.DeleteProduct)

	e.POST("/users/signup", uh.CreateUser)

	e.Logger.Info("Listening on port %s:%s", cfg.Host, cfg.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)))
}
