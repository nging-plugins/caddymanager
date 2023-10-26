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

package handler

import (
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/nging-plugins/caddymanager/application/dbschema"
	"github.com/nging-plugins/caddymanager/application/library/cmder"
	"github.com/nging-plugins/caddymanager/application/library/engine"
	"github.com/nging-plugins/caddymanager/application/library/form"
	"github.com/nging-plugins/caddymanager/application/model"
)

func DeleteCaddyfileByID(ctx echo.Context, serverIdent string, id uint) error {
	cfg, err := getServerConfig(ctx, serverIdent)
	if err != nil || cfg == nil {
		return err
	}
	saveDir, err := cfg.GetVhostConfigDirAbsPath()
	if err != nil {
		return err
	}
	saveFile := filepath.Join(saveDir, fmt.Sprint(id)+`.conf`)
	err = os.Remove(saveFile)
	if err == nil {
		item := engine.Engines.GetItem(cfg.Engine())
		if item != nil {
			err = item.X.(engine.Enginer).ReloadServer(ctx, cfg)
		}
	} else if os.IsNotExist(err) {
		err = nil
	}
	return err
}

var reSplitRegexp = regexp.MustCompile(`[\s]+`)

func hasEnvVar(v string) bool {
	for _, r := range v {
		if r == '$' || r == '%' {
			return true
		}
	}
	return false
}

func generateHostURL(currentHost string, hosts string) []template.HTML {
	hosts = strings.TrimSpace(hosts)
	hostsSlice := reSplitRegexp.Split(hosts, -1)
	urls := make([]template.HTML, 0, len(hostsSlice))
	for _, v := range hostsSlice {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		parsedValue := com.ParseEnvVar(v)
		if len(parsedValue) > 0 {
			switch {
			case parsedValue[0] == ':':
				v = `<a href="http://` + currentHost + parsedValue + `" target="_blank" rel="noopener noreferrer">` + v + `</a>`
			case strings.HasPrefix(parsedValue, `0.0.0.0:`):
				v = `<a href="http://` + currentHost + strings.TrimPrefix(parsedValue, `0.0.0.0`) + `" target="_blank" rel="noopener noreferrer">` + v + `</a>`
			case !strings.Contains(parsedValue, `//`):
				v = `<a href="http://` + parsedValue + `" target="_blank" rel="noopener noreferrer">` + v + `</a>`
			default:
				v = `<a href="` + strings.ReplaceAll(parsedValue, `*`, `test`) + `" target="_blank" rel="noopener noreferrer">` + v + `</a>`
			}
		}
		urls = append(urls, template.HTML(v))
	}
	return urls
}

func removeAllConf(rootDir string) error {
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
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

func getSaveDir(cfg engine.Configer) (saveDir string, err error) {
	if cfg == nil {
		cfg = cmder.GetCaddyConfig()
	}
	saveDir, err = cfg.GetVhostConfigDirAbsPath()
	if err != nil {
		return
	}
	err = com.MkdirAll(saveDir, os.ModePerm)
	return
}

func saveVhostConf(ctx echo.Context, cfg engine.Configer, id uint, values url.Values) error {
	if cfg == nil {
		cfg = cmder.GetCaddyConfig()
	}
	ctx.Set(`values`, form.NewValues(values))
	b, err := ctx.Fetch(`caddy/makeconfig/`+cfg.TemplateFile(), nil)
	if err != nil {
		return err
	}
	b = com.CleanSpaceLine(b)
	saveFile, err := cfg.GetVhostConfigDirAbsPath()
	if err != nil {
		return err
	}
	if !com.FileExists(saveFile) {
		com.MkdirAll(saveFile, os.ModePerm)
	}
	saveFile = filepath.Join(saveFile, fmt.Sprint(id)+`.conf`)
	log.Info(`Generate a Caddy configuration file: `, saveFile)
	err = os.WriteFile(saveFile, b, os.ModePerm)
	//jsonb, _ := caddyfile.ToJSON(b)
	//err = os.WriteFile(saveFile+`.json`, jsonb, os.ModePerm)
	return err
}

func getServerConfig(ctx echo.Context, serverIdent string) (engine.Configer, error) {
	var cfg engine.Configer
	if serverIdent == `default` {
		cfg = cmder.GetCaddyConfig()
	} else {
		svrM := model.NewVhostServer(ctx)
		err := svrM.Get(nil, `ident`, serverIdent)
		if err != nil {
			return cfg, err
		}
		item := engine.Engines.GetItem(svrM.Engine)
		if item == nil {
			return cfg, err
		}
		cfg = item.X.(engine.Enginer).BuildConfig(ctx, svrM.NgingVhostServer)
	}
	return cfg, nil
}

func saveVhostData(ctx echo.Context, m *dbschema.NgingVhost, values url.Values, restart bool) (err error) {
	var cfg engine.Configer
	cfg, err = getServerConfig(ctx, m.ServerIdent)
	if err != nil || cfg == nil {
		return
	}
	var saveDir string
	saveDir, err = getSaveDir(cfg)
	if err != nil {
		return
	}
	if m.Disabled == `Y` {
		saveFile := filepath.Join(saveDir, fmt.Sprint(m.Id)+`.conf`)
		if err = os.Remove(saveFile); os.IsNotExist(err) {
			err = nil
		}
	} else {
		err = saveVhostConf(ctx, cfg, m.Id, values)
	}
	if err == nil && restart {
		item := engine.Engines.GetItem(cfg.Engine())
		if item != nil {
			err = item.X.(engine.Enginer).ReloadServer(ctx, cfg)
		}
	}
	return
}

func receiveFormData(ctx echo.Context, m *dbschema.NgingVhost) {
	m.Domain = ctx.Form(`domain`)
	m.Disabled = ctx.Form(`disabled`)
	m.Root = ctx.Form(`root`)
	m.Name = ctx.Form(`name`)
	m.GroupId = ctx.Formx(`groupId`).Uint()
	m.ServerIdent = ctx.Form(`serverIdent`)
}
