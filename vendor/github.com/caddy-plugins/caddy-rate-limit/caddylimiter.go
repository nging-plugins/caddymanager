package ratelimit

import (
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type CaddyLimiter struct {
	keys map[string]*rate.Limiter
	sync.RWMutex
}

func NewCaddyLimiter() *CaddyLimiter {

	return &CaddyLimiter{
		keys: make(map[string]*rate.Limiter),
	}
}

// Allow is just a shortcut for AllowN
func (cl *CaddyLimiter) Allow(keys []string, rule Rule, limiter ...*rate.Limiter) bool {

	return cl.AllowN(keys, rule, 1, limiter...)
}

// AllowN check if n count are allowed for a specific key
func (cl *CaddyLimiter) AllowN(keys []string, rule Rule, n int, limiter ...*rate.Limiter) bool {
	if len(limiter) > 0 && limiter[0] != nil {
		return limiter[0].AllowN(time.Now(), n)
	}
	keysJoined := strings.Join(keys, "|")
	curLimiter, found := cl.GetLimiterOk(keysJoined)
	if !found {
		switch rule.Unit {
		case "second":
			curLimiter = rate.NewLimiter(rate.Limit(rule.Rate)/rate.Limit(time.Second.Seconds()), rule.Burst)
		case "minute":
			curLimiter = rate.NewLimiter(rate.Limit(rule.Rate)/rate.Limit(time.Minute.Seconds()), rule.Burst)
		case "hour":
			curLimiter = rate.NewLimiter(rate.Limit(rule.Rate)/rate.Limit(time.Hour.Seconds()), rule.Burst)
		case "day":
			curLimiter = rate.NewLimiter(rate.Limit(rule.Rate)/rate.Limit(24*time.Hour.Seconds()), rule.Burst)
		case "week":
			curLimiter = rate.NewLimiter(rate.Limit(rule.Rate)/rate.Limit(7*24*time.Hour.Seconds()), rule.Burst)
		default:
			// Infinite
			curLimiter = rate.NewLimiter(rate.Inf, rule.Burst)
		}
		cl.SetLimiter(keysJoined, curLimiter)
	}

	return curLimiter.AllowN(time.Now(), n)
}

func (cl *CaddyLimiter) GetLimiter(key string) *rate.Limiter {
	cl.RLock()
	limiter := cl.keys[key]
	cl.RUnlock()
	return limiter
}

func (cl *CaddyLimiter) SetLimiter(key string, limiter *rate.Limiter) {
	cl.Lock()
	cl.keys[key] = limiter
	cl.Unlock()
}

func (cl *CaddyLimiter) GetLimiterOk(key string) (*rate.Limiter, bool) {
	cl.RLock()
	limiter, ok := cl.keys[key]
	cl.RUnlock()
	return limiter, ok
}

func (cl *CaddyLimiter) HasLimiter(key string) bool {
	cl.RLock()
	_, ok := cl.keys[key]
	cl.RUnlock()
	return ok
}

// RetryAfter return a helper message for client
func (cl *CaddyLimiter) RetryAfter(keys []string, limiter ...*rate.Limiter) time.Duration {
	var curLimiter *rate.Limiter
	if len(limiter) > 0 && limiter[0] != nil {
		curLimiter = limiter[0]
	} else {
		keysJoined := strings.Join(keys, "|")
		curLimiter = cl.GetLimiter(keysJoined)
	}
	reserve := curLimiter.Reserve()
	defer reserve.Cancel()

	if reserve.OK() {
		retryAfter := reserve.Delay()
		return retryAfter
	}

	return rate.InfDuration
}

// Reserve will consume 1 token from `token bucket`
func (cl *CaddyLimiter) Reserve(keys []string) bool {
	keysJoined := strings.Join(keys, "|")
	r := cl.GetLimiter(keysJoined).Reserve()
	return r.OK()
}

// buildKeys combine client-ip / request-header, methods, status code and resource
func buildKeys(limitedKey, methods, status, res string) [][]string {

	sliceKeys := make([][]string, 0)

	if len(limitedKey) != 0 {
		sliceKeys = append(sliceKeys, []string{limitedKey, methods, status, res})
	}

	return sliceKeys
}

// buildKeysOnlyWithKey will use client ip or request header as key
func buildKeysOnlyWithLimitedKey(limitedKey string) [][]string {
	sliceKeys := make([][]string, 0)

	if len(limitedKey) != 0 {
		sliceKeys = append(sliceKeys, []string{limitedKey})
	}

	return sliceKeys
}
