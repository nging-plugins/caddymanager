package engine

import (
	"path/filepath"

	"github.com/nging-plugins/caddymanager/application/dbschema"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type Enginer interface {
	Name() string
	ListConfig(ctx echo.Context) ([]Configer, error)
	BuildConfig(ctx echo.Context, m *dbschema.NgingVhostServer) Configer
	ReloadServer(ctx echo.Context, cfg Configer) error
	DefaultConfigDir() string
}

var Engines = echo.NewKVData()

func Thirdparty() []echo.KV {
	r := make([]echo.KV, 0, len(Engines.Slice())-1)
	for _, v := range Engines.Slice() {
		if v.K == `default` {
			continue
		}
		r = append(r, *v)
	}
	return r
}

func certSaveDir(dir string) string {
	if com.IsWindows {
		return filepath.Join(echo.Wd(), `data`, dir)
	}
	return filepath.Join(`/etc/`, dir)
}

func BuildCertSaveDir(dir string) string {
	return certSaveDir(dir)
}
