package api

import (
	"log"
	"log/slog"
	"os"

	"github.com/abdulazizax/udevslab-lesson3/internal/config"
	"github.com/abdulazizax/udevslab-lesson3/internal/http/app"
	"github.com/abdulazizax/udevslab-lesson3/internal/http/handler"
	"github.com/abdulazizax/udevslab-lesson3/internal/service"
	"github.com/abdulazizax/udevslab-lesson3/internal/storage"
	mongo "github.com/abdulazizax/udevslab-lesson3/internal/storage/mongodb"
)

func Run() error {
	// Load configuration and handle errors
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
		return err // Add return statement after log.Fatal for better error handling
	}

	// Set up logger with file output
	logFile, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer logFile.Close()

	// Initialize structured logger with JSON format
	logger := slog.New(slog.NewJSONHandler(logFile, nil))

	// Initialize MongoDB connection
	db, err := mongo.ConnectDB(cfg)
	if err != nil {
		logger.Error("Error while connecting to MongoDB", slog.String("err", err.Error()))
		return err
	}

	// Initialize storage layer with MongoDB and Redis
	storage := storage.New(db, cfg, logger)

	// Initialize service layer
	service := service.NewService(logger, storage)

	// Initialize HTTP handler
	handler := handler.NewHandler(logger, service, cfg)

	// Start the HTTP server
	return app.Run(handler, logger, cfg)
}
