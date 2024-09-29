//go:build !bindata
// +build !bindata

/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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

package main

import (
	"path/filepath"

	"github.com/coscms/webcore/library/bindata"
	"github.com/coscms/webcore/library/buildinfo"
	"github.com/coscms/webcore/library/buildinfo/nginginfo"
)

func init() {
	nginginfo.SetNgingDir(buildinfo.NgingDir())
	bindata.Initialize()
	backendTemplateDir, err := filepath.Abs(filepath.Join(buildinfo.NgingPluginsDir(), `dockermanager/template/backend`))
	if err != nil {
		panic(err)
	}
	nginginfo.WatchTemplateDir(backendTemplateDir)
}