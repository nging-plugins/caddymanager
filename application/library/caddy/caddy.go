/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package caddy

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gitee.com/admpub/certmagic"
	"github.com/admpub/caddy"
	_ "github.com/admpub/caddy/caddyhttp"
	"github.com/admpub/caddy/caddytls"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/msgbox"
	"github.com/admpub/nging/v5/application/library/common"
)

const UnsupportedPathInfo = `unsupported-nging-default-webserver-config`

var (
	DefaultConfig = &Config{
		Agreed:                  true,
		CAUrl:                   certmagic.Default.CA,
		CATimeout:               int64(certmagic.HTTPTimeout.Seconds()),
		DisableHTTPChallenge:    certmagic.Default.DisableHTTPChallenge,
		DisableTLSALPNChallenge: certmagic.Default.DisableTLSALPNChallenge,
		ServerType:              `http`,
		CPU:                     `80%`,
		PidFile:                 `./caddy.pid`,
		AppName:                 `nging`,
	}
	DefaultVersion = `2.0.0`
	EnableReload   = true
)

func TrapSignals() {
	caddy.TrapSignals()
}

func SetDefaults(c *Config) {
	if len(c.CAUrl) == 0 {
		c.CAUrl = DefaultConfig.CAUrl
	}
	if c.CATimeout == 0 {
		c.CATimeout = DefaultConfig.CATimeout
	}
	if len(c.ServerType) == 0 {
		c.ServerType = DefaultConfig.ServerType
	}
	if len(c.CPU) == 0 {
		c.CPU = DefaultConfig.CPU
	}
	pidFile := filepath.Join(echo.Wd(), `data` + echo.FilePathSeparator + `pid`)
	err := com.MkdirAll(pidFile, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	pidFile = filepath.Join(pidFile, `caddy.pid`)
	c.PidFile = pidFile
	if len(c.LogFile) == 0 {
		logFile := filepath.Join(echo.Wd(), `data` + echo.FilePathSeparator + `logs`)
		err := com.MkdirAll(logFile, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
		c.LogFile = filepath.Join(logFile, `caddy.log`)
	} else {
		c.LogFile = common.OSAbsPath(c.LogFile)
		err := com.MkdirAll(filepath.Dir(c.LogFile), os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	if len(c.AppName) == 0 {
		c.AppName = DefaultConfig.AppName
	}
	c.AppVersion = echo.String(`VERSION`, DefaultVersion)
	c.Agreed = true
	c.ctx, c.cancel = context.WithCancel(context.Background())
}

type Config struct {
	Agreed                  bool   `json:"agreed"` //Agree to the CA's Subscriber Agreement
	CAUrl                   string `json:"caURL"`  //URL to certificate authority's ACME server directory
	DisableHTTPChallenge    bool   `json:"disableHTTPChallenge"`
	DisableTLSALPNChallenge bool   `json:"disableTLSALPNChallenge"`
	Caddyfile               string `json:"caddyFile,omitempty"` //Caddyfile to load (default caddy.DefaultConfigFile)
	CPU                     string `json:"cpu"`                 //CPU cap
	CAEmail                 string `json:"caEmail"`             //Default ACME CA account email address
	CATimeout               int64  `json:"caTimeout"`           //Default ACME CA HTTP timeout
	LogFile                 string `json:"logFile"`             //Process log file
	PidFile                 string `json:"-"`                   //Path to write pid file
	Quiet                   bool   `json:"quiet"`               //Quiet mode (no initialization output)
	Revoke                  string `json:"revoke"`              //Hostname for which to revoke the certificate
	ServerType              string `json:"serverType"`          //Type of server to run
	VhostConfigDir          string `json:"vhostConfigDir"`

	//---
	EnvFile string `json:"envFile"` //Path to file with environment variables to load in KEY=VALUE format
	Plugins bool   `json:"plugins"` //List installed plugins
	Version bool   `json:"version"` //Show version

	//---
	AppVersion            string
	AppName               string
	instance              *caddy.Instance
	stopped               bool
	ctx                   context.Context
	cancel                context.CancelFunc
	caddyfileAbsPath      string
	vhostConfigDirAbsPath string
}

func now() string {
	return time.Now().Format(`2006-01-02 15:04:05`)
}

func DefaultConfigDir() string {
	return filepath.Join(config.FromCLI().ConfDir(), `vhosts`)
}

func (c *Config) GetVhostConfigLocalDirAbs() (string, error) {
	var err error
	if len(c.vhostConfigDirAbsPath) == 0 {
		if len(c.VhostConfigDir) > 0 {
			c.vhostConfigDirAbsPath, err = filepath.Abs(c.VhostConfigDir)
		} else {
			c.vhostConfigDirAbsPath = DefaultConfigDir()
		}
	}
	return c.vhostConfigDirAbsPath, err
}

func (c *Config) GetTemplateFile() string {
	return `caddyfile`
}

func (c *Config) GetIdent() string {
	return `default`
}

func (c *Config) GetEngine() string {
	return `default`
}

func (c *Config) GetEnviron() string {
	return `local`
}

func (c *Config) GetCertLocalDir() string {
	return UnsupportedPathInfo
}

func (c *Config) GetCertContainerDir() string {
	return UnsupportedPathInfo
}

func (c *Config) GetVhostConfigLocalDir() string {
	return UnsupportedPathInfo
}

func (c *Config) GetVhostConfigContainerDir() string {
	return UnsupportedPathInfo
}

func (c *Config) GetEngineConfigLocalFile() string {
	return UnsupportedPathInfo
}

func (c *Config) GetEngineConfigContainerFile() string {
	return UnsupportedPathInfo
}

func (c *Config) setDefaultCaddyfile() (err error) {
	var saveDir string
	saveDir, err = c.GetVhostConfigLocalDirAbs()
	if err != nil {
		return err
	}
	importPattern := saveDir + echo.FilePathSeparator + `*.conf`
	c.caddyfileAbsPath, err = filepath.Abs(importPattern)
	return
}

func (c *Config) fixedCaddyfile() error {
	if len(c.caddyfileAbsPath) > 0 {
		return nil
	}
	c.caddyfileAbsPath = c.Caddyfile
	if len(c.caddyfileAbsPath) == 0 {
		return c.setDefaultCaddyfile()
	}
	if c.caddyfileAbsPath == `stdin` {
		return nil
	}
	var err error
	if strings.Contains(c.caddyfileAbsPath, `*`) {
		c.caddyfileAbsPath, err = filepath.Abs(c.caddyfileAbsPath)
		return err
	}
	c.caddyfileAbsPath, err = filepath.Abs(c.caddyfileAbsPath)
	if err != nil {
		return err
	}
	if !com.FileExists(c.caddyfileAbsPath) {
		return c.setDefaultCaddyfile()
	}
	b, err := os.ReadFile(c.caddyfileAbsPath)
	if err != nil {
		return c.setDefaultCaddyfile()
	}
	b = bytes.TrimSpace(b)
	if bytes.Equal(b, []byte(`import ./config/vhosts/*.conf`)) {
		var actualDir string
		actualDir, err = c.GetVhostConfigLocalDirAbs()
		if err != nil {
			return err
		}
		expectedDir, _ := filepath.Abs(`./config/vhosts`)
		if actualDir != expectedDir {
			err = c.setDefaultCaddyfile()
		}
	}
	return err
}

func (c *Config) Start() error {
	if err := c.fixedCaddyfile(); err != nil {
		return err
	}
	caddy.AppName = c.AppName
	caddy.AppVersion = c.AppVersion
	certmagic.UserAgent = c.AppName + "/" + c.AppVersion
	c.stopped = false

	// Executes Startup events
	caddy.EmitEvent(caddy.StartupEvent, nil)

	// Get Caddyfile input
	caddyfile, err := caddy.LoadCaddyfile(c.ServerType)
	if err != nil {
		return err
	}

	if EnableReload {
		c.watchingSignal()
	}

	// Start your engines
	c.instance, err = caddy.Start(caddyfile)
	if err != nil {
		return err
	}

	msgbox.Success(`Caddy`, `Server has been successfully started at `+now())

	// Twiddle your thumbs
	c.instance.Wait()
	return nil
}

func (c *Config) RemoveCachedCert(domain string) {
	certKey := certmagic.StorageKeys.SiteCert(c.CAUrl, domain)
	keyKey := certmagic.StorageKeys.SitePrivateKey(c.CAUrl, domain)
	metaKey := certmagic.StorageKeys.SiteMeta(c.CAUrl, domain)
	os.Remove(certKey)
	os.Remove(keyKey)
	os.Remove(metaKey)
}

// Listen to keypress of "return" and restart the app automatically
func (c *Config) watchingSignal() {
	debug := false
	go func() {
		if debug {
			msgbox.Info(`Caddy`, `listen return ==> `+now())
		}
		in := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				if debug {
					msgbox.Info(`Caddy`, `reading ==> `+now())
				}
				input, err := in.ReadString(com.LF)
				if err != nil && err != io.EOF {
					log.Println(`watchingSignal:`, err.Error())
					return
				}
				if input == com.StrLF || input == com.StrCRLF {
					if debug {
						msgbox.Info(`Caddy`, `restart ==> `+now())
					}
					c.Restart()
				} else {
					if debug {
						msgbox.Info(`Caddy`, `waiting ==> `+now())
					}
				}
			}
		}
	}()
}

func (c *Config) Restart() error {
	defer msgbox.Success(`Caddy`, `Server has been successfully reloaded at `+now())
	if c.instance == nil {
		return nil
	}
	// Get Caddyfile input
	caddyfile, err := caddy.LoadCaddyfile(c.ServerType)
	if err != nil {
		return err
	}
	c.instance, err = c.instance.Restart(caddyfile)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Stop() error {
	c.stopped = true
	defer func() {
		c.cancel()
		msgbox.Success(`Caddy`, `Server has been successfully stopped at `+now())
	}()
	if c.instance == nil {
		return nil
	}
	return c.instance.Stop()
}

func (c *Config) Init() *Config {
	certmagic.Default.Agreed = c.Agreed
	certmagic.Default.CA = c.CAUrl
	certmagic.Default.DisableHTTPChallenge = c.DisableHTTPChallenge
	certmagic.Default.DisableTLSALPNChallenge = c.DisableTLSALPNChallenge
	certmagic.Default.Email = c.CAEmail
	certmagic.HTTPTimeout = time.Duration(c.CATimeout) * time.Second

	caddy.PidFile = c.PidFile
	caddy.Quiet = c.Quiet
	caddy.RegisterCaddyfileLoader("flag", caddy.LoaderFunc(c.confLoader))
	caddy.SetDefaultCaddyfileLoader("default", caddy.LoaderFunc(c.defaultLoader))

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// Set up process log before anything bad happens
	switch c.LogFile {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	case "":
		log.SetOutput(io.Discard)
	default:
		log.SetOutput(&lumberjack.Logger{
			Filename:   c.LogFile,
			MaxSize:    100,
			MaxAge:     14,
			MaxBackups: 10,
		})
	}

	//Load all additional envs as soon as possible
	if err := LoadEnvFromFile(c.EnvFile); err != nil {
		mustLogFatalf("%v", err)
	}

	// Check for one-time actions
	if len(c.Revoke) > 0 {
		err := caddytls.Revoke(c.Revoke)
		if err != nil {
			mustLogFatalf(err.Error())
		}
		fmt.Printf("Revoked certificate for %s\n", c.Revoke)
		os.Exit(0)
	}
	if c.Version {
		fmt.Printf("%s %s\n", c.AppName, c.AppVersion)
		os.Exit(0)
	}
	if c.Plugins {
		fmt.Println(caddy.DescribePlugins())
		os.Exit(0)
	}

	// Set CPU cap
	err := setCPU(c.CPU)
	if err != nil {
		mustLogFatalf(err.Error())
	}
	return c
}

// confLoader loads the Caddyfile using the -conf flag.
func (c *Config) confLoader(serverType string) (caddy.Input, error) {
	if c.caddyfileAbsPath == "" {
		return nil, nil
	}
	if c.caddyfileAbsPath == "stdin" {
		return caddy.CaddyfileFromPipe(os.Stdin, serverType)
	}
	var contents []byte
	if strings.Contains(c.caddyfileAbsPath, "*") {
		// Let caddyfile.doImport logic handle the globbed path
		contents = []byte("import " + c.caddyfileAbsPath)
	} else {
		var err error
		contents, err = os.ReadFile(c.caddyfileAbsPath)
		if err != nil {
			return nil, err
		}
	}
	return caddy.CaddyfileInput{
		Contents:       contents,
		Filepath:       c.caddyfileAbsPath,
		ServerTypeName: serverType,
	}, nil
}

// defaultLoader loads the Caddyfile from the current working directory.
func (c *Config) defaultLoader(serverType string) (caddy.Input, error) {
	contents, err := os.ReadFile(caddy.DefaultConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return caddy.CaddyfileInput{
		Contents:       contents,
		Filepath:       caddy.DefaultConfigFile,
		ServerTypeName: serverType,
	}, nil
}

// mustLogFatalf wraps log.Fatalf() in a way that ensures the
// output is always printed to stderr so the user can see it
// if the user is still there, even if the process log was not
// enabled. If this process is an upgrade, however, and the user
// might not be there anymore, this just logs to the process
// log and exits.
func mustLogFatalf(format string, args ...interface{}) {
	if !caddy.IsUpgrade() {
		log.SetOutput(os.Stderr)
	}
	log.Fatalf(format, args...)
}

// setCPU parses string cpu and sets GOMAXPROCS
// according to its value. It accepts either
// a number (e.g. 3) or a percent (e.g. 50%).
func setCPU(cpu string) error {
	var numCPU int

	availCPU := runtime.NumCPU()

	if strings.HasSuffix(cpu, "%") {
		// Percent
		var percent float32
		pctStr := cpu[:len(cpu)-1]
		pctInt, err := strconv.Atoi(pctStr)
		if err != nil || pctInt < 1 || pctInt > 100 {
			return errors.New("invalid CPU value: percentage must be between 1-100")
		}
		percent = float32(pctInt) / 100
		numCPU = int(float32(availCPU) * percent)
		if numCPU < 1 {
			numCPU = 1
		}
	} else {
		// Number
		num, err := strconv.Atoi(cpu)
		if err != nil || num < 1 {
			return errors.New("invalid CPU value: provide a number or percent greater than 0")
		}
		numCPU = num
	}

	if numCPU > availCPU {
		numCPU = availCPU
	}

	runtime.GOMAXPROCS(numCPU)
	return nil
}

// LoadEnvFromFile loads additional envs if file provided and exists
// Envs in file should be in KEY=VALUE format
func LoadEnvFromFile(envFile string) error {
	if envFile == "" {
		return nil
	}

	file, err := os.Open(envFile)
	if err != nil {
		return err
	}
	defer file.Close()

	envMap, err := ParseEnvFile(file)
	if err != nil {
		return err
	}

	for k, v := range envMap {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}

	return nil
}

// ParseEnvFile implements parse logic for environment files
func ParseEnvFile(envInput io.Reader) (map[string]string, error) {
	envMap := make(map[string]string)

	scanner := bufio.NewScanner(envInput)
	var line string
	lineNumber := 0

	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		lineNumber++

		// skip lines starting with comment
		if strings.HasPrefix(line, "#") {
			continue
		}

		// skip empty line
		if len(line) == 0 {
			continue
		}

		fields := strings.SplitN(line, "=", 2)
		if len(fields) != 2 {
			return nil, fmt.Errorf("Can't parse line %d; line should be in KEY=VALUE format", lineNumber)
		}

		if strings.Contains(fields[0], " ") {
			return nil, fmt.Errorf("Can't parse line %d; KEY contains whitespace", lineNumber)
		}

		key := fields[0]
		val := fields[1]

		if key == "" {
			return nil, fmt.Errorf("Can't parse line %d; KEY can't be empty string", lineNumber)
		}
		envMap[key] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envMap, nil
}
