package main

import (
	"log"
	"net/http"
	"os"

	"service-test-runner/internal/config"
	"service-test-runner/internal/db"
	httpDelivery "service-test-runner/internal/delivery/http"
	"service-test-runner/internal/repository/selenium"
	"service-test-runner/internal/usecase/automation"
	"service-test-runner/internal/usecase/testsuite"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration (from .env if it exists, otherwise from config.json)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize the database connection using GORM (see internal/db/db.go)
	if err := db.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Load project mappings from the database using GORM.
	projects, err := db.LoadProjects()
	if err != nil {
		log.Fatalf("Failed to load projects from database: %v", err)
	}

	// Initialize repository with the project mappings.
	seleniumRepo := selenium.NewSeleniumRepository(projects)

	// Initialize use cases.
	automationUsecase := automation.NewAutomationUsecase(seleniumRepo)
	testsuiteUsecase := testsuite.NewTestSuiteUsecase(seleniumRepo)

	// Setup HTTP router.
	router := mux.NewRouter()
	handler := httpDelivery.NewHandler(automationUsecase, testsuiteUsecase)
	httpDelivery.RegisterRoutes(router, handler)

	// Start the server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "6000"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
