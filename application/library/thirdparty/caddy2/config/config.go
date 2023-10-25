package config

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/webx-top/com"
)

type Config struct {
	Command   string
	Endpoint  string
	Caddyfile string
}

func (c *Config) Init() error {
	var err error
	return err
}

func (c *Config) Start(ctx context.Context) error {
	_, err := c.exec(ctx, `run`, `--config`, c.Caddyfile)
	return err
}

func (c *Config) Reload(ctx context.Context) error {
	_, err := c.exec(ctx, `reload`, `--config`, c.Caddyfile)
	return err
}

func (c *Config) Stop(ctx context.Context) error {
	_, err := c.exec(ctx, `stop`)
	return err
}

func (c *Config) exec(ctx context.Context, args ...string) ([]byte, error) {
	if len(c.Command) == 0 {
		c.Command = `caddy`
		if com.IsWindows {
			c.Command += `.exe`
		}
	}
	cmd := exec.CommandContext(ctx, c.Command, args...)
	result, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf(`%w: %s`, err, cmd.String())
	}
	return result, err
}
