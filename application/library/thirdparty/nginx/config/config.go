package config

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/nging-plugins/caddymanager/application/library/thirdparty"
	"github.com/webx-top/com"
)

const Name = `nginx`

var (
	regexConfigFile    = regexp.MustCompile(`[\s]+configuration file (.+\.conf)[\s]+`)
	regexConfigInclude = regexp.MustCompile(`[\s]+include[\s]+(.+\.conf)[\s]*;[\s]*` + "\n")
)

type Config struct {
	Command       string
	CmdWithConfig bool
	Version       string
	ConfigPath    string
	ConfigInclude string
	ID            string
	WorkDir       string
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
	if len(c.ConfigInclude) == 0 {
		if len(c.ConfigPath) == 0 {
			c.ConfigPath, err = c.getConfigFilePath(ctx)
			if err != nil {
				return err
			}
		}
		c.ConfigInclude, err = c.getConfigIncludePath(c.ConfigPath)
	}
	return err
}

func DefaultConfigDir() string {
	return filepath.Join(config.FromCLI().ConfDir(), `vhosts-nginx`)
}

func (c *Config) GetVhostConfigDirAbsPath() (string, error) {
	if len(c.ConfigInclude) == 0 {
		if err := c.Init(); err != nil {
			log.Error(err)
			if len(c.ConfigInclude) == 0 {
				c.ConfigInclude = filepath.Join(DefaultConfigDir(), c.ID)
			}
		}
	}
	return c.ConfigInclude, nil
}

func (c *Config) TemplateFile() string {
	return Name
}

func (c *Config) Ident() string {
	return c.ID
}

func (c *Config) Engine() string {
	return Name
}

func (c *Config) Start(ctx context.Context) error {
	args := []string{}
	if c.CmdWithConfig && len(c.ConfigPath) > 0 && strings.HasSuffix(c.ConfigPath, `.conf`) {
		args = append(args, `-c`, c.ConfigPath)
	}
	_, err := c.exec(ctx)
	return err
}

func (c *Config) TestConfig(ctx context.Context) error {
	args := []string{`-t`}
	if c.CmdWithConfig && len(c.ConfigPath) > 0 {
		args = append(args, `-c`, c.ConfigPath)
	}
	//nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
	//nginx: configuration file /etc/nginx/nginx.conf test is successful
	_, err := c.exec(ctx, args...)
	return err
}

var versionRegex = regexp.MustCompile(`[\d]+\.[\d]+\.[\d]+`)

func (c *Config) getVersion(ctx context.Context) (string, error) {
	result, err := c.exec(ctx, `-v`)
	if err != nil {
		return ``, err
	}
	result = bytes.TrimSpace(result)
	parts := bytes.SplitN(result, []byte(`:`), 2)
	if len(parts) != 2 {
		matches := versionRegex.FindStringSubmatch(string(result))
		if len(matches) > 0 {
			return matches[0], err
		}
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
	command := c.Command
	rootArgs := com.ParseArgs(command)
	if len(rootArgs) > 1 {
		command = rootArgs[0]
		args = append(rootArgs[1:], args...)
	}
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = c.WorkDir
	if stderr := thirdparty.GetCtxStderr(ctx); stderr != nil {
		cmd.Stderr = stderr
	}
	if stdout := thirdparty.GetCtxStdout(ctx); stdout != nil {
		cmd.Stdout = stdout
	}
	if cmd.Stderr == nil && cmd.Stdout == nil {
		result, err := cmd.CombinedOutput()
		if err != nil {
			err = fmt.Errorf(`%s: %w`, result, err)
		}
		return result, err
	}
	err := cmd.Run()
	return nil, err
}
