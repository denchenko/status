# status

[![Run Tests](https://github.com/denchenko/status/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/denchenko/status/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/denchenko/status/branch/master/graph/badge.svg)](https://codecov.io/gh/denchenko/status)
[![Go Report Card](https://goreportcard.com/badge/github.com/denchenko/status)](https://goreportcard.com/report/github.com/denchenko/status)
[![GoDoc](https://godoc.org/github.com/denchenko/status?status.svg)](https://godoc.org/github.com/denchenko/status)

# Example

```
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
```
