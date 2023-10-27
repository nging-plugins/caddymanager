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
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/nging-plugins/caddymanager/application/library/engine"
	"github.com/nging-plugins/caddymanager/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func ServerIndex(ctx echo.Context) error {
	m := model.NewVhostServer(ctx)
	cond := db.NewCompounds()
	err := m.ListPage(cond)
	ctx.Set(`listData`, m.Objects())
	ctx.SetFunc(`engineName`, engine.Engines.Get)
	ctx.SetFunc(`environName`, engine.Environs.Get)
	return ctx.Render(`caddy/server`, handler.Err(ctx, err))
}

func ServerAdd(ctx echo.Context) error {
	var err error
	m := model.NewVhostServer(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingVhostServer, echo.ExcludeFieldName(`created`, `updated`))
		if err == nil {
			_, err = m.Add()
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/caddy/server`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, db.Cond{`id`: id})
			if err == nil {
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}
	ctx.Set(`activeURL`, `/caddy/server`)
	ctx.Set(`isAdd`, true)
	setServerForm(ctx)
	return ctx.Render(`caddy/server_edit`, handler.Err(ctx, err))
}

func setServerForm(ctx echo.Context) {
	thirdpartyEngines := engine.Thirdparty()
	ctx.Set(`engineList`, thirdpartyEngines)
	configDirs := map[string]string{}
	for _, eng := range thirdpartyEngines {
		configDirs[eng.K] = eng.X.(engine.Enginer).DefaultConfigDir()
	}
	ctx.Set(`configDirs`, configDirs)
	ctx.Set(`environList`, engine.Environs.Slice())
}

func ServerEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewVhostServer(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/caddy/server`))
	}
	if ctx.IsPost() {
		old := *m.NgingVhostServer
		err = ctx.MustBind(m.NgingVhostServer, echo.ExcludeFieldName(`created`, `updated`, `engine`, `ident`))
		if err != nil {
			goto END
		}
		m.Id = id
		err = m.Edit(nil, `id`, id)
		if err != nil {
			goto END
		}
		if old.Disabled != m.Disabled {
			if m.Disabled == `Y` {
				err = deleteCaddyfileByServer(ctx, &old, true)
			} else {
				err = vhostbuild(ctx, 0, ``, ``, m.NgingVhostServer)
			}
			if err != nil {
				ctx.Logger().Error(err)
			}
		} else if old.VhostConfigDir != m.VhostConfigDir {
			err = deleteCaddyfileByServer(ctx, &old, true)
			if err != nil {
				ctx.Logger().Error(err)
			}
			err = vhostbuild(ctx, 0, ``, ``, m.NgingVhostServer)
			if err != nil {
				ctx.Logger().Error(err)
			}
		}
		handler.SendOk(ctx, ctx.T(`修改成功`))
		return ctx.Redirect(handler.URLFor(`/caddy/server`))
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
				err = deleteCaddyfileByServer(ctx, m.NgingVhostServer, true)
			} else {
				err = vhostbuild(ctx, 0, ``, ``, m.NgingVhostServer)
			}
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
		}
		return ctx.JSON(data)
	}
	echo.StructToForm(ctx, m.NgingVhostServer, ``, echo.LowerCaseFirstLetter)

END:
	ctx.Set(`activeURL`, `/caddy/server`)
	ctx.Set(`isAdd`, false)
	setServerForm(ctx)
	return ctx.Render(`caddy/server_edit`, handler.Err(ctx, err))
}

func ServerDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewVhostServer(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		err = deleteCaddyfileByServer(ctx, m.NgingVhostServer, true)
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}
	return ctx.Redirect(handler.URLFor(`/caddy/server`))
}
