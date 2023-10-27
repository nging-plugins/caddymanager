package form

import (
	"html/template"
	"mime"
	"sort"
	"strconv"
	"strings"

	"github.com/webx-top/com"
)

type NginxDomainInfo struct {
	Port    int
	Args    []string
	Domains []string
}

/*
$remote_addr             客户端地址                                    211.28.65.253
$remote_user             客户端用户名称                                --
$time_local              访问时间和时区                                18/Jul/2012:17:00:01 +0800
$request                 请求的URI和HTTP协议                           "GET /article-10000.html HTTP/1.1"
$http_host               请求地址，即浏览器中你输入的地址（IP或域名）     www.wang.com 192.168.100.100
$status                  HTTP请求状态                                  200
$upstream_status         upstream状态                                  200
$body_bytes_sent         发送给客户端文件内容大小                        1547
$http_referer            url跳转来源                                   https://www.baidu.com/
$http_user_agent         用户终端浏览器等信息                           "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0; SV1; GTB7.0; .NET4.0C;
$ssl_protocol            SSL协议版本                                   TLSv1
$ssl_cipher              交换数据中的算法                               RC4-SHA
$upstream_addr           后台upstream的地址，即真正提供服务的主机地址     10.10.10.100:80
$request_time            整个请求的总时间                               0.205
$upstream_response_time  请求过程中，upstream响应时间                    0.002
*/
//{remote} - {user} [{when}] "{method} {scheme} {host} {uri} {proto}" {status} {size} "{>Referer}" "{>User-Agent}" {latency}
var NginxLogFormatReplacer = strings.NewReplacer(
	`{remote}`, `$remote_addr`,
	`{user}`, `$remote_user`,
	`{when}`, `$time_local`,
	//`{method}`,`$remote_addr`,
	//`{scheme}`,`$remote_addr`,
	`{host}`, `$http_host`,
	`{uri}`, `$request_uri`,
	`{method} {scheme} {host} {uri} {proto}`, `$request`,
	`{method} {uri} {proto}`, `$request`,
	//`{proto}`,`$remote_addr`,
	`{status}`, `$status`,
	`{size}`, `$body_bytes_sent`,
	`{>Referer}`, `$http_referer`,
	`{>User-Agent}`, `$http_user_agent`,
	`{latency}`, `$request_time`,
)

func AsNginxLogFormat(value string) string {
	value = ExplodeCombinedLogFormat(value)
	return NginxLogFormatReplacer.Replace(value)
}

func (v Values) AsNginxLogFormat() string {
	return AsNginxLogFormat(v.Get(`log_format`))
}

func (v Values) ExtensionsToMime(value string) []string {
	extensions := SplitBySpace(value)
	mimes := make([]string, 0, len(extensions))
	for _, ext := range extensions {
		mimeType := mime.TypeByExtension(ext)
		if len(mimeType) > 0 {
			mimes = append(mimes, mimeType)
		}
	}
	return mimes
}

func (v Values) GetNginxDomainList() []NginxDomainInfo {
	domainList := v.GetDomainList()
	var list []NginxDomainInfo
	portsDomains := map[int][]string{}
	var ports []int
	for _, domain := range domainList {
		domain = com.ParseEnvVar(domain)
		parts := strings.SplitN(domain, `://`, 2)
		var scheme, host, port string
		if len(parts) == 2 {
			scheme = parts[1]
			host, port = com.SplitHostPort(parts[1])
		} else {
			host, port = com.SplitHostPort(domain)
		}
		portNumber, _ := strconv.ParseUint(port, 10, 16)
		portN := int(portNumber)
		if portN == 0 {
			switch scheme {
			case `http`:
				portN = 80
			case `https`:
				portN = 443
			default:
				portN = 80
			}
		}
		if len(host) == 0 {
			host = `127.0.0.1`
		}
		if _, ok := portsDomains[portN]; !ok {
			portsDomains[portN] = []string{}
			ports = append(ports, portN)
		}
		portsDomains[portN] = append(portsDomains[portN], host)
	}
	sort.Sort(sort.IntSlice(ports))
	isTLS := v.Values.Get(`tls`) == `1`
	for _, portN := range ports {
		info := NginxDomainInfo{
			Port:    portN,
			Domains: portsDomains[portN],
		}
		if isTLS {
			info.Args = append(info.Args, `ssl`, `http2`)
		}
		list = append(list, info)
	}
	return list
}

type UpstreamInfo struct {
	Scheme       string
	Host         string
	Path         string
	UpstreamName string
	withQuote    bool
}

func (u UpstreamInfo) String() string {
	host := u.Host
	if len(u.UpstreamName) > 0 {
		host = u.UpstreamName
	}
	value := u.Scheme + `://` + host + u.Path
	if u.withQuote {
		value = `"` + com.AddCSlashes(value, '"') + `"`
	}
	return value
}

func (v Values) ServerGroup(key string, customHost string, withQuotes ...bool) interface{} {
	var withQuote bool
	if len(withQuotes) > 0 {
		withQuote = withQuotes[0]
	}
	val := v.Get(key)
	sh := strings.SplitN(val, `://`, 2)
	var scheme string
	var host string
	var path string
	if len(sh) == 2 {
		scheme = sh[0]
		host = sh[1]
	} else {
		scheme = `http`
		host = val
	}
	hp := strings.SplitN(host, `/`, 2)
	if len(hp) == 2 {
		host = hp[0]
		path = `/` + hp[1]
	}
	return UpstreamInfo{
		Scheme:       scheme,
		Host:         host,
		Path:         path,
		UpstreamName: customHost,
		withQuote:    withQuote,
	}
}

func (v Values) IteratorHeaderKV(addon string, item string, plusPrefix string, minusPrefix string, withValueAndQuotes ...bool) interface{} {
	if len(addon) > 0 && len(item) > 0 {
		addon += `_`
	}
	k := addon + item + `_k`
	keys, _ := v.Values[k]

	k = addon + item + `_v`
	values, _ := v.Values[k]

	var r, t string
	var withValueAndQuote bool
	if len(withValueAndQuotes) > 0 {
		withValueAndQuote = withValueAndQuotes[0]
	}
	l := len(values)
	var suffix string
	if v.Config.GetEngine() == `nginx` {
		suffix = `;`
	}
	for i, k := range keys {
		if i < l {
			prefix := plusPrefix
			if strings.HasPrefix(k, `-`) {
				if len(minusPrefix) == 0 {
					continue
				}
				k = strings.TrimPrefix(k, `-`)
				prefix = minusPrefix
			} else {
				if len(plusPrefix) == 0 {
					continue
				}
				k = strings.TrimPrefix(k, `+`)
			}
			if withValueAndQuote {
				v := values[i]
				v = `"` + com.AddCSlashes(v, '"') + `"`
				r += t + prefix + k + `   ` + v + suffix
			} else {
				r += t + prefix + k + suffix
			}
			t = "\n"
		}
	}
	if withValueAndQuote {
		return template.HTML(r)
	}
	return r
}
