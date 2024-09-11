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
	"os"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/coscms/webcore/cmd"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/nging-plugins/caddymanager/application/model"
	"github.com/spf13/cobra"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

// Lego 钩子

var legoHookCmd = &cobra.Command{
	Use:   "legohook",
	Short: "lego hook",
	RunE:  legoHookRunE,
}

type EnvVarFromLego struct {
	Email    string
	Domain   string
	PathCert string
	PathKey  string
}

func legoHookRunE(cmd *cobra.Command, args []string) error {
	/*
	   Some information is provided through environment variables:
	       LEGO_ACCOUNT_EMAIL: the email of the account.
	       LEGO_CERT_DOMAIN: the main domain of the certificate.
	       LEGO_CERT_PATH: the path of the certificate.
	       LEGO_CERT_KEY_PATH: the path of the certificate key.
	*/
	data := &EnvVarFromLego{
		Email:    os.Getenv(`LEGO_ACCOUNT_EMAIL`),
		Domain:   os.Getenv(`LEGO_CERT_DOMAIN`),
		PathCert: os.Getenv(`LEGO_CERT_PATH`),
		PathKey:  os.Getenv(`LEGO_CERT_KEY_PATH`),
	}
	log.Okayf("LEGO_ACCOUNT_EMAIL: %v", data.Email)
	log.Okayf("LEGO_CERT_DOMAIN: %v", data.Domain)
	log.Okayf("LEGO_CERT_PATH: %v", data.PathCert)
	log.Okayf("LEGO_CERT_KEY_PATH: %v", data.PathKey)
	if len(data.Domain) == 0 {
		log.Error("LEGO_CERT_DOMAIN is empty")
		return nil
	}
	conf, err := config.InitConfig()
	config.MustOK(err)
	conf.AsDefault()
	ctx := defaults.NewMockContext()
	m := model.NewVhost(ctx)
	cond := db.NewCompounds()
	cond.AddKV(`disabled`, common.BoolN)
	domain := data.Domain
	domainReverserReplacer := strings.NewReplacer(
		"-", ":",
		"_", "*",
	)
	domain = domainReverserReplacer.Replace(domain)
	cond.Add(db.Or(
		db.Cond{`domain`: domain},
		db.Cond{`domain`: db.Like(domain + ` %`)},
		db.Cond{`domain`: db.Like(domain + `:%`)},
		db.Cond{`domain`: db.Like(`% ` + domain)},
		db.Cond{`domain`: db.Like(`% ` + domain + ` %`)},
		db.Cond{`domain`: db.Like(`% ` + domain + `:%`)},
	))
	err = m.Get(nil, cond.And())
	if err != nil {
		return err
	}
	set := echo.H{}
	if m.SslObtained == 0 {
		m.SslObtained = uint(time.Now().Unix())
		set.Set(`ssl_obtained`, m.SslObtained)
	} else {
		m.SslRenewed = uint(time.Now().Unix())
		set.Set(`ssl_renewed`, m.SslRenewed)
	}
	return m.UpdateFields(nil, set, `id`, m.Id)
}

func init() {
	cmd.Add(legoHookCmd)
}
