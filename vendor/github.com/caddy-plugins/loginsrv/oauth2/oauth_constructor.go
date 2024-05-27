package oauth2

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/bitbucket"
	"github.com/admpub/goth/providers/gitea"
	"github.com/admpub/goth/providers/github"
	"github.com/admpub/goth/providers/gitlab"
	"github.com/admpub/goth/providers/google"
	"github.com/admpub/goth/providers/paypal"
	"github.com/admpub/goth/providers/salesforce"
	"github.com/admpub/goth/providers/stripe"
	"github.com/admpub/goth/providers/wechat"
	"github.com/admpub/goth/providers/yahoo"
)

var constructors = map[string]func(*Config) goth.Provider{
	`gitea`: func(cfg *Config) goth.Provider {
		return gitea.New(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURI, cfg.GetScopes()...)
	},
	`google`: func(cfg *Config) goth.Provider {
		return google.New(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURI)
	},
	`github`: func(cfg *Config) goth.Provider {
		return github.New(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURI, cfg.GetScopes()...)
	},
	`bitbucket`: func(cfg *Config) goth.Provider {
		return bitbucket.New(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURI, cfg.GetScopes()...)
	},
	`gitlab`: func(cfg *Config) goth.Provider {
		return gitlab.New(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURI)
	},
	`salesforce`: func(cfg *Config) goth.Provider {
		return salesforce.New(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURI, cfg.GetScopes()...)
	},
	`yahoo`: func(cfg *Config) goth.Provider {
		return yahoo.New(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURI, cfg.GetScopes()...)
	},
	`stripe`: func(cfg *Config) goth.Provider {
		return stripe.New(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURI, cfg.GetScopes()...)
	},
	`paypal`: func(cfg *Config) goth.Provider {
		return paypal.New(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURI, cfg.GetScopes()...)
	},
	`wechat`: func(cfg *Config) goth.Provider {
		return wechat.New(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURI, wechat.WECHAT_LANG_CN)
	},
}

// ConstructorList returns the names of all registered constructor
func ConstructorList() []string {
	list := make([]string, 0, len(constructors))
	for k := range constructors {
		list = append(list, k)
	}
	return list
}

func Register(name string, newi func(cfg *Config) goth.Provider) {
	constructors[name] = newi
}

func GetConstructor(name string) func(cfg *Config) goth.Provider {
	return constructors[name]
}
