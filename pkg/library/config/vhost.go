package config

import (
	"net/url"

	"github.com/nging-plugins/caddymanager/pkg/dbschema"
)

type VhostConfig struct {
	*dbschema.NgingVhost
	FormData url.Values
}
