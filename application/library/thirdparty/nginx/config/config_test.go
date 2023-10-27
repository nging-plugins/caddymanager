package config

import (
	"context"
	"testing"

	"github.com/nging-plugins/caddymanager/application/library/engine"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
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
	c := &Config{
		//Environ: engine.EnvironContainer,
		//Command: `docker exec nginx nginx`,
	}
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
	assert.Equal(t, `/etc/nginx/sites-enabled/`, c.ConfigInclude)

	err = c.TestConfig(ctx)
	assert.NoError(t, err)

	c.CmdWithConfig = true
	c.ConfigPath = `/not-exists.conf`
	err = c.TestConfig(ctx)
	assert.Error(t, err)
}

func _TestConfig2(t *testing.T) {
	c := &Config{
		Environ: engine.EnvironContainer,
		Command: `docker exec nginx nginx`,
	}
	ctx := context.Background()
	var err error
	assert.NoError(t, err)
	var b []byte
	b, err = c.exec(ctx, `-v`)
	assert.NoError(t, err)
	t.Logf(`%s`, b)
	matches := regexVersion.FindStringSubmatch(string(b))
	com.Dump(matches)
	_, _ = c, ctx
}
