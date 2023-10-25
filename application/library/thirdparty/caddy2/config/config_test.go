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

	r, err := c.exec(ctx, `version`)
	assert.NoError(t, err)
	t.Log(string(r))

	r, err = c.exec(ctx, `help`)
	assert.NoError(t, err)
	t.Log(string(r))
	_, _ = c, ctx
}
