package config

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/admpub/log"
	"github.com/nging-plugins/caddymanager/application/library/thirdparty"
	"github.com/webx-top/com"
)

type Config struct {
	Command   string
	Endpoint  string
	Caddyfile string
}

func (c *Config) Init() *Config {
	return c
}

func (c *Config) Start(ctx context.Context) error {
	err := c.exec(ctx, `start`, `--config`, c.Caddyfile)
	return err
}

func (c *Config) Reload(ctx context.Context) error {
	err := c.exec(ctx, `reload`, `--config`, c.Caddyfile)
	return err
}

func (c *Config) Stop(ctx context.Context) error {
	err := c.exec(ctx, `stop`)
	return err
}

func (c *Config) exec(ctx context.Context, args ...string) error {
	if len(c.Command) == 0 {
		c.Command = `caddy`
		if com.IsWindows {
			c.Command += `.exe`
		}
	}
	cmd := exec.CommandContext(ctx, c.Command, args...)
	if stderr := thirdparty.GetCtxStderr(ctx); stderr != nil {
		cmd.Stderr = stderr
	}
	if stdout := thirdparty.GetCtxStdout(ctx); stdout != nil {
		cmd.Stdout = stdout
	}
	if cmd.Stderr == nil && cmd.Stdout == nil {
		result, err := cmd.CombinedOutput()
		if err != nil {
			err = fmt.Errorf(`%w: %s`, err, cmd.String())
		} else {
			log.Infof(`%s`, result)
		}
		return err
	}
	err := cmd.Run()
	return err
}
