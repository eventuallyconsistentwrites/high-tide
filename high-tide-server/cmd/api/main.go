package main

import (
	"log"
	"net/http"

	"github.com/eventuallyconsistentwrites/high-tide-server/internal/domain"
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

	// 6. START SERVER
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
