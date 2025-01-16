package caddymanager

import (
	"github.com/coscms/webcore/library/config/cmder"
	"github.com/coscms/webcore/library/config/extend"
	"github.com/coscms/webcore/library/module"
	"github.com/coscms/webcore/library/nlog"

	"github.com/nging-plugins/caddymanager/application/handler"
	pluginCmder "github.com/nging-plugins/caddymanager/application/library/cmder"
	"github.com/nging-plugins/caddymanager/application/library/setup"

	_ "github.com/nging-plugins/caddymanager/application/library/thirdparty/caddy2/cmder"
	_ "github.com/nging-plugins/caddymanager/application/library/thirdparty/nginx/cmder"
)

const ID = `caddy`

var Module = module.Module{
	Startup: pluginCmder.Name,
	Extend: map[string]extend.Initer{
		ID: pluginCmder.Initer,
	},
	Cmder: map[string]cmder.Cmder{
		ID: pluginCmder.New(),
	},
	TemplatePath: map[string]string{
		ID: `caddymanager/template/backend`,
	},
	AssetsPath: []string{
		`caddymanager/public/assets/backend`,
	},
	SQLCollection: setup.RegisterSQL,
	Dashboard:     RegisterDashboard,
	Navigate:      RegisterNavigate,
	Route:         handler.RegisterRoute,
	LogParser: map[string]nlog.LogParser{
		`access`: handler.ParseTailLine,
	},
	DBSchemaVer: 0.3000,
}
