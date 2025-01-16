package config

import "github.com/webx-top/echo"

var VhostModules = echo.NewKVData().
	Add(`basic`, echo.T(`通用设置`)).
	//Add(`log`, echo.T(`访问日志`), echo.KVOptHKV(`turnable`, true)).
	Add(`header`, echo.T(`响应Header`), echo.KVOptHKV(`turnable`, true), echo.KVOptHKV(`docURL`, `https://caddyserver.com/docs/header`)).
	// Add(`tls`, echo.T(`HTTPS`), echo.KVOptHKV(`turnable`, true), echo.KVOptHKV(`docURLs`, map[string]string{
	// 	`default`: `https://caddyserver.com/docs/tls`,
	// 	`caddy2`:  `https://caddyserver.com/docs/caddyfile/directives/tls`,
	// 	`nginx`:   `https://nginx.org/en/docs/http/ngx_http_ssl_module.html#ssl`,
	// })).
	// Add(`gzip`, echo.T(`GZIP`), echo.KVOptHKV(`turnable`, true), echo.KVOptHKV(`docURLs`, map[string]string{
	// 	`default`: `https://caddyserver.com/docs/gzip`,
	// 	`caddy2`:  `https://caddyserver.com/docs/caddyfile/directives/encode`,
	// 	`nginx`:   `https://nginx.org/en/docs/http/ngx_http_gzip_module.html#gzip`,
	// })).
	Add(`fastcgi`, echo.T(`FastCGI`), echo.KVOptHKV(`turnable`, true), echo.KVOptHKV(`docURLs`, map[string]string{
		`default`: `https://caddyserver.com/docs/fastcgi`,
		`caddy2`:  `https://caddyserver.com/docs/caddyfile/directives/php_fastcgi#syntax`,
		`nginx`:   `https://nginx.org/en/docs/http/ngx_http_fastcgi_module.html#fastcgi_pass`,
	})).
	Add(`proxy`, echo.T(`反向代理`), echo.KVOptHKV(`turnable`, true), echo.KVOptHKV(`docURLs`, map[string]string{
		`default`: `https://caddyserver.com/docs/proxy`,
		`caddy2`:  `https://caddyserver.com/docs/caddyfile/directives/reverse_proxy#syntax`,
		`nginx`:   `https://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_pass`,
	})).
	Add(`browse`, echo.T(`文件服务`), echo.KVOptHKV(`turnable`, true), echo.KVOptHKV(`docURLs`, map[string]string{
		`default`: `https://caddyserver.com/docs/browse`,
		`caddy2`:  `https://caddyserver.com/docs/caddyfile/directives/file_server#syntax`,
		`nginx`:   `https://nginx.org/en/docs/http/ngx_http_autoindex_module.html`,
	})).
	//Add(`expires`, echo.T(`静态文件缓存`), echo.KVOptHKV(`turnable`, true), echo.KVOptHKV(`docURL`, `https://caddyserver.com/docs/http.expires`)).
	Add(`ipfilter`, echo.T(`IP过滤`), echo.KVOptHKV(`turnable`, true), echo.KVOptHKV(`docURL`, `https://caddyserver.com/docs/http.ipfilter`)).
	Add(`filter`, echo.T(`响应内容过滤`), echo.KVOptHKV(`turnable`, true), echo.KVOptHKV(`noWrapper`, true), echo.KVOptHKV(`docURL`, `https://caddyserver.com/docs/http.filter`)).
	Add(`rewrite`, echo.T(`网址重写`), echo.KVOptHKV(`docURLs`, map[string]string{
		`default`: `https://caddyserver.com/docs/rewrite`,
		`caddy2`:  `https://caddyserver.com/docs/caddyfile/directives/rewrite`,
		`nginx`:   `https://nginx.org/en/docs/http/ngx_http_rewrite_module.html#rewrite`,
	})).
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
