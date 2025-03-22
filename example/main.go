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
	healthChecker := status.NewHealthChecker().
		WithTarget("db", status.TargetImportanceHigh, func(ctx context.Context) error {
			return generateRandomError()
		}).
		WithTarget("network", status.TargetImportanceLow, func(ctx context.Context) error {
			return generateRandomError()
		})

	statusPage, err := status.NewPage()
	if err != nil {
		log.Printf("constructing page: %v", err)
		return
	}
	statusPage.WithHealthChecker(healthChecker)
	statusPage.WithURL("Swagger", "/swagger")

	http.HandleFunc("/health", healthChecker.Handler())
	http.HandleFunc("/status", statusPage.Handler())

	http.ListenAndServe(":8080", nil)
}

func generateRandomError() error {
	if rand.Intn(2) == 0 {
		return nil
	}

	return errors.New("dependency is not healthy")
}
