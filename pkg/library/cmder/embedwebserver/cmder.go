package embedwebserver

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/library/config"

	"github.com/nging-plugins/caddymanager/pkg/dbschema"
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

func (c *caddyCmd) VhostConfigFile(id uint) (string, error) {
	saveFile, err := c.ConfigDirectory()
	if err != nil {
		return saveFile, err
	}
	saveFile = filepath.Join(saveFile, fmt.Sprint(id)+`.conf`)
	return saveFile, nil
}

func (c *caddyCmd) SaveVhostConfig(ctx echo.Context, m *dbschema.NgingVhost, values url.Values) (string, error) {
	saveFile, err := c.VhostConfigFile(m.Id)
	if err != nil {
		return saveFile, err
	}
	ctx.Set(`values`, NewFormValues(values))
	b, err := ctx.Fetch(`caddy/caddyfile`, nil)
	if err != nil {
		return saveFile, err
	}
	b = com.CleanSpaceLine(b)
	log.Info(`Generate a Caddy configuration file: `, saveFile)
	err = os.WriteFile(saveFile, b, os.ModePerm)
	//jsonb, _ := caddyfile.ToJSON(b)
	//err = os.WriteFile(saveFile+`.json`, jsonb, os.ModePerm)
	return saveFile, err
}

func (c *caddyCmd) ConfigDirectory() (configDir string, err error) {
	configDir, err = filepath.Abs(config.DefaultConfig.Sys.VhostsfileDir)
	if err != nil {
		return
	}
	err = com.MkdirAll(configDir, os.ModePerm)
	return
}

func (c *caddyCmd) ClearVhostConfig() error {
	saveFile, err := c.ConfigDirectory()
	if err == nil {
		err = filepath.Walk(saveFile, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, `.conf`) {
				return nil
			}
			log.Info(`Delete the Caddy configuration file: `, path)
			return os.Remove(path)
		})
	}
	return err
}

func (c *caddyCmd) RemoveVhostConfig(m *dbschema.NgingVhost) (err error) {
	var saveFile string
	saveFile, err = c.VhostConfigFile(m.Id)
	if err != nil {
		return
	}
	err = os.Remove(saveFile)
	if os.IsNotExist(err) {
		err = nil
	}
	return
}
