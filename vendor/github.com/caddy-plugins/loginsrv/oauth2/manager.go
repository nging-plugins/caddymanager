package oauth2

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/admpub/goth"
	"github.com/caddy-plugins/loginsrv/model"
)

var (
	errMissingParameterClientID     = errors.New("missing parameter client_id")
	errMissingParameterClientSecret = errors.New("missing parameter client_secret")
)

// Manager has the responsibility to handle the user user requests in an oauth flow.
// It has to pick the right configuration and start the oauth redirecting.
type Manager struct {
	configs      map[string]*Config
	startFlow    func(cfg *Config, w http.ResponseWriter, r *http.Request) error
	authenticate func(cfg *Config, w http.ResponseWriter, r *http.Request) (model.UserInfo, error)
}

// NewManager creates a new Manager
func NewManager() *Manager {
	return &Manager{
		configs:      map[string]*Config{},
		startFlow:    StartFlow,
		authenticate: Authenticate,
	}
}

// Handle is managing the oauth flow.
// Dependent on the code parameter of the url, the oauth flow is started or
// the call is interpreted as the redirect callback and the token exchange is done.
// Return parameters:
//
//	startedFlow - true, if this was the initial call to start the oauth flow
//	authenticated - if the authentication was successful or not
//	userInfo - the user info from the provider in case of a successful authentication
//	err - an error
func (manager *Manager) Handle(w http.ResponseWriter, r *http.Request) (
	startedFlow bool,
	authenticated bool,
	userInfo model.UserInfo,
	err error) {

	if r.FormValue("error") != "" {
		return false, false, model.UserInfo{}, fmt.Errorf("error: %v", r.FormValue("error"))
	}

	cfg, err := manager.GetConfigFromRequest(r)
	if err != nil {
		return false, false, model.UserInfo{}, err
	}

	if r.FormValue("code") != "" {
		userInfo, err = manager.authenticate(cfg, w, r)
		if err != nil {
			return false, false, model.UserInfo{}, err
		}
		return false, true, userInfo, err
	}

	manager.startFlow(cfg, w, r)
	return true, false, model.UserInfo{}, nil
}

// GetConfigFromRequest returns the oauth configuration matching the current path.
// The configuration name is taken from the last path segment.
func (manager *Manager) GetConfigFromRequest(r *http.Request) (*Config, error) {
	configName := manager.getConfigNameFromPath(r.URL.Path)
	cfg, exist := manager.configs[configName]
	if !exist {
		return nil, fmt.Errorf("no oauth configuration for %v", configName)
	}

	if cfg.RedirectURI == "" {
		cfg.RedirectURI = redirectURIFromRequest(r)
	}

	if cfg.Provider == nil {
		err := cfg.InitProvider()
		if err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

func (manager *Manager) getConfigNameFromPath(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

// AddConfig for a provider
func (manager *Manager) AddConfig(providerName string, opts map[string]string) error {
	cfg := &Config{
		ProviderName: providerName,
	}
	cfg.Extra = map[string]string{}
	for k, v := range opts {
		switch k {
		case `client_id`:
			cfg.ClientID = v
		case `client_secret`:
			cfg.ClientSecret = v
		case `scope`:
			cfg.Scope = v
		case `redirect_uri`:
			cfg.RedirectURI = v
		default:
			cfg.Extra[k] = v
		}
	}

	if len(cfg.ClientID) == 0 {
		return errMissingParameterClientID
	}
	if len(cfg.ClientSecret) == 0 {
		return errMissingParameterClientSecret
	}

	p, err := goth.GetProvider(providerName)
	if err != nil {
		p, err = cfg.NewProvider()
		if err != nil {
			return err
		}
	}

	cfg.Provider = p
	manager.configs[providerName] = cfg
	return nil
}

// GetConfigs of the manager
func (manager *Manager) GetConfigs() map[string]*Config {
	return manager.configs
}

func redirectURIFromRequest(r *http.Request) string {
	u := url.URL{}
	u.Path = r.URL.Path

	if ffh := r.Header.Get("X-Forwarded-Host"); ffh == "" {
		u.Host = r.Host
	} else {
		u.Host = ffh
	}

	if ffp := r.Header.Get("X-Forwarded-Proto"); ffp == "" {
		if r.TLS != nil {
			u.Scheme = "https"
		} else {
			u.Scheme = "http"
		}
	} else {
		u.Scheme = ffp
	}

	return u.String()
}
