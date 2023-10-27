package engine

import (
	"context"
	"os"
)

type Configer interface {
	GetVhostConfigLocalDirAbs() (string, error)
	GetTemplateFile() string
	GetIdent() string
	GetEngine() string
	GetEnviron() string
	GetCertLocalDir() string
	GetCertContainerDir() string
	GetVhostConfigLocalDir() string
	GetVhostConfigContainerDir() string
	GetEngineConfigLocalFile() string
	GetEngineConfigContainerFile() string
}

type CertRenewaler interface {
	RenewalCert(ctx context.Context, domain, email string) error
}

type EngineConfigFileFixer interface {
	FixEngineConfigFile(deleteMode ...bool) (bool, error)
}

func FixEngineConfigFile(cfg Configer, deleteMode ...bool) (bool, error) {
	if fx, ok := cfg.(EngineConfigFileFixer); ok {
		hasUpdate, err := fx.FixEngineConfigFile(deleteMode...)
		if err != nil && os.IsNotExist(err) {
			return hasUpdate, nil
		}
		return hasUpdate, err
	}
	return false, nil
}
