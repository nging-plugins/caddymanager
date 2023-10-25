package model

import (
	"github.com/nging-plugins/caddymanager/application/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func NewVhostServer(ctx echo.Context) *VhostServer {
	return &VhostServer{
		NgingVhostServer: dbschema.NewNgingVhostServer(ctx),
	}
}

type VhostServer struct {
	*dbschema.NgingVhostServer
}

func (f *VhostServer) check() error {
	ctx := f.Context()
	if len(f.Name) == 0 {
		return ctx.NewError(code.InvalidParameter, `名称不能为空`).SetZone(`name`)
	}
	if len(f.Ident) == 0 {
		return ctx.NewError(code.InvalidParameter, `唯一标识不能为空`).SetZone(`ident`)
	}
	var exists bool
	var err error
	if f.Id > 0 {
		exists, err = f.Exists(f.Ident, f.Id)
	} else {
		exists, err = f.Exists(f.Ident)
	}
	if err != nil {
		return err
	}
	if exists {
		return ctx.NewError(code.DataAlreadyExists, `唯一标识“%s”已经存在`, f.Ident).SetZone(`ident`)
	}
	return nil
}

func (f *VhostServer) Exists(ident string, exclude ...uint) (bool, error) {
	cond := db.NewCompounds()
	cond.AddKV(`ident`, ident)
	if len(exclude) > 0 {
		cond.AddKV(`id`, db.NotEq(exclude[0]))
	}
	return f.NgingVhostServer.Exists(nil, cond)
}

func (f *VhostServer) Add() (interface{}, error) {
	if err := f.check(); err != nil {
		return nil, err
	}
	return f.NgingVhostServer.Insert()
}

func (f *VhostServer) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	if err = f.check(); err != nil {
		return err
	}
	return f.NgingVhostServer.Update(mw, args...)
}

func (f *VhostServer) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	return f.NgingVhostServer.Delete(mw, args...)
}
