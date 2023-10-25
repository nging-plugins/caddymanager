package engine

type Configer interface {
	GetVhostConfigDirAbsPath() (string, error)
	TemplateFile() string
	Ident() string
	Engine() string
}
