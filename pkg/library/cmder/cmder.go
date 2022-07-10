package cmder

import (
	"sort"

	"github.com/admpub/nging/v4/application/library/config/cmder"

	"github.com/nging-plugins/caddymanager/pkg/library/caddy"
)

type NewServerCmd func() ServerCmder

type ServerDef struct {
	Ident   string
	Name    string
	newFunc NewServerCmd
}

var serverCmders = map[string]*ServerDef{}
var serverType = `embedwebserver`

func SetServerType(ident string) {
	serverType = ident
}

type ServerCmder interface {
	cmder.Cmder
	RemoveCachedCert(string) error
	LogFile() string
}

func Register(ident string, name string, newServerCmd NewServerCmd) {
	serverCmders[ident] = &ServerDef{Ident: ident, Name: name, newFunc: newServerCmd}
}

func Initer() interface{} {
	return &caddy.Config{}
}

func Get() ServerCmder {
	return cmder.Get(serverType).(ServerCmder)
}

func New() ServerCmder {
	if _, ok := serverCmders[serverType]; !ok {
		return nil
	}
	return serverCmders[serverType].newFunc()
}

func ServerDefs() []ServerDef {
	names := make([]string, 0, len(serverCmders))
	for name := range serverCmders {
		names = append(names, name)
	}
	sort.Strings(names)
	defs := make([]ServerDef, len(names))
	for index, name := range names {
		defs[index] = *serverCmders[name]
	}
	return defs
}
