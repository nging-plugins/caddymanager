package embedwebserver

import (
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/admpub/log"
	"github.com/webx-top/com"

	"github.com/admpub/nging/v4/application/library/config"

	"github.com/nging-plugins/caddymanager/pkg/library/caddy"
	scmder "github.com/nging-plugins/caddymanager/pkg/library/cmder"
)

const Name = "embedwebserver"

func init() {
	scmder.Register(Name, `Caddy(v1)`, New)
}

func Initer() interface{} {
	return &caddy.Config{}
}

func New() scmder.ServerCmder {
	return &caddyCmd{
		CLIConfig: config.DefaultCLIConfig,
		once:      sync.Once{},
	}
}

type caddyCmd struct {
	CLIConfig   *config.CLIConfig
	caddyConfig *caddy.Config
	caddyCh     *com.CmdChanReader
	once        sync.Once
}

func (c *caddyCmd) Init() error {
	caddy.TrapSignals()
	return c.CaddyConfig().Init().Start()
}

func (c *caddyCmd) getConfig() *config.Config {
	if config.DefaultConfig == nil {
		c.CLIConfig.ParseConfig()
	}
	return config.DefaultConfig
}

func (c *caddyCmd) parseConfig() {
	c.caddyConfig, _ = c.getConfig().Extend.Get(Name).(*caddy.Config)
	if c.caddyConfig == nil {
		c.caddyConfig = &caddy.Config{}
	}
	if len(c.caddyConfig.Caddyfile) == 0 {
		c.caddyConfig.Caddyfile = `./Caddyfile`
	} else if strings.HasSuffix(c.caddyConfig.Caddyfile, `/`) || strings.HasSuffix(c.caddyConfig.Caddyfile, `\`) {
		c.caddyConfig.Caddyfile = path.Join(c.caddyConfig.Caddyfile, `Caddyfile`)
	}
	caddy.SetDefaults(c.caddyConfig)
}

func (c *caddyCmd) CaddyConfig() *caddy.Config {
	c.once.Do(c.parseConfig)
	return c.caddyConfig
}

func (c *caddyCmd) StopHistory(_ ...string) error {
	if c.getConfig() == nil {
		return nil
	}
	return com.CloseProcessFromPidFile(c.CaddyConfig().PidFile)
}

func (c *caddyCmd) Start(writer ...io.Writer) error {
	err := c.StopHistory()
	if err != nil {
		log.Error(err.Error())
	}
	params := []string{os.Args[0], `--config`, c.CLIConfig.Conf, `--type`, Name}
	var cmd *exec.Cmd
	if caddy.EnableReload {
		c.caddyCh = com.NewCmdChanReader()
		cmd = com.RunCmdWithReaderWriter(params, c.caddyCh, writer...)
	} else {
		cmd = com.RunCmdWithWriter(params, writer...)
	}
	c.CLIConfig.CmdSet(Name, cmd)
	return nil
}

func (c *caddyCmd) Stop() error {
	defer func() {
		if c.caddyCh != nil {
			c.caddyCh.Close()
			c.caddyCh = nil
		}
	}()
	return c.CLIConfig.CmdStop(Name)
}

func (c *caddyCmd) Reload() error {
	if c.caddyCh == nil || com.IsWindows { //windows不支持重载，需要重启
		return c.Restart()
	}
	c.caddyCh.Send(com.BreakLine)
	return nil
	//c.CmdSendSignal("caddy", os.Interrupt)
}

func (c *caddyCmd) Restart(writer ...io.Writer) error {
	err := c.Stop()
	if err != nil {
		log.Error(err)
	}
	return c.Start(writer...)
}

func (c *caddyCmd) LogFile() string {
	return c.CaddyConfig().LogFile
}

func (c *caddyCmd) RemoveCachedCert(domain string) error {
	return c.CaddyConfig().RemoveCachedCert(domain)
}
