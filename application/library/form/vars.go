package form

import (
	"github.com/caddyserver/certmagic"
	"github.com/webx-top/echo"
)

var CAList = echo.NewKVData().
	Add(`LetsEncrypt`, `Let's Encrypt`, echo.KVOptX(certmagic.LetsEncryptProductionCA), echo.KVOptHKV(`Staging`, certmagic.LetsEncryptStagingCA)).
	Add(`ZeroSSL`, `Zero SSL`, echo.KVOptX(certmagic.ZeroSSLProductionCA)).
	Add(`GoogleTrust`, `Google Trust`, echo.KVOptX(certmagic.GoogleTrustProductionCA), echo.KVOptHKV(`Staging`, certmagic.GoogleTrustStagingCA))
