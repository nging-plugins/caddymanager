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
	"net/http"

	"github.com/admpub/log"
	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/model"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func ValidAddonName(addon string) bool {
	return com.IsAlphaNumericUnderscore(addon)
}

func AddonIndex(ctx echo.Context) error {
	return ctx.Render(`caddy/addon/index`, nil)
}

func AddonForm(ctx echo.Context) error {
	addon := ctx.Query(`addon`)
	if len(addon) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, ctx.T("参数 addon 的值不能为空"))
	}
	if !ValidAddonName(addon) {
		return echo.NewHTTPError(http.StatusBadRequest, ctx.T("参数 addon 的值包含非法字符"))
	}
	ctx.SetFunc(`Val`, func(name, defaultValue string) string {
		return defaultValue
	})
	ctx.SetFunc(`Key`, func(name string) string {
		return name
	})
	ctx.SetFunc(`ListS3Accounts`, func() []*dbschema.NgingCloudStorage {
		return listS3Accounts(ctx)
	})
	index := ctx.Queryx(`index`, `0`).Int()
	return ctx.Render(`caddy/addon/form/`+addon, index)
}

func listS3Accounts(ctx echo.Context) []*dbschema.NgingCloudStorage {
	m := model.NewCloudStorage(ctx)
	_, err := m.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, 0, -1)
	if err != nil {
		log.Warn(err)
		return nil
	}
	return m.Objects()
}
