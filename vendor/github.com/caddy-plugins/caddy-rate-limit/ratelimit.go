package ratelimit

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/caddy/caddyhttp/httpserver"
	"github.com/admpub/realip"
)

// RateLimit is an http.Handler that can limit request rate to specific paths or files
type RateLimit struct {
	Next  httpserver.Handler
	Rules []Rule
}

const HTTPTimeLayout = "Mon, 02 Jan 2006 15:04:05 GMT"

// Rule is a configuration for ratelimit
type Rule struct {
	Methods       string
	Rate          int64
	Burst         int
	Unit          string
	Whitelist     []string
	LimitByHeader string
	Status        string
	Resources     []string

	whitelistIPNets []*net.IPNet
}

const (
	ignoreSymbol = "^"
)

var (
	caddyLimiter *CaddyLimiter
)

func init() {
	caddyLimiter = NewCaddyLimiter()
}

func SetRetryAfterHeader(header http.Header, retryAfter time.Duration) {
	header.Add("X-Retry-After", retryAfter.String())
	header.Add("Retry-After", time.Now().UTC().Add(retryAfter).Format(HTTPTimeLayout))
}

func ParseHTTPTime(value string) (time.Time, error) {
	return time.Parse(HTTPTimeLayout, value)
}

func getLimitedKeyByHeader(r *http.Request, rule Rule) (string, bool) {
	switch rule.LimitByHeader {
	case realip.HeaderForwarded:
		return realip.Default().ClientIP(r.RemoteAddr, r.Header.Get), true
	case realip.HeaderXForwardedFor:
		return realip.XRealIP(r.Header.Get(realip.HeaderXRealIP), r.Header.Get(realip.HeaderXForwardedFor), r.RemoteAddr), true
	case realip.HeaderXRealIP:
		return realip.XRealIP(r.Header.Get(realip.HeaderXRealIP), ``, r.RemoteAddr), true
	default:
		return r.Header.Get(rule.LimitByHeader), false
	}
}

type limitedItem struct {
	Key  string
	IsIP bool
}

// ServeHTTP is the method handling every request
func (rl RateLimit) ServeHTTP(w http.ResponseWriter, r *http.Request) (nextResponseStatus int, err error) {

	retryAfter := time.Duration(0)
	// get request ip address
	ipAddress, err := GetRemoteIP(r)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	ruleLimitedKeys := make([]limitedItem, len(rl.Rules))
	parsedHeaders := map[string]limitedItem{}

	for index, rule := range rl.Rules {
		limitedKey := ipAddress
		if len(rule.LimitByHeader) == 0 {
			ruleLimitedKeys[index] = limitedItem{
				Key:  ipAddress,
				IsIP: true,
			}
		} else {
			item, ok := parsedHeaders[rule.LimitByHeader]
			if !ok {
				var isIP bool
				limitedKey, isIP = getLimitedKeyByHeader(r, rule)
				item = limitedItem{Key: limitedKey, IsIP: isIP}
				parsedHeaders[rule.LimitByHeader] = item
			}
			ruleLimitedKeys[index] = item
		}
		for _, res := range rule.Resources {

			// handle `ignore`
			if strings.HasPrefix(res, ignoreSymbol) {
				res = strings.TrimPrefix(res, ignoreSymbol)
				if httpserver.Path(r.URL.Path).Matches(res) {
					return rl.Next.ServeHTTP(w, r)
				}
				continue
			}

			// handle path mismatch
			if !httpserver.Path(r.URL.Path).Matches(res) {
				continue
			}

			// handle whitelist ip & method mismatch
			if IsWhitelistIPAddress(ipAddress, rule.whitelistIPNets) || !MatchMethod(rule.Methods, r.Method) {
				continue
			}

			/*
				check if this ip has already exceeded quota
				if so, reject all the subsequent requests

				note: this won't block resources outside of the plugin's config
			*/
			sliceKeysOnlyWithKey := buildKeysOnlyWithLimitedKey(limitedKey)
			for _, keys := range sliceKeysOnlyWithKey {
				keysJoined := strings.Join(keys, "|")
				if limiter, ok := caddyLimiter.GetLimiterOk(keysJoined); ok {
					ret := caddyLimiter.Allow(keys, rule, limiter)
					if !ret {
						retryAfter = caddyLimiter.RetryAfter(keys, limiter)
						SetRetryAfterHeader(w.Header(), retryAfter)
						return http.StatusTooManyRequests, err
					}
				}
			}

			// check limit
			if len(rule.Status) == 0 || rule.Status == "*" {
				sliceKeys := buildKeys(limitedKey, rule.Methods, rule.Status, res)
				for _, keys := range sliceKeys {
					ret := caddyLimiter.Allow(keys, rule)
					if !ret {
						retryAfter = caddyLimiter.RetryAfter(keys)
						SetRetryAfterHeader(w.Header(), retryAfter)
						return http.StatusTooManyRequests, err
					}
				}
			}
		}
	}

	/*
		special case for limiting by response status code
	*/
	nextResponseStatus, err = rl.Next.ServeHTTP(w, r)

	for index, rule := range rl.Rules {
		// handle response status code mismatch
		if len(rule.Status) == 0 || rule.Status == "*" || !MatchStatus(rule.Status, strconv.Itoa(nextResponseStatus)) {
			continue
		}
		ipAddr := ipAddress
		if ruleLimitedKeys[index].IsIP {
			ipAddr = ruleLimitedKeys[index].Key
		}
		limitedKey := ruleLimitedKeys[index].Key
		for _, res := range rule.Resources {

			// handle `ignore`
			if strings.HasPrefix(res, ignoreSymbol) {
				res = strings.TrimPrefix(res, ignoreSymbol)
				if httpserver.Path(r.URL.Path).Matches(res) {
					return nextResponseStatus, err
				}
				continue
			}

			// handle path mismatch
			if !httpserver.Path(r.URL.Path).Matches(res) {
				continue
			}

			// handle whitelist ip & method mismatch
			if IsWhitelistIPAddress(ipAddr, rule.whitelistIPNets) || !MatchMethod(rule.Methods, r.Method) {
				continue
			}

			sliceKeys := buildKeysOnlyWithLimitedKey(limitedKey)
			for _, keys := range sliceKeys {
				// consume one token if status code matches
				caddyLimiter.Allow(keys, rule)
			}
		}
	}

	return nextResponseStatus, err
}
