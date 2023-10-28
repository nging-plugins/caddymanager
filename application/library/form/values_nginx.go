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

// {remote} - {user} [{when}] "{method} {scheme} {host} {uri} {proto}" {status} {size} "{>Referer}" "{>User-Agent}" {latency}
var NginxLogFormatReplacer = strings.NewReplacer(
	//`{method} {scheme} {host} {uri} {proto}`, `$request`,
	`{method} {uri} {proto}`, `$request`,
	`{remote}`, `$remote_addr`,
	`{user}`, `$remote_user`,
	`{when}`, `$time_local`,
	`{method}`, `$request_method`,
	`{scheme}`, `$scheme`,
	`{host}`, `$http_host`,
	`{uri}`, `$request_uri`,
	`{proto}`, `$server_protocol`,
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
		// text/css; charset=utf-8
		mimeType := strings.SplitN(mime.TypeByExtension(ext), `;`, 2)[0]
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

func (v Values) IteratorNginxProxyHeaderKV() interface{} {
	addon := `proxy`
	item := `header_downstream`
	if len(addon) > 0 && len(item) > 0 {
		addon += `_`
	}
	k := addon + item + `_k`
	keys, _ := v.Values[k]

	k = addon + item + `_v`
	values, _ := v.Values[k]

	var r, t string
	l := len(values)
	suffix := `;`
	for i, k := range keys {
		if i < l {
			var prefix string
			if strings.HasPrefix(k, `-`) {
				k = strings.TrimPrefix(k, `-`)
				prefix = `proxy_hide_header`
				r += t + prefix + ` ` + k + suffix
			} else {
				k = strings.TrimPrefix(k, `+`)
				prefix = `proxy_pass_header`
				v := values[i]
				if len(v) > 0 {
					prefix = `add_header`
					v = `"` + com.AddCSlashes(v, '"') + `"`
					r += t + prefix + ` ` + k + ` ` + v + suffix
				} else {
					r += t + prefix + ` ` + k + suffix
				}
			}
			t = "\n"
		}
	}
	return template.HTML(r)
}
