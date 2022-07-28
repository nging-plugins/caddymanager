package caddymanager

import (
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config/cmder"
	"github.com/admpub/nging/v4/application/library/config/extend"
	"github.com/admpub/nging/v4/application/library/module"

	"github.com/nging-plugins/caddymanager/pkg/handler"
	pluginCmder "github.com/nging-plugins/caddymanager/pkg/library/cmder"
	_ "github.com/nging-plugins/caddymanager/pkg/library/cmder/embedcaddy"
	"github.com/nging-plugins/caddymanager/pkg/library/setup"
)

const ID = `caddy`

var Module = module.Module{
	Startup: ID,
	Extend: map[string]extend.Initer{
		ID: pluginCmder.Initer,
	},
	Cmder: map[string]cmder.Cmder{
		ID: pluginCmder.New(),
	},
	TemplatePath: map[string]string{
		ID: `caddymanager/template/backend`,
	},
	AssetsPath:    []string{},
	SQLCollection: setup.RegisterSQL,
	Dashboard:     RegisterDashboard,
	Navigate:      RegisterNavigate,
	Route:         handler.RegisterRoute,
	LogParser: map[string]common.LogParser{
		`access`: handler.ParseTailLine,
	},
	DBSchemaVer: 0.0000,
}
