package form

import (
	"github.com/caddyserver/certmagic"
	//"github.com/coscms/forms/config"
	"github.com/webx-top/echo"
)

var CAList = echo.NewKVData().
	Add(`LetsEncrypt`, `Let's Encrypt`, echo.KVOptX(certmagic.LetsEncryptProductionCA), echo.KVOptHKV(`Staging`, certmagic.LetsEncryptStagingCA), echo.KVOptHKV(`Issuer`, `acme`)).
	Add(`ZeroSSL`, `Zero SSL`, echo.KVOptX(certmagic.ZeroSSLProductionCA), echo.KVOptHKV(`NeedToken`, true), echo.KVOptHKV(`Issuer`, `zerossl`)).
	Add(`GoogleTrust`, `Google Trust`, echo.KVOptX(certmagic.GoogleTrustProductionCA), echo.KVOptHKV(`Staging`, certmagic.GoogleTrustStagingCA), echo.KVOptHKV(`Issuer`, `acme`), echo.KVOptHKV(`NeedEAB`, true))
