package config

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/nging-plugins/caddymanager/application/library/thirdparty"
	"github.com/webx-top/com"
)

type Config struct {
	Command               string
	CmdWithConfig         bool
	Endpoint              string
	Caddyfile             string
	ConfigInclude         string
	ID                    string
	vhostConfigDirAbsPath string
}

func (c *Config) Init() *Config {
	return c
}

func DefaultConfigDir() string {
	return filepath.Join(config.FromCLI().ConfDir(), `vhosts-caddy2`)
}

func (c *Config) GetVhostConfigDirAbsPath() (string, error) {
	if len(c.ConfigInclude) == 0 {
		c.ConfigInclude = filepath.Join(DefaultConfigDir(), c.ID)
	}
	return c.ConfigInclude, nil
}

func (c *Config) TemplateFile() string {
	return `caddy2`
}

func (c *Config) Ident() string {
	return c.ID
}

func (c *Config) Engine() string {
	return `caddy2`
}

func (c *Config) Start(ctx context.Context) error {
	args := []string{`start`}
	if c.CmdWithConfig && len(c.Caddyfile) > 0 && com.IsFile(c.Caddyfile) {
		args = append(args, `--config`, c.Caddyfile)
	}
	err := c.exec(ctx, args...)
	return err
}

func (c *Config) Reload(ctx context.Context) error {
	args := []string{`reload`}
	if c.CmdWithConfig && len(c.Caddyfile) > 0 && com.IsFile(c.Caddyfile) {
		args = append(args, `--config`, c.Caddyfile)
	}
	err := c.exec(ctx, args...)
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
	command := c.Command
	rootArgs := com.ParseArgs(command)
	if len(rootArgs) > 1 {
		command = rootArgs[0]
		args = append(rootArgs[1:], args...)
	}
	cmd := exec.CommandContext(ctx, command, args...)
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
