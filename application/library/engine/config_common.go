package engine

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/webx-top/com"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/nging-plugins/caddymanager/application/dbschema"
)

func NewConfig(engineName, templateFile string) *CommonConfig {
	return &CommonConfig{engineName: engineName, templateFile: templateFile}
}

type ParsedCommand struct{
	Command string
	Args string
	ContainerEnging string
	ContainerName string
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
	parsedCommand *ParsedCommand
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

func ParseContainerInfo(parts []string) string, string {
	return "", ""
}

func (c *CommonConfig) Exec(ctx context.Context, args ...string) ([]byte, error) {
	command := c.Command
	if c.Environ == EnvironContainer {
		if c.parsedCommand=nil{
			c.parsedCommand=&ParsedCommand{}
		rootArgs := com.ParseArgs(command)
		if len(rootArgs) > 1 {
			c.parsedCommand.ContainerEngine, c.parsedCommand.ContainerName = ParseContainerInfo(rootArgs)
			c.parsedCommand.Command = rootArgs[0]
			c.parsedCommand.Args = rootArgs[1:]
		}
		}
		if len(c.parsedCommand.Command)>0{
			command=c.parsedCommand.Command
			rootArgs:=append([]string{},c.parsedCommand.Args...)
			args=append(rootArgs,args)
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

func (c *CommonConfig) RemoveDir(typeName string, rootDir string, prefix string, extensions ...string) error {
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if len(extensions) > 0 {
			ext := filepath.Ext(path)
			if !com.InSlice(ext, extensions) {
				return nil
			}
		}
		if len(prefix) > 0 && !strings.HasPrefix(info.Name(), prefix) {
			return nil
		}
		log.Info(`Delete the `+typeName+` file: `, path)
		return os.Remove(path)
	})
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
	} else {
		os.Remove(rootDir)
	}
	return err
}
