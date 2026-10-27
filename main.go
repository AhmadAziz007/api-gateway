package main

import (
	"ecommerce/api-gateway/handlers"
	"ecommerce/api-gateway/middleware/cors"
	gatewayjwt "ecommerce/api-gateway/middleware/jwt"
	"ecommerce/api-gateway/middleware/logger"
	"ecommerce/api-gateway/middleware/ratelimit"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.JSONFormatter{})

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(logger.Middleware())
	e.Use(cors.Middleware())
	e.Use(ratelimit.RedisMiddleware())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	api := e.Group("/api/v1")
	handlers.RegisterAllRoutes(api, gatewayjwt.Middleware())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	logrusLogger.Infof("API Gateway starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		logrusLogger.Fatal(err)
	}
}
