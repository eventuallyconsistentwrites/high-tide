package main

import (
	"log"
	"net/http"

	"github.com/eventuallyconsistentwrites/high-tide-server/countmin"
	"github.com/eventuallyconsistentwrites/high-tide-server/internal/domain"
	"github.com/eventuallyconsistentwrites/high-tide-server/internal/middleware"
	"github.com/eventuallyconsistentwrites/high-tide-server/internal/repository"
	"github.com/eventuallyconsistentwrites/high-tide-server/internal/routes"
	"github.com/eventuallyconsistentwrites/high-tide-server/internal/service"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Connect to sqlite database
	db, err := gorm.Open(sqlite.Open("high_tide.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
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

	// Register routes
	postHandler.RegisterRoutes(mux)

	// Initialize Count-Min Sketch and Rate Limiter middleware.
	// These values can be tuned for your specific needs.
	cms := countmin.NewCountMinSketch(5, 25)
	rateLimiter := middleware.NewRateLimiter(cms, 20) // Block IPs after 20 requests

	// Wrap the main router with the rate-limiting middleware.
	// All requests will now pass through the rate limiter first.
	var handler http.Handler = mux
	handler = rateLimiter.Limit(handler)

	// START SERVER
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
