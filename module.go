package caddymanager

import (
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config/cmder"
	"github.com/admpub/nging/v4/application/library/config/extend"
	"github.com/admpub/nging/v4/application/library/module"

	"github.com/nging-plugins/caddymanager/pkg/handler"
	pluginCmder "github.com/nging-plugins/caddymanager/pkg/library/cmder"
	"github.com/nging-plugins/caddymanager/pkg/library/setup"
)

var Module = module.Module{
	Startup: `caddy`,
	Extend: map[string]extend.Initer{
		`caddy`: pluginCmder.Initer,
	},
	Cmder: map[string]cmder.Cmder{
		`caddy`: pluginCmder.New(),
	},
	TemplatePath: map[string]string{
		`caddy`: `caddymanager/template/backend`,
	},
	AssetsPath:    []string{},
	SQLCollection: setup.RegisterSQL,
	Dashboard:     RegisterDashboard,
	Navigate:      RegisterRoute,
	Route:         handler.RegisterRoute,
	LogParser: map[string]common.LogParser{
		`access`: handler.ParseTailLine,
	},
}
