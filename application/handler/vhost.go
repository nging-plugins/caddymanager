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
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/nging-plugins/caddymanager/application/dbschema"
	"github.com/nging-plugins/caddymanager/application/library/cmder"
	"github.com/nging-plugins/caddymanager/application/library/engine"
	"github.com/nging-plugins/caddymanager/application/library/form"
	"github.com/nging-plugins/caddymanager/application/model"
)

func VhostIndex(ctx echo.Context) error {
	m := model.NewVhost(ctx)
	groupID := ctx.Formx(`groupId`).Uint()
	serverIdent := ctx.Formx(`serverIdent`).String()
	engineType := ctx.Formx(`engine`).String()
	cond := db.Compounds{}
	if groupID > 0 {
		cond.AddKV(`a.group_id`, groupID)
	}
	if len(engineType) > 0 {
		if engineType == `default` {
			serverIdent = engineType
		} else {
			cond.AddKV(`b.engine`, engineType)
		}
	}
	if len(serverIdent) > 0 {
		cond.AddKV(`a.server_ident`, serverIdent)
	}
	q := ctx.Formx(`q`).String()
	if len(q) > 0 {
		cond.AddKV(`a.name`, db.Like(`%`+q+`%`))
	}
	var rowAndGroup []*model.VhostAndGroup
	p := m.NewParam().SetCols(`a.*`, `b.name AS serverName`, `b.engine AS serverEngine`).SetAlias(`a`).SetMW(func(r db.Result) db.Result {
		return r.OrderBy(`-a.id`)
	}).SetRecv(&rowAndGroup)
	p.AddJoin(`LEFT`, dbschema.NewNgingVhostServer(ctx).Short_(), `b`, `b.ident=a.server_ident`)
	p.AddArgs(cond.And())
	_, err := handler.PagingWithList(ctx, p)
	mg := dbschema.NewNgingVhostGroup(ctx)
	var groupList []*dbschema.NgingVhostGroup
	mg.ListByOffset(&groupList, nil, 0, -1)
	ctx.Set(`listData`, rowAndGroup)
	ctx.Set(`groupList`, groupList)
	ctx.Set(`engineList`, engine.Engines.Slice())
	ctx.Set(`groupId`, groupID)
	currentHost := ctx.Host()
	ctx.SetFunc(`generateHostURL`, func(hosts string) []template.HTML {
		return generateHostURL(currentHost, hosts)
	})
	ctx.SetFunc(`engineName`, engine.Engines.Get)
	return ctx.Render(`caddy/vhost`, handler.Err(ctx, err))
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

func Vhostbuild(ctx echo.Context) error {
	var err error
	configs := map[string]engine.Configer{}
	for _, v := range engine.Engines.Slice() {
		eng := v.X.(engine.Enginer)
		if eng == nil {
			continue
		}
		var rows []engine.Configer
		rows, err = eng.ListConfig(ctx)
		if err != nil {
			return err
		}
		for _, cfg := range rows {
			var saveDir string
			saveDir, err = getSaveDir(cfg)
			if err != nil {
				break
			}
			err = removeAllConf(saveDir)
			if err != nil {
				break
			}
			configs[cfg.Ident()] = cfg
		}
	}
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
	}
	m := model.NewVhost(ctx)
	n := 100
	cnt, err := m.ListByOffset(nil, nil, 0, n, `disabled`, `N`)
	if err != nil {
		return err
	}
	for i, j := 0, cnt(); int64(i) < j; i += n {
		if i > 0 {
			_, err = m.ListByOffset(nil, nil, i, n, `disabled`, `N`)
			if err != nil {
				handler.SendFail(ctx, err.Error())
				return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
			}
		}
		for _, m := range m.Objects() {
			var formData url.Values
			err := json.Unmarshal([]byte(m.Setting), &formData)
			if err == nil {
				cfg, ok := configs[m.ServerIdent]
				if ok {
					err = saveVhostConf(ctx, cfg, m.Id, formData)
				}
			}
			if err != nil {
				handler.SendFail(ctx, err.Error())
				return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
			}
		}
	}
	for _, cfg := range configs {
		item := engine.Engines.GetItem(cfg.Engine())
		if item == nil {
			continue
		}
		err = item.X.(engine.Enginer).ReloadServer(ctx, cfg)
		if err != nil {
			ctx.Logger().Error(err)
		}
	}
	handler.SendOk(ctx, ctx.T(`操作成功`))
	return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
}

