package oauth2

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/admpub/goth"
	"github.com/admpub/goth/gothic"
	"github.com/caddy-plugins/loginsrv/model"
)

// Config describes a typical 3-legged OAuth2 flow, with both the
// client application information and the server's endpoint URLs.
type Config struct {
	// ClientID is the application's ID.
	ClientID string

	// ClientSecret is the application's secret.
	ClientSecret string

	// RedirectURL is the URL to redirect users going through
	// the OAuth flow, after the resource owner's URLs.
	RedirectURI string

	// Scope specifies optional requested permissions, this is a *space* separated list.
	Scope string

	Extra map[string]string

	// The oauth provider
	Provider     goth.Provider
	mu           sync.RWMutex
	ProviderName string
}

func (cfg *Config) SetProvider(provider goth.Provider) *Config {
	cfg.mu.Lock()
	cfg.Provider = provider
	cfg.mu.Unlock()
	return cfg
}

func (cfg *Config) GetProvider() goth.Provider {
	cfg.mu.RLock()
	provider := cfg.Provider
	cfg.mu.RUnlock()
	return provider
}

func (cfg *Config) InitProvider() error {
	provider, err := cfg.NewProvider()
	if err != nil {
		return err
	}
	cfg.SetProvider(provider)
	return nil
}

func (cfg *Config) GetScopes() []string {
	var scopes []string
	if len(cfg.Scope) > 0 {
		scopes = strings.Split(cfg.Scope, ` `)
	}
	return scopes
}

func (cfg *Config) NewProvider() (p goth.Provider, err error) {
	if newI, ok := constructors[cfg.ProviderName]; ok {
		p = newI(cfg)
		goth.UseProviders(p)
		return
	}
	err = fmt.Errorf("no provider for name %v", cfg.ProviderName)
	return
}

// StartFlow by redirecting the user to the login provider.
// A state parameter to protect against cross-site request forgery attacks is randomly generated and stored in a cookie
func StartFlow(cfg *Config, w http.ResponseWriter, r *http.Request) error {
	r = r.WithContext(context.WithValue(r.Context(), gothic.ProviderParamKey, cfg.ProviderName))
	gothic.BeginAuthHandler(w, r)
	return nil
}

// Authenticate after coming back from the oauth flow.
// Verify the state parameter againt the state cookie from the request.
func Authenticate(cfg *Config, w http.ResponseWriter, r *http.Request) (model.UserInfo, error) {
	r = r.WithContext(context.WithValue(r.Context(), gothic.ProviderParamKey, cfg.ProviderName))
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		return model.UserInfo{}, err
	}
	return model.UserInfo{
		Sub:     user.Name,
		Picture: user.AvatarURL,
		Name:    user.NickName,
		Email:   user.Email,
		Origin:  user.Provider,
	}, nil
}

// ProviderList returns the names of all registered provider
func ProviderList() []string {
	providers := goth.GetProviders()
	list := make([]string, 0, len(providers))
	for k := range providers {
		list = append(list, k)
	}
	return list
}
