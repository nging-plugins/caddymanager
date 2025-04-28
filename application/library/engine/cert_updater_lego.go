package engine

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/admpub/once"
	"github.com/nging-plugins/caddymanager/application/library/github"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"golang.org/x/net/idna"
)

var legoCommand string
var legoOnce once.Once

func initLegoCommand() {
	if _, err := exec.LookPath(`lego`); err == nil {
		legoCommand = `lego`
		return
	}
	legoBinary := filepath.Join(echo.Wd(), `support`, `lego`, `lego`)
	if com.IsWindows {
		legoBinary += `.exe`
	}
	if com.IsFile(legoBinary) {
		legoCommand = legoBinary
		return
	}
	legoCommand = `lego`
	legoOnce.Reset()
}

func LegoCommand() string {
	legoOnce.Do(initLegoCommand)
	return legoCommand
}

var legoPathFormat = CertPathFormat{
	Cert:  certSaveDir(`letsencrypt/.lego/certificates/{domain}.crt`),
	Key:   certSaveDir(`letsencrypt/.lego/certificates/{domain}.key`),
	Trust: ``,
}

func init() {
	CertUpdaters.Add(`lego`, `Lego`, echo.KVxOptX[CertUpdater, any](CertUpdater{
		MakeCommand:     MakeLegoCommand,
		Update:          RenewCertByLego,
		PathFormat:      legoPathFormat,
		DomainSanitizer: LegoSanitizedDomain,
	}))
}

// 申请：
// lego --accept-tos --email you@example.com --http --http.webroot /path/to/webroot --domains example.com run
// https://go-acme.github.io/lego/usage/cli/obtain-a-certificate/
// 更新：
// lego --email="you@example.com" --domains="example.com" --http renew
// https://go-acme.github.io/lego/usage/cli/renew-a-certificate/
//
// CLOUDFLARE_EMAIL="you@example.com" \
// CLOUDFLARE_API_KEY="yourprivatecloudflareapikey" \
// lego --email "you@example.com" --dns cloudflare --domains "example.org" run
func MakeLegoCommand(data RequestCertUpdate) (command string, args []string, env []string) {
	command = LegoCommand()
	args = []string{command}
	saveDir := certSaveDir(`letsencrypt`)
	if len(data.CertSaveDir) > 0 {
		saveDir = data.CertSaveDir
	}
	args = append(args, `--path`, saveDir)
	args = append(args, `--email`, data.Email)
	if len(data.DNSProvider) > 0 {
		args = append(args, `--dns`, data.DNSProvider)
	} else {
		args = append(args, `--http`)
	}
	for _, domain := range data.Domains {
		args = append(args, `--domains`, domain)
	}
	if data.Obtain {
		args = append(args,
			`--http.webroot`, data.CertVerifyDir,
			`--agree-tos`,
			`--pem`,
			`run`,
		)
	} else {
		args = append(args, `renew`)
	}
	if len(data.CmdPathPrefix) > 0 {
		rootArgs := com.ParseArgs(data.CmdPathPrefix)
		if len(rootArgs) > 1 {
			command = rootArgs[0]
			args = append(rootArgs[1:], args...)
		}
	}
	env = data.Env
	return
}

func RenewCertByLego(ctx context.Context, data RequestCertUpdate) error {
	if len(data.Domains) == 0 {
		return nil
	}
	command, args, env := MakeLegoCommand(data)
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, env...)
	result, err := cmd.CombinedOutput()
	//log.Okay(cmd.String())
	if err != nil {
		err = fmt.Errorf(`%s: %w`, result, err)
	}
	return err
}

var (
	domainReplacer = strings.NewReplacer(
		":", "-",
		"*", "_",
	)
)

// LegoSanitizedDomain Make sure no funny chars are in the cert names (like wildcards ;)).
func LegoSanitizedDomain(domain string) (string, error) {
	return idna.ToASCII(domainReplacer.Replace(domain))
}

func LegoInstall() error {
	savePath, err := github.Download(`go-acme`, `lego`)
	if err != nil {
		return fmt.Errorf("download error: %w", err)
	}
	return legoInstall(savePath)
}

func legoInstall(tarFile string) error {
	installDir := filepath.Join(echo.Wd(), `support`, `lego`)
	err := os.MkdirAll(installDir, 0755)
	if err != nil {
		return fmt.Errorf("mkdir error: %w", err)
	}
	_, err = com.UnTarGz(tarFile, installDir)
	if err != nil {
		return fmt.Errorf("untar error: %w", err)
	}
	binaryPath := filepath.Join(installDir, `lego`)
	if com.IsWindows {
		binaryPath += `.exe`
	}
	err = os.Chmod(binaryPath, 0755)
	return err
}
