package main

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/eventuallyconsistentwrites/high-tide-server/countmin"
	"github.com/eventuallyconsistentwrites/high-tide-server/internal/domain"
	"github.com/eventuallyconsistentwrites/high-tide-server/internal/middleware"
	"github.com/eventuallyconsistentwrites/high-tide-server/internal/repository"
	"github.com/eventuallyconsistentwrites/high-tide-server/internal/routes"
	"github.com/eventuallyconsistentwrites/high-tide-server/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

func main() {
	// Create a structured JSON logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Load environment variables from .env file for local development.
	// In a container, this file is not expected to exist, so we ignore the "not found" error.
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			logger.Error("Failed to load .env file", "error", err)
			os.Exit(1)
		}
	}

	dbPath := os.Getenv("DB_PATH")
	resetIntervalStr := os.Getenv("RESET_INTERVAL")
	listenAddr := os.Getenv("LISTEN_ADDR")
	logger.Info(
		"From .env",
		"dbPath", dbPath,
		"resetIntervalStr", resetIntervalStr,
		"listenAddr", listenAddr,
	)
	if dbPath == "" {
		dbPath = "high_tide.db" // Default for local, non-containerized runs
	}

	resetInterval := 30 // Default value in seconds
	if resetIntervalStr != "" {
		var err error
		resetInterval, err = strconv.Atoi(resetIntervalStr)
		if err != nil {
			logger.Warn("Invalid RESET_INTERVAL, using default value", "value", resetIntervalStr, "defaultSeconds", resetInterval)
			resetInterval = 30 // Fallback to default if parsing fails
		}
	} else {

	}
	logger.Info("Count-Min Sketch reset interval configured", "seconds", resetInterval)

	// Connect to sqlite database
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		logger.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	// Sync DB schema with GoLang Structs
	db.AutoMigrate(&domain.Post{})

	// Inject database into repository
	postRepo := repository.NewPostSQLRepository(db)

	// Inject repository into service
	postSvc := service.NewPostService(postRepo)

	// Inject service into routes
	postHandler := routes.NewPostRoutes(postSvc)

	// Setup router
	mux := http.NewServeMux()

	// Add a healthcheck endpoint for monitoring and robust startup.
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Register routes
	postHandler.RegisterRoutes(mux)

	// Initialize Count-Min Sketch and Rate Limiter middleware.
	// These values can be tuned for your specific needs.
	cms := countmin.NewCountMinSketch(5, 25)

	// Periodically reset the Count-Min Sketch to avoid saturation with old values.
	// This goroutine will create a new ticker that fires at the specified interval.
	// On each tick, it will reset the sketch.
	go func() {
		ticker := time.NewTicker(time.Duration(resetInterval) * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			logger.Info("Resetting count-min sketch")
			cms.Reset()
		}
	}()
	rateLimiter := middleware.NewRateLimiter(cms, 20, logger) // Block IPs after 20 requests

	// Wrap the main router with the rate-limiting middleware.
	// All requests will now pass through the rate limiter first.
	var handler http.Handler = mux
	handler = rateLimiter.Limit(handler)

	if listenAddr == "" {
		listenAddr = ":8080" // Default for local, non-containerized runs
	}

	// START SERVER
	logger.Info("Server starting", "address", listenAddr)
	if err := http.ListenAndServe(listenAddr, handler); err != nil {
		logger.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
