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

// 也可以以服务的方式启动nging
// 服务支持的操作有：
// nging service install  	-- 安装服务
// nging service uninstall  -- 卸载服务
// nging service start 		-- 启动服务
// nging service stop 		-- 停止服务
// nging service restart 	-- 重启服务
package main

import (
	"path/filepath"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/cmd"
	_ "github.com/admpub/nging/v5/application/ico"
	_ "github.com/admpub/nging/v5/upgrade"

	"github.com/webx-top/com"

	//register

	_ "github.com/admpub/nging/v5/application"
	_ "github.com/admpub/nging/v5/application/initialize/manager"
	"github.com/admpub/nging/v5/application/library/buildinfo"
	"github.com/admpub/nging/v5/application/library/buildinfo/nginginfo"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/module"
	_ "github.com/admpub/nging/v5/application/library/sqlite"

	"github.com/admpub/nging/v5/application/version"

	// module
	"github.com/nging-plugins/caddymanager"
	"github.com/nging-plugins/servermanager"
)

var (
	BUILD_TIME string
	BUILD_OS   string
	BUILD_ARCH string
	CLOUD_GOX  string
	COMMIT     string
	LABEL      = `dev` //beta/alpha/stable
	VERSION    = `4.1.6`
	PACKAGE    = `free`

	schemaVer = version.DBSCHEMA //数据表结构版本
)

func main() {
	config.FromCLI().Conf = filepath.Join(buildinfo.NgingDir(), `config/config.yaml`)
	log.SetEmoji(com.IsMac)
	defer log.Close()
	buildinfo.New(
		buildinfo.Time(BUILD_TIME),
		buildinfo.OS(BUILD_OS),
		buildinfo.Arch(BUILD_ARCH),
		buildinfo.CloudGox(CLOUD_GOX),
		buildinfo.Commit(COMMIT),
		buildinfo.Label(LABEL),
		buildinfo.Version(VERSION),
		buildinfo.Package(PACKAGE),
		buildinfo.SchemaVer(schemaVer),
	).Apply()
	initModule()
	exec()
}

func exec() {
	cmd.Execute()
}

func initModule() {
	nginginfo.SetNgingPluginsDir(buildinfo.NgingPluginsDir())
	module.Register(
		&caddymanager.Module,
		&servermanager.Module,
	)
}
