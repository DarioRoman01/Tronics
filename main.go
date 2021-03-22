package main

import (
	"fmt"
	"log"

	"github.com/Haizza1/tronics/config"
	"github.com/Haizza1/tronics/handlers"
	"github.com/Haizza1/tronics/lib"
	"github.com/Haizza1/tronics/middlewares"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	userscol *mongo.Collection
	prodscol *mongo.Collection
	cfg      config.Properties
)

func init() {
	userscol, prodscol = lib.GetConnection()
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Unable to load configuration: %+v", err)
	}
}

func main() {
	// instance a new echo server and general middlewares
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middlewares.LoggerMiddleware())

	// instance handlers ph(products handlers) uh(users handlers)
	ph := &handlers.ProductHandler{Col: prodscol}
	uh := &handlers.UserHandler{Col: userscol}

	// products endpoints
	e.GET("/products/:id", ph.GetProduct)
	e.GET("/products", ph.GetProducts)
	e.POST("/products", ph.CreateProducts, middleware.BodyLimit("1M"), middlewares.JwtMiddleware())
	e.PUT("/products/:id", ph.UpdateProduct, middleware.BodyLimit("1M"), middlewares.JwtMiddleware())
	e.DELETE("/products/:id", ph.DeleteProduct, middlewares.JwtMiddleware(), middlewares.AdminMiddleware)

	// users endpoints
	e.POST("/users/signup", uh.CreateUser)
	e.POST("/users/login", uh.LoginUser)

	// start the server
	e.Logger.Info("Listening on port %s:%s", cfg.Host, cfg.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)))
}
