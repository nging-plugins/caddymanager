package config

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/nging-plugins/caddymanager/application/library/engine"
	"github.com/nging-plugins/caddymanager/application/library/thirdparty"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

const Name = `nginx`

var (
	regexConfigFile    = regexp.MustCompile(`[\s]+configuration file (.+\.conf)[\s]+`)
	regexConfigInclude = regexp.MustCompile(`^include[\s]+([\S]+(?:\*|\*\.conf))[\s]*;(?:[\s]*#.*)?[\s]*$`)
	regexVersion       = regexp.MustCompile(`[\d]+\.[\d]+\.[\d]+`)
)

type Config struct {
	Command          string
	CmdWithConfig    bool
	Version          string
	ConfigPath       string
	ConfigInclude    string
	ID               string
	WorkDir          string
	Environ          string
	CertLocalDir     string
	CertContainerDir string
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

func (c *Config) GetTemplateFile() string {
	return Name
}

func (c *Config) GetIdent() string {
	return c.ID
}

func (c *Config) GetEngine() string {
	return Name
}

func (c *Config) GetEnviron() string {
	return c.Environ
}

func (c *Config) GetCertLocalDir() string {
	return c.CertLocalDir
}

func (c *Config) GetCertContainerDir() string {
	return c.CertContainerDir
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

func (c *Config) getVersion(ctx context.Context) (string, error) {
	result, err := c.exec(ctx, `-v`)
	if err != nil {
		return ``, err
	}
	result = bytes.TrimSpace(result)
	parts := bytes.SplitN(result, []byte(`:`), 2)
	if len(parts) != 2 {
		matches := regexVersion.FindStringSubmatch(string(result))
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
	var includeConfD string
	var includeSitesEnabled string
	var includeDir string
	err := com.SeekFileLines(confPath, func(line string) error {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			return nil
		}
		if strings.HasPrefix(line, `#`) {
			return nil
		}
		matches := regexConfigInclude.FindAllStringSubmatch(line, 1)
		if len(matches) == 0 || len(matches[0]) < 2 {
			return nil
		}
		line = matches[0][1]
		if strings.Contains(line, `sites-`) /*strings.Contains(line, `sites-enabled`) || strings.Contains(line, `sites-available`)*/ {
			includeSitesEnabled = line
			return echo.ErrExit
		}
		if strings.Contains(line, `conf.d`) {
			includeConfD = line
			return nil
		}
		includeDir = line
		return nil
	})
	if err != nil && err != echo.ErrExit {
		return ``, err
	}
	if len(includeSitesEnabled) > 0 {
		includeDir = includeSitesEnabled
	} else if len(includeConfD) > 0 {
		includeDir = includeConfD
	}
	includeDir = com.TrimFileName(includeDir)
	return includeDir, nil
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
	if c.Environ == engine.EnvironContainer {
		rootArgs := com.ParseArgs(command)
		if len(rootArgs) > 1 {
			command = rootArgs[0]
			args = append(rootArgs[1:], args...)
		}
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

func (c *Config) RenewalCert(ctx context.Context, domain, email string) error {
	command := strings.TrimSpace(c.Command)
	command = strings.TrimSuffix(command, `.exe`)
	command = strings.TrimSuffix(command, `nginx`)
	command += `certbot`
	if c.Environ == engine.EnvironContainer {
		return RenewalCert(ctx, command, domain, email, c.CertContainerDir)
	}
	return RenewalCert(ctx, command, domain, email, c.CertLocalDir)
}

func RenewalCert(ctx context.Context, customCmd, domain, email, certDir string) error {
	//certbot certonly --webroot -d example.com --email info@example.com -w /var/www/_letsencrypt -n --agree-tos --force-renewal
	command := `certbot`
	args := []string{
		`certonly`,
		`--webroot`,
		`-d`, domain,
		`--email`, email,
		`-w`, certDir,
		`-n`,
		`--agree-tos`,
		`--force-renewal`,
	}
	if len(customCmd) > 0 {
		rootArgs := com.ParseArgs(customCmd)
		if len(rootArgs) > 1 {
			command = rootArgs[0]
			args = append(rootArgs[1:], args...)
		}
	}
	cmd := exec.CommandContext(ctx, command, args...)
	result, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf(`%s: %w`, result, err)
	}
	return err
}
