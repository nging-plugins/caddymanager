package config

import (
	"golang.org/x/crypto/ssh"
)

func NewHostConfig(cfg *ssh.ClientConfig, host string, port int) *HostConfig {
	return &HostConfig{ClientConfig: cfg, Host: host, Port: port}
}

type HostConfig struct {
	*ssh.ClientConfig
	Host    string
	Port    int
	Account *AccountConfig
}

func (c *HostConfig) SetAccount(account *AccountConfig) *HostConfig {
	c.Account = account
	return c
}

// NewSSHConfig
func NewSSHConfig(end ...*HostConfig) *SSHConfig {
	c := &SSHConfig{
		Transform: NewTransformConfig(),
	}
	l := len(end)
	if l > 0 {
		j := l - 1
		if j > 0 {
			c.Jumps = append(c.Jumps, end[0:j]...)
		}
		c.End = end[j]
	}
	return c
}

type SSHConfig struct {
	End       *HostConfig
	Jumps     []*HostConfig
	Transform *TransformConfig
}

func (c *SSHConfig) SetEnd(endHostConfig *HostConfig) *SSHConfig {
	c.End = endHostConfig
	return c
}

func (c *SSHConfig) AddJump(jumpHostConfigs ...*HostConfig) *SSHConfig {
	c.Jumps = append(c.Jumps, jumpHostConfigs...)
	return c
}
