package config

import (
	"context"
	"os"
	"path/filepath"
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
	c := New()
	ctx := context.Background()
	var err error
	c.Version, err = c.getVersion(ctx)
	assert.NoError(t, err)
	assert.Equal(t, `1.18.0`, c.Version)
	c.EngineConfigLocalFile, err = c.getEngineConfigLocalFile(ctx)
	assert.NoError(t, err)
	assert.Equal(t, `/etc/nginx/nginx.conf`, c.EngineConfigLocalFile)

	c.VhostConfigLocalDir, err = c.getVhostConfigLocalDir(c.EngineConfigLocalFile)
	assert.NoError(t, err)
	assert.Equal(t, `/etc/nginx/sites-enabled/`, c.VhostConfigLocalDir)

	err = c.TestConfig(ctx)
	assert.NoError(t, err)

	c.CmdWithConfig = true
	c.EngineConfigLocalFile = `/not-exists.conf`
	err = c.TestConfig(ctx)
	assert.Error(t, err)
}

func _TestConfig2(t *testing.T) {
	c := New()
	c.Environ = engine.EnvironContainer
	c.Command = `docker exec nginx nginx`
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

func TestFixEngineConfigFile(t *testing.T) {
	c := New()
	var err error
	os.Remove(`./testdata/nginx.conf`)
	originalFile, err := filepath.Abs(`./testdata/nginx.conf.original`)
	com.Copy(originalFile, `./testdata/nginx.conf`)
	assert.NoError(t, err)
	c.Environ = engine.EnvironLocal
	c.VhostConfigLocalDir = `/etc/nginx/sites-nging/`
	c.EngineConfigLocalFile, _ = filepath.Abs(`./testdata/nginx.conf`)
	var hasUpdate bool
	hasUpdate, err = c.FixEngineConfigFile()
	assert.NoError(t, err)
	assert.True(t, hasUpdate)

	hasUpdate, err = c.FixEngineConfigFile()
	assert.NoError(t, err)
	assert.False(t, hasUpdate)

	hasUpdate, err = c.FixEngineConfigFile(true)
	assert.NoError(t, err)
	assert.True(t, hasUpdate)
}
