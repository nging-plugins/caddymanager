package engine

import "context"

type Configer interface {
	GetVhostConfigDirAbsPath() (string, error)
	GetTemplateFile() string
	GetIdent() string
	GetEngine() string
	GetEnviron() string
	GetCertLocalDir() string
	GetCertContainerDir() string
}

type CertRenewaler interface {
	RenewalCert(ctx context.Context, domain, email string) error
}
