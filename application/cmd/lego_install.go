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

package cmd

import (
	"github.com/coscms/webcore/cmd"
	"github.com/spf13/cobra"

	"github.com/nging-plugins/caddymanager/application/library/engine"
)

// Lego 安装

var legoInstallCmd = &cobra.Command{
	Use:   "legoinstall",
	Short: "lego install",
	RunE:  legoInstallRunE,
}

func legoInstallRunE(cmd *cobra.Command, args []string) error {
	return engine.LegoInstall()
}

func init() {
	cmd.Add(legoInstallCmd)
}
