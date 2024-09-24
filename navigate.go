package caddymanager

import (
	"github.com/coscms/webcore/library/module"
	"github.com/nging-plugins/caddymanager/application/handler"
	_ "github.com/nging-plugins/caddymanager/application/library/cmder"
	_ "github.com/nging-plugins/caddymanager/application/library/setup"
)

var LeftNavigate = handler.LeftNavigate

func RegisterNavigate(nc module.Navigate) {
	nc.Backend().AddLeftItems(-1, LeftNavigate)
}
