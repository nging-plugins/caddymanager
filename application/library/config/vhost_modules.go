package config

import "github.com/webx-top/echo"

var VhostModules = echo.NewKVData().
	Add(`basic`, echo.T(`基本设置`)).
	Add(`log`, echo.T(`访问日志`), echo.KVOptHKV(`turnable`, true)).
	Add(`header`, echo.T(`响应Header`)).
	Add(`tls`, echo.T(`HTTPS`), echo.KVOptHKV(`turnable`, true)).
	Add(`gzip`, echo.T(`GZIP`), echo.KVOptHKV(`turnable`, true)).
	Add(`fastcgi`, echo.T(`FastCGI`), echo.KVOptHKV(`turnable`, true)).
	Add(`proxy`, echo.T(`反向代理`), echo.KVOptHKV(`turnable`, true)).
	Add(`browse`, echo.T(`文件服务`), echo.KVOptHKV(`turnable`, true)).
	Add(`expires`, echo.T(`静态文件缓存`), echo.KVOptHKV(`turnable`, true)).
	Add(`ipfilter`, echo.T(`IP过滤`), echo.KVOptHKV(`turnable`, true)).
	Add(`filter`, echo.T(`响应内容过滤`), echo.KVOptHKV(`turnable`, true)).
	Add(`rewrite`, echo.T(`网址重写`)).
	Add(`basicauth`, echo.T(`BasicAuth`), echo.KVOptHKV(`turnable`, true)).
	Add(`ratelimit`, echo.T(`限流`), echo.KVOptHKV(`turnable`, true)).
	Add(`cors`, echo.T(`跨域支持(CORS)`), echo.KVOptHKV(`turnable`, true)).
	Add(`jwt`, echo.T(`JWT`), echo.KVOptHKV(`turnable`, true)).
	Add(`login`, echo.T(`Login`), echo.KVOptHKV(`turnable`, true)).
	Add(`locale`, echo.T(`多语言`), echo.KVOptHKV(`turnable`, true)).
	Add(`nobots`, echo.T(`反爬虫`), echo.KVOptHKV(`turnable`, true)).
	Add(`prometheus`, echo.T(`Prometheus`), echo.KVOptHKV(`turnable`, true)).
	Add(`s3browser`, echo.T(`S3Browser`), echo.KVOptHKV(`turnable`, true)).
	Add(`webdav`, echo.T(`WebDAV`), echo.KVOptHKV(`turnable`, true))
