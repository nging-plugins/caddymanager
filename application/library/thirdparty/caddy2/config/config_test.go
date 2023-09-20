package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c := &Config{}
	ctx := context.Background()
	var err error
	assert.NoError(t, err)

	_, _ = c, ctx
}
