package main

import (
	"log"
	"net/http"
	"os"

	"service-test-runner/internal/config"
	httpDelivery "service-test-runner/internal/delivery/http"
	"service-test-runner/internal/repository/selenium"
	"service-test-runner/internal/usecase/automation"
	"service-test-runner/internal/usecase/testsuite"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	seleniumRepo := selenium.NewSeleniumRepository(cfg.Projects)
	automationUsecase := automation.NewAutomationUsecase(seleniumRepo)
	testsuiteUsecase := testsuite.NewTestSuiteUsecase(seleniumRepo)

	router := mux.NewRouter()
	handler := httpDelivery.NewHandler(automationUsecase, testsuiteUsecase)
	httpDelivery.RegisterRoutes(router, handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "6000"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
