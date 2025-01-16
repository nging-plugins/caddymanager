package config

import "github.com/webx-top/echo"

var VhostModules = echo.NewKVData().
	Add(`basic`, echo.T(`基本设置`)).
	Add(`log`, echo.T(`访问日志`)).
	Add(`header`, echo.T(`响应Header`)).
	Add(`tls`, echo.T(`HTTPS`)).
	Add(`gzip`, echo.T(`GZIP`)).
	Add(`fastcgi`, echo.T(`FastCGI`)).
	Add(`proxy`, echo.T(`反向代理`)).
	Add(`browse`, echo.T(`文件服务`)).
	Add(`expires`, echo.T(`静态文件缓存`)).
	Add(`ipfilter`, echo.T(`IP过滤`)).
	Add(`filter`, echo.T(`响应内容过滤`)).
	Add(`rewrite`, echo.T(`网址重写`)).
	Add(`basicauth`, echo.T(`BasicAuth`)).
	Add(`ratelimit`, echo.T(`限流`)).
	Add(`cors`, echo.T(`跨域支持(CORS)`)).
	Add(`jwt`, echo.T(`JWT`)).
	Add(`login`, echo.T(`Login`)).
	Add(`locale`, echo.T(`多语言`)).
	Add(`nobots`, echo.T(`反爬虫`)).
	Add(`prometheus`, echo.T(`Prometheus`)).
	Add(`s3browser`, echo.T(`S3Browser`)).
	Add(`webdav`, echo.T(`WebDAV`))
