package status

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type HealthTarget struct {
	Name       string           `json:"name"`
	Importance TargetImportance `json:"importance"`
	check      HealthCheckFunc
}

type TargetImportance string

const (
	TargetImportanceLow  = TargetImportance("low")
	TargetImportanceHigh = TargetImportance("high")
)

type HealthCheckFunc func(ctx context.Context) error

type HealthChecker struct {
	targets []HealthTarget
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{}
}

func (c *HealthChecker) WithTarget(name string, importance TargetImportance, check HealthCheckFunc) *HealthChecker {
	c.targets = append(c.targets, HealthTarget{
		Name:       name,
		Importance: importance,
		check:      check,
	})

	return c
}

func (c *HealthChecker) Handler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, noDeps := r.URL.Query()["no_deps"]; noDeps {
			w.WriteHeader(http.StatusOK)
			return
		}

		ctx := r.Context()

		results, err := c.Check(ctx)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err)
			return
		}

		status := http.StatusOK

		for _, result := range results {
			if result.Target.Importance == TargetImportanceHigh &&
				(result.Status != HealthTargetStatusOk || result.err != nil) {
				status = http.StatusInternalServerError
				break
			}
		}

		respondJSON(w, status, results)
	})
}

type HealthTargetStatus string

const (
	HealthTargetStatusOk   = HealthTargetStatus("ok")
	HealthTargetStatusFail = HealthTargetStatus("fail")
)

type HealthCheckResult struct {
	Target       HealthTarget       `json:"target"`
	Status       HealthTargetStatus `json:"status"`
	ErrorMessage string             `json:"error,omitempty"`
	err          error
}

func (c *HealthChecker) Check(ctx context.Context) ([]HealthCheckResult, error) {
	results := make([]HealthCheckResult, len(c.targets))

	g, ctx := errgroup.WithContext(ctx)

	for i, target := range c.targets {
		g.Go(func() error {
			err := target.check(ctx)
			if err != nil {
				results[i] = HealthCheckResult{
					Target:       target,
					Status:       HealthTargetStatusFail,
					err:          err,
					ErrorMessage: err.Error(),
				}
			} else {
				results[i] = HealthCheckResult{
					Target: target,
					Status: HealthTargetStatusOk,
				}
			}

			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		return nil, fmt.Errorf("waiting errgroup: %w", err)
	}

	return results, nil
}

// respondJSON responds JSON body with a given code. It sets
// Content-Type header.
func respondJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(&data); err != nil {
		log.Printf("encoding data to respond with json: %v", err)
	}
}
