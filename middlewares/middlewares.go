package middlewares

import (
	"net/http"
	"strings"

	"github.com/Haizza1/tronics/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

var cfg config.Properties

// read configuration
func init() {
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Errorf("Configuration cannot be read: %+v", err)
	}
}

// LoggerMiddleware print logs in the console for information about request
func LoggerMiddleware() echo.MiddlewareFunc {
	logger := middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339_nano} ${remote_ip} ${host} ${method} ${uri} ${user_agent}` +
			`${status} ${error} ${latency_human}` + "\n",
	})

	return logger
}

// valodate token before any request
func JwtMiddleware() echo.MiddlewareFunc {
	jwtMiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(cfg.JwtTokenSecret),
		TokenLookup: "header:x-auth-token",
	})

	return jwtMiddleware
}

// check if user had admin status
func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get the token
		headerToken := c.Request().Header.Get("x-auth-token")
		token := strings.Split(headerToken, " ")[1]
		claims := jwt.MapClaims{}

		// parse the token
		_, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtTokenSecret), nil
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Unable to parse token")
		}

		// verify token authorization
		if !claims["authorized"].(bool) {
			return echo.NewHTTPError(http.StatusForbidden, "You not have permissions to perform this action")
		}

		return next(c)
	}
}
