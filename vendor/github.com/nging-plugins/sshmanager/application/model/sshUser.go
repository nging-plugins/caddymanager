/*
Nging is a toolbox for webmasters
Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

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

package model

import (
	"io"
	"strings"

	"github.com/webx-top/echo"

	"github.com/admpub/go-sshclient"
	webTerminalSSH "github.com/admpub/web-terminal/library/ssh"
	"github.com/nging-plugins/sshmanager/application/dbschema"
	sshconf "github.com/nging-plugins/sshmanager/application/library/config"
)

var (
	Decode = func(r string) string { return r }
)

type SshUserAndGroup struct {
	*dbschema.NgingSshUser
	Group *dbschema.NgingSshUserGroup
}

func NewSshUser(ctx echo.Context) *SshUser {
	return &SshUser{
		NgingSshUser: dbschema.NewNgingSshUser(ctx),
	}
}

type SshUser struct {
	*dbschema.NgingSshUser
}

func (s *SshUser) ExecMultiCMD(writer io.Writer, commands ...string) error {
	return ExecMultiCMD(s.NgingSshUser, writer, commands...)
}

func (s *SshUser) Connect() (*webTerminalSSH.SSH, error) {
	return Connect(s.NgingSshUser)
}

func ExecMultiCMD(s *dbschema.NgingSshUser, writer io.Writer, commands ...string) error {
	if len(commands) == 0 {
		return nil
	}
	client, err := Connect(s)
	if err != nil {
		return err
	}
	defer client.Close()
	err = sshclient.NewRemoteScript(
		client.Client,
		sshclient.RSScript(strings.Join(commands, "\r\n")),
	).SetStdio(writer, writer).Run()
	return err
}

func Connect(s *dbschema.NgingSshUser) (*webTerminalSSH.SSH, error) {
	c := sshconf.ToSFTPConfig(s)
	client, err := c.MakeClient()
	if err != nil {
		return nil, err
	}
	err = client.Connect()
	if err != nil {
		return nil, err
	}
	//defer client.Close()
	return client, nil
}
