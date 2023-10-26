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

	err = c.exec(ctx, `version`)
	assert.NoError(t, err)

	err = c.exec(ctx, `help`)
	assert.NoError(t, err)

	err = c.TestConfig(ctx)
	assert.NoError(t, err)

	c.CmdWithConfig = true
	c.Caddyfile = `/not-exists.conf`
	err = c.TestConfig(ctx)
	assert.Error(t, err)
	_, _ = c, ctx
}