func VhostAdd(ctx echo.Context) error {
	var err error
	m := model.NewVhost(ctx)
	if ctx.IsPost() {
		m.Domain = ctx.Form(`domain`)
		m.Disabled = ctx.Form(`disabled`)
		m.Root = ctx.Form(`root`)
		m.Name = ctx.Form(`name`)
		m.GroupId = ctx.Formx(`groupId`).Uint()
		var b []byte
		b, err = json.Marshal(ctx.Forms())
		switch {
		case err == nil:
			m.Setting = string(b)
			_, err = m.Insert()
			if err != nil {
				break
			}
			fallthrough
		case 0 == 1:
			err = saveVhostData(ctx, m.NgingVhost, ctx.Forms(), false)
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, db.Cond{`id`: id})
			if err == nil {
				var formData url.Values
				if e := json.Unmarshal([]byte(m.Setting), &formData); e == nil {
					for key, values := range formData {
						for k, v := range values {
							if k == 0 {
								ctx.Request().Form().Set(key, v)
								continue
							}
							ctx.Request().Form().Add(key, v)
						}
					}
				}
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}
	setVhostForm(ctx)
	ctx.Set(`isAdd`, true)
	ctx.Set(`title`, ctx.T(`添加网站`))
	return ctx.Render(`caddy/vhost_edit`, err)
}

func setVhostForm(ctx echo.Context) {
	ctx.SetFunc(`Val`, func(name, defaultValue string) string {
		return ctx.Form(name, defaultValue)
	})
	g := model.NewVhostGroup(ctx)
	g.ListByOffset(nil, nil, 0, -1)
	ctx.Set(`groupList`, g.Objects())
	svr := model.NewVhostServer(ctx)
	svr.ListByOffset(nil, nil, 0, -1, db.Cond{`disabled`: common.BoolN})
	ctx.Set(`serverList`, svr.Objects())
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
	b, err := ctx.Fetch(`caddy/`+cfg.TemplateFile(), nil)
	if err != nil {
		return err
	}
	b = com.CleanSpaceLine(b)
	saveFile, err := cfg.GetVhostConfigDirAbsPath()
	if err != nil {
		return err
	}
	saveFile = filepath.Join(saveFile, fmt.Sprint(id)+`.conf`)
	log.Info(`Generate a Caddy configuration file: `, saveFile)
	err = os.WriteFile(saveFile, b, os.ModePerm)
	//jsonb, _ := caddyfile.ToJSON(b)
	//err = os.WriteFile(saveFile+`.json`, jsonb, os.ModePerm)
	return err
}

func saveVhostData(ctx echo.Context, m *dbschema.NgingVhost, values url.Values, restart bool) (err error) {
	var saveDir string
	var cfg engine.Configer
	if m.ServerIdent == `default` {
		cfg = cmder.GetCaddyConfig()
	} else {
		svrM := model.NewVhostServer(ctx)
		err = svrM.Get(nil, `ident`, m.ServerIdent)
		if err != nil {
			return
		}
		item := engine.Engines.GetItem(svrM.Engine)
		if item == nil {
			return
		}
		cfg = item.X.(engine.Enginer).BuildConfig(ctx, svrM.NgingVhostServer)
		if cfg == nil {
			return
		}
	}
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

func VhostDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	if id < 1 {
		handler.SendFail(ctx, ctx.T(`id无效`))
		return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
	}
	m := model.NewVhost(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err != nil {
		handler.SendFail(ctx, err.Error())
	} else {
		err = DeleteCaddyfileByID(ctx, id)
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
		}
	}
	return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
}

func DeleteCaddyfileByID(ctx echo.Context, id uint) error {
	cfg := cmder.GetCaddyConfig()
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

func VhostEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	if id < 1 {
		handler.SendFail(ctx, ctx.T(`id无效`))
		return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
	}

	var err error
	m := model.NewVhost(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
	}
	if ctx.IsPost() {
		m.Domain = ctx.Form(`domain`)
		m.Disabled = ctx.Form(`disabled`)
		m.Root = ctx.Form(`root`)
		m.Name = ctx.Form(`name`)
		m.GroupId = ctx.Formx(`groupId`).Uint()
		var b []byte
		b, err = json.Marshal(ctx.Forms())
		switch {
		case err == nil:
			m.Setting = string(b)
			err = m.Update(nil, db.Cond{`id`: id})
			if err != nil {
				break
			}
			fallthrough
		case 0 == 1:
			removeCachedCert := ctx.Form(`removeCachedCert`)
			if len(removeCachedCert) > 0 && removeCachedCert == `1` {
				m.RemoveCachedCert()
			}
			err = saveVhostData(ctx, m.NgingVhost, ctx.Forms(), len(ctx.Form(`restart`)) > 0)
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
		}
	} else if ctx.IsAjax() {
		data := ctx.Data()
		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
			if !common.IsBoolFlag(disabled) {
				return ctx.NewError(code.InvalidParameter, ``).SetZone(`disabled`)
			}
			m.Disabled = disabled
			err = m.UpdateField(nil, `disabled`, disabled, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			if m.Disabled == `Y` {
				err = DeleteCaddyfileByID(ctx, id)
			} else {
				var formData url.Values
				err = json.Unmarshal([]byte(m.Setting), &formData)
				if err == nil {
					err = saveVhostData(ctx, m.NgingVhost, formData, true)
				}
			}
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
		}
		return ctx.JSON(data)
	} else {
		var formData url.Values
		if e := json.Unmarshal([]byte(m.Setting), &formData); e == nil {
			for key, values := range formData {
				for k, v := range values {
					if k == 0 {
						ctx.Request().Form().Set(key, v)
						continue
					}
					ctx.Request().Form().Add(key, v)
				}
			}
		}
		echo.StructToForm(ctx, m.NgingVhost, ``, echo.LowerCaseFirstLetter)
	}
	setVhostForm(ctx)
	ctx.Set(`activeURL`, `/caddy/vhost`)
	ctx.Set(`isAdd`, false)
	ctx.Set(`title`, ctx.T(`修改网站`))
	return ctx.Render(`caddy/vhost_edit`, err)
}
