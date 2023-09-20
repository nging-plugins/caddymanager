package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfigFilePath(t *testing.T) {
	line := []byte(` configuration file /etc/nginx/nginx.conf `)
	var configFilePath string
	matches := regexConfigFile.FindAllSubmatch(line, 1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		configFilePath = string(matches[0][1])
	}
	assert.Equal(t, `/etc/nginx/nginx.conf`, configFilePath)
}

func TestConfig(t *testing.T) {
	c := &Config{}
	ctx := context.Background()
	var err error
	c.Version, err = c.getVersion(ctx)
	assert.NoError(t, err)
	assert.Equal(t, `1.18.0`, c.Version)
	c.ConfigPath, err = c.getConfigFilePath(ctx)
	assert.NoError(t, err)
	assert.Equal(t, `/etc/nginx/nginx.conf`, c.ConfigPath)

	c.ConfigInclude, err = c.getConfigIncludePath(c.ConfigPath)
	assert.NoError(t, err)
	assert.Equal(t, `/etc/nginx/modules-enabled/`, c.ConfigInclude)
}
