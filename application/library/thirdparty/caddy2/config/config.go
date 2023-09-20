package config

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/webx-top/com"
)

type Config struct {
	Command  string
	Endpoint string
}

func (c *Config) Init() error {
	var err error
	return err
}

func (c *Config) Start() error {
	return nil
}

func (c *Config) Reload(ctx context.Context) error {
	return c.sendSignal(ctx, `reload`)
}

func (c *Config) Stop(ctx context.Context) error {
	return c.sendSignal(ctx, `stop`)
}

func (c *Config) Quit(ctx context.Context) error {
	return c.sendSignal(ctx, `quit`)
}

func (c *Config) Reopen(ctx context.Context) error {
	return c.sendSignal(ctx, `reopen`)
}

// signal: stop, quit, reopen, reload
func (c *Config) sendSignal(ctx context.Context, signal string) error {
	_, err := c.exec(ctx, `-s`, signal)
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
