package status

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthChecker_Check(t *testing.T) {
	ctx := context.Background()
	hc := NewHealthChecker()
	results, err := hc.Check(ctx)
	assert.Empty(t, results)
	assert.NoError(t, err)
}
