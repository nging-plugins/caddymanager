package cmder

import (
	"context"
	"io"

	"github.com/admpub/once"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/config/cmder"

	caddy2Config "github.com/nging-plugins/caddymanager/application/library/thirdparty/caddy2/config"
)

const Name = `caddy2`

func Initer() interface{} {
	return &caddy2Config.Config{}
}

func Get() cmder.Cmder {
	return cmder.Get(Name)
}

func GetNginxConfig() *caddy2Config.Config {
	return GetNginxCmd().NginxConfig()
}

func GetNginxCmd() *caddy2Cmd {
	cm := cmder.Get(Name).(*caddy2Cmd)
	return cm
}

func New() cmder.Cmder {
	return &caddy2Cmd{
		CLIConfig: config.FromCLI(),
		once:      once.Once{},
	}
}

type caddy2Cmd struct {
	CLIConfig    *config.CLIConfig
	caddy2Config *caddy2Config.Config
	once         once.Once
}

func (c *caddy2Cmd) Boot() error {
	return c.NginxConfig().Init()
}

func (c *caddy2Cmd) getConfig() *config.Config {
	if config.FromFile() == nil {
		c.CLIConfig.ParseConfig()
	}
	return config.FromFile()
}

func (c *caddy2Cmd) parseConfig() {
	c.caddy2Config, _ = c.getConfig().Extend.Get(Name).(*caddy2Config.Config)
	if c.caddy2Config == nil {
		c.caddy2Config = &caddy2Config.Config{}
	}
}

func (c *caddy2Cmd) NginxConfig() *caddy2Config.Config {
	c.once.Do(c.parseConfig)
	return c.caddy2Config
}

func (c *caddy2Cmd) StopHistory(_ ...string) error {
	return nil
}

func (c *caddy2Cmd) Start(writer ...io.Writer) error {
	return nil
}

func (c *caddy2Cmd) Stop() error {
	return nil
}

func (c *caddy2Cmd) Reload() error {
	return c.NginxConfig().Reload(context.Background())
}

func (c *caddy2Cmd) Restart(writer ...io.Writer) error {
	return nil
}
