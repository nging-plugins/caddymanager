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

package caddy

import (
	"fmt"
	"strconv"
	"strings"

	_ "github.com/caddy-plugins/caddy-expires"
	_ "github.com/caddy-plugins/caddy-filter"
	_ "github.com/caddy-plugins/caddy-jwt/v3"
	_ "github.com/caddy-plugins/caddy-locale"
	_ "github.com/caddy-plugins/caddy-prometheus"
	_ "github.com/caddy-plugins/caddy-rate-limit"
	s3browser "github.com/caddy-plugins/caddy-s3browser"
	_ "github.com/caddy-plugins/cors/caddy"
	_ "github.com/caddy-plugins/ipfilter"
	_ "github.com/caddy-plugins/loginsrv/caddy"
	_ "github.com/caddy-plugins/nobots"
	_ "github.com/caddy-plugins/webdav"

	//_ "github.com/caddy-plugins/caddy-iplimit/iplimit"

	"github.com/admpub/goth"
	"github.com/caddy-plugins/loginsrv/oauth2"
	_ "github.com/caddy-plugins/loginsrv/oauth2/register/gitea"
	_ "github.com/caddy-plugins/loginsrv/oauth2/register/github"
	_ "github.com/caddy-plugins/loginsrv/oauth2/register/gitlab"
	_ "github.com/caddy-plugins/loginsrv/oauth2/register/google"
	_ "github.com/caddy-plugins/loginsrv/oauth2/register/wechat"
	"github.com/coscms/webcore/library/backend/oauth2nging"
	"github.com/coscms/webcore/model"
	"github.com/webx-top/echo/defaults"

	// TLS DNS providers
	_ "github.com/admpub/caddy/dnsproviders/acmedns"
	//_ "github.com/admpub/caddy/dnsproviders/alidns"
	_ "github.com/admpub/caddy/dnsproviders/cloudflare"
	//_ "github.com/admpub/caddy/dnsproviders/dnspod"
	_ "github.com/admpub/caddy/dnsproviders/edgeone"
	_ "github.com/admpub/caddy/dnsproviders/he"
	_ "github.com/admpub/caddy/dnsproviders/rfc2136"
	_ "github.com/admpub/caddy/dnsproviders/tencentcloud"
)

func init() {
	oauth2.Register(`nging`, func(cfg *oauth2.Config) goth.Provider {
		hostURL := cfg.Extra[`host_url`]
		if len(hostURL) > 0 {
			hostURL = strings.TrimSuffix(hostURL, `/`)
		}
		return oauth2nging.New(cfg.ClientID, cfg.ClientSecret, cfg.GetRedirectURI(), hostURL, `profile`)
	})
	s3browser.AccountGetter = getS3Account
}

func getS3Account(arg string) (s3browser.Account, error) {
	id, err := strconv.ParseUint(arg, 10, 0)
	if err != nil {
		return s3browser.Account{}, fmt.Errorf(`[s3browser]failed to parse uint %q: %w`, arg, err)
	}
	ctx := defaults.NewMockContext()
	m := model.NewCloudStorage(ctx)
	err = m.Get(nil, `id`, id)
	if err != nil {
		return s3browser.Account{}, fmt.Errorf(`[s3browser]failed to get account by id %d: %w`, id, err)
	}
	return s3browser.Account{
		Key:      m.Key,
		Bucket:   m.Bucket,
		Secret:   m.Secret,
		Region:   m.Region,
		CDNURL:   m.Baseurl,
		Endpoint: m.Endpoint,
		Secure:   m.Secure == `Y`,
	}, nil
}
