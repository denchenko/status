# status

[![Run Tests](https://github.com/denchenko/status/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/denchenko/status/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/denchenko/status)](https://goreportcard.com/report/github.com/denchenko/status)
[![GoDoc](https://godoc.org/github.com/denchenko/status?status.svg)](https://godoc.org/github.com/denchenko/status)

A Go package for health checking and status page generation. It provides a simple way to monitor the health of various dependencies and services in your application, with a status page dashboard.

## Installation

```bash
go get github.com/denchenko/status
```

## Usage

```go
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
		WithTarget("database", status.TargetImportanceHigh, func(ctx context.Context) error {
			// Implement your database health check here
			return generateRandomError()
		}).
		WithTarget("network", status.TargetImportanceLow, func(ctx context.Context) error {
			// Implement your network health check here
			return generateRandomError()
		})

	// Create a status page
	statusPage, err := status.NewPage()
	if err != nil {
		log.Fatal(err)
	}

	// Add health checker to status page
	statusPage.WithHealthChecker(healthChecker)

	// Add additional URLs to the status page
	statusPage.WithURL("OpenAPI Documentation", "/swagger")
	statusPage.WithURL("Metrics", "/metrics")

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
```
