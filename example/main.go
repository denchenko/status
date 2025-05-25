package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net/http"

	"github.com/denchenko/status"
)

func main() {
	// Create a health checker
	healthChecker := status.NewHealthChecker().
		WithTarget("database", status.TargetImportanceHigh, func(_ context.Context) error {
			// Implement your database health check here
			return generateRandomError()
		}).
		WithTarget("network", status.TargetImportanceLow, func(_ context.Context) error {
			// Implement your network health check here
			return generateRandomError()
		})

	// Create a status page
	statusPage := status.NewPage(
		// Add health checker to status page
		status.WithHealthChecker(healthChecker),
		// Add additional links
		status.WithLink("OpenAPI Documentation", "/swagger"),
		status.WithLink("Metrics", "/metrics"),
	)

	// Set up HTTP handlers
	http.HandleFunc("/health", healthChecker.Handler())
	http.HandleFunc("/status", statusPage.Handler())

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func generateRandomError() error {
	if rand.Intn(2) == 0 {
		return nil
	}

	return errors.New("dependency is not healthy")
}
