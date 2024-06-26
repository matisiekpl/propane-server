package main

import (
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/matisiekpl/propane-server/internal/client"
	"github.com/matisiekpl/propane-server/internal/controller"
	"github.com/matisiekpl/propane-server/internal/dto"
	"github.com/matisiekpl/propane-server/internal/repository"
	"github.com/matisiekpl/propane-server/internal/service"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Info("Error loading .env file")
	}

	config := dto.Config{DSN: os.Getenv("DSN"), SigningSecret: os.Getenv("JWT_SECRET")}

	db, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
	if err != nil {
		logrus.Panic(err)
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.CORS())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				var appError dto.AppError
				switch {
				case errors.As(err, &appError):
					return echo.NewHTTPError(400, err.Error())
				}
			}
			return err
		}
	})

	opt := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logrus.Panic(err)
	}
	m, err := app.Messaging(context.Background())
	if err != nil {
		logrus.Panic(err)
	}

	clients := client.NewClients(config)
	repositories := repository.NewRepositories(db)
	services := service.NewServices(repositories, config, clients, m)
	controllers := controller.NewControllers(services)
	controllers.Route(e)
	go controllers.Loop()

	//go controllers.Measurement().InsertSamples()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	logrus.Info("Starting HTTP server on port " + port)
	logrus.Fatal(e.Start(":" + port))
}
