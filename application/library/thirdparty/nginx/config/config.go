package config

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/webx-top/com"
)

var (
	regexConfigFile    = regexp.MustCompile(`[\s]+configuration file (.+\.conf)[\s]+`)
	regexConfigInclude = regexp.MustCompile(`[\s]+include[\s]+(.+\.conf)[\s]*;[\s]*` + "\n")
)

type Config struct {
	Command       string
	Version       string
	ConfigPath    string
	ConfigInclude string
}

func (c *Config) Init() error {
	var err error
	ctx := context.Background()
	if len(c.Version) == 0 {
		c.Version, err = c.getVersion(ctx)
		if err != nil {
			return err
		}
	}
	if len(c.ConfigPath) == 0 {
		c.ConfigPath, err = c.getConfigFilePath(ctx)
		if err != nil {
			return err
		}
	}
	if len(c.ConfigInclude) == 0 {
		c.ConfigInclude, err = c.getConfigIncludePath(c.ConfigPath)
	}
	return err
}

func (c *Config) Start() error {
	_, err := c.exec(
		context.Background(),
		//`-c`, `/config/nginx/nginx.conf`,
	)
	return err
}

func (c *Config) getVersion(ctx context.Context) (string, error) {
	result, err := c.exec(ctx, `-v`)
	if err != nil {
		return ``, err
	}
	result = bytes.TrimSpace(result)
	parts := bytes.SplitN(result, []byte(`:`), 2)
	if len(parts) != 2 {
		return ``, err
	}
	result = parts[1]
	parts = bytes.SplitN(result, []byte(`/`), 2)
	if len(parts) != 2 {
		return ``, err
	}
	parts[1] = bytes.TrimSpace(bytes.SplitN(parts[1], []byte(` `), 2)[0])
	return string(parts[1]), err
}

func (c *Config) getConfigFilePath(ctx context.Context) (string, error) {
	result, err := c.exec(ctx, `-t`)
	if err != nil {
		return ``, err
	}
	result = bytes.TrimSpace(result)
	lines := bytes.Split(result, []byte{com.LF})
	var configFilePath string
	for _, line := range lines {
		matches := regexConfigFile.FindAllSubmatch(line, 1)
		if len(matches) > 0 && len(matches[0]) > 1 {
			configFilePath = string(matches[0][1])
			break
		}
	}
	return configFilePath, err
}

func (c *Config) getConfigIncludePath(confPath string) (string, error) {
	if len(confPath) == 0 {
		return ``, nil
	}
	content, err := os.ReadFile(confPath)
	if err != nil {
		return ``, err
	}
	matches := regexConfigInclude.FindAllSubmatch(content, 1)
	var include string
	if len(matches) > 0 && len(matches[0]) > 1 {
		include = string(matches[0][1])
		include = strings.TrimSpace(include)
		include = strings.SplitN(include, `*`, 2)[0]
	}
	return include, nil
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
		c.Command = `nginx`
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
