package form

import (
	"github.com/caddyserver/certmagic"
	//"github.com/coscms/forms/config"
	"github.com/webx-top/echo"
)

/*
// ZeroSSL 表单配置
var zeroSSLFormConfig = &config.Config{
	Theme:    `bootstrap3`,
	Template: `allfields`,
	Elements: []*config.Element{
		{
			Type:  `text`,
			Name:  `tls_zerossl_apikey`,
			Label: echo.T(`API Key`),
			Attributes: [][]string{
				{`required`, `true`},
			},
			HelpText: echo.T(`ZeroSSL API Key`),
		},
	},
}
*/

var CAList = echo.NewKVData().
	Add(`LetsEncrypt`, `Let's Encrypt`, echo.KVOptX(certmagic.LetsEncryptProductionCA), echo.KVOptHKV(`Staging`, certmagic.LetsEncryptStagingCA)).
	Add(`ZeroSSL`, `Zero SSL`, echo.KVOptX(certmagic.ZeroSSLProductionCA) /*echo.KVOptHKV(`FormConfig`, zeroSSLFormConfig)*/).
	Add(`GoogleTrust`, `Google Trust`, echo.KVOptX(certmagic.GoogleTrustProductionCA), echo.KVOptHKV(`Staging`, certmagic.GoogleTrustStagingCA))
