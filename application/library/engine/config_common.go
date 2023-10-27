package engine

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/nging-plugins/caddymanager/application/dbschema"
	"github.com/webx-top/com"
)

func NewConfig(engineName, templateFile string) *CommonConfig {
	return &CommonConfig{engineName: engineName, templateFile: templateFile}
}

type CommonConfig struct {
	ID                        string
	Command                   string
	CmdWithConfig             bool
	WorkDir                   string
	Environ                   string
	EngineConfigLocalFile     string
	EngineConfigContainerFile string
	VhostConfigLocalDir       string
	VhostConfigContainerDir   string
	CertLocalDir              string
	CertContainerDir          string

	engineName   string
	templateFile string
}

func (c *CommonConfig) GetIdent() string {
	return c.ID
}

func (c *CommonConfig) GetTemplateFile() string {
	return c.templateFile
}

func (c *CommonConfig) GetEngine() string {
	return c.engineName
}

func (c *CommonConfig) GetEnviron() string {
	return c.Environ
}

func (c *CommonConfig) GetCertLocalDir() string {
	return c.CertLocalDir
}

func (c *CommonConfig) GetCertContainerDir() string {
	return c.CertContainerDir
}

func (c *CommonConfig) GetVhostConfigLocalDir() string {
	return c.VhostConfigLocalDir
}

func (c *CommonConfig) GetVhostConfigContainerDir() string {
	return c.VhostConfigContainerDir
}

func (c *CommonConfig) GetEngineConfigLocalFile() string {
	return c.EngineConfigLocalFile
}

func (c *CommonConfig) GetEngineConfigContainerFile() string {
	return c.EngineConfigContainerFile
}

func (c *CommonConfig) EngineConfigFile() string {
	if c.Environ == EnvironContainer {
		return c.EngineConfigContainerFile
	}
	return c.EngineConfigLocalFile
}

func (c *CommonConfig) VhostConfigDir() string {
	if c.Environ == EnvironContainer {
		return c.VhostConfigContainerDir
	}
	return c.VhostConfigLocalDir
}

func (c *CommonConfig) CopyFrom(m *dbschema.NgingVhostServer) {
	c.ID = m.Ident
	c.Command = m.ExecutableFile
	c.CmdWithConfig = m.CmdWithConfig == common.BoolY
	c.WorkDir = m.WorkDir
	c.Environ = m.Environ
	c.EngineConfigLocalFile, _ = filepath.Abs(m.ConfigLocalFile)
	c.EngineConfigContainerFile = m.ConfigContainerFile
	c.VhostConfigLocalDir, _ = filepath.Abs(m.VhostConfigLocalDir)
	c.VhostConfigContainerDir = m.VhostConfigContainerDir
	c.CertLocalDir, _ = filepath.Abs(m.CertLocalDir)
	c.CertContainerDir = m.CertContainerDir
}

func (c *CommonConfig) Exec(ctx context.Context, args ...string) ([]byte, error) {
	command := c.Command
	if c.Environ == EnvironContainer {
		rootArgs := com.ParseArgs(command)
		if len(rootArgs) > 1 {
			command = rootArgs[0]
			args = append(rootArgs[1:], args...)
		}
	}
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = c.WorkDir
	if stderr := GetCtxStderr(ctx); stderr != nil {
		cmd.Stderr = stderr
	}
	if stdout := GetCtxStdout(ctx); stdout != nil {
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
