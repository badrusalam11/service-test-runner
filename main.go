package main

import (
	"log"
	"net/http"
	"os"

	"service-test-runner/internal/config"
	"service-test-runner/internal/db"
	httpDelivery "service-test-runner/internal/delivery/http"
	handler "service-test-runner/internal/delivery/http/handler"
	automationRepo "service-test-runner/internal/repository/automation"
	"service-test-runner/internal/repository/project"
	"service-test-runner/internal/repository/selenium"
	"service-test-runner/internal/usecase"

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
	projectRepo := project.NewProjectRepository(projects)
	queueAutomationRepository := automationRepo.NewQueueAutomationRepository()

	// Initialize use cases.
	automationUsecase := usecase.NewAutomationUsecase(seleniumRepo)
	queueAutomationUsecase := usecase.NewQueueAutomationUseCase(queueAutomationRepository)
	testsuiteUsecase := usecase.NewTestSuiteUsecase(seleniumRepo)
	projectUsecase := usecase.NewProjectUsecase(projectRepo)

	// Setup HTTP router.
	router := mux.NewRouter()
	handler := handler.NewHandler(
		automationUsecase,
		queueAutomationUsecase,
		testsuiteUsecase,
		projectUsecase)
	httpDelivery.RegisterRoutes(router, handler)

	// Start the server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "6000"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
