package engine

import "github.com/webx-top/echo"

const (
	EnvironLocal     = `local`
	EnvironContainer = `container`
)

var Environs = echo.NewKVData().Add(EnvironLocal, echo.T(`本机`)).Add(EnvironContainer, echo.T(`容器`))
