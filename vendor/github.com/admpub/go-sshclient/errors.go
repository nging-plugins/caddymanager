package sshclient

import (
	"errors"

	"github.com/hashicorp/go-multierror"
)

type multiError struct {
	*multierror.Error
}

func newMultiError() *multiError {
	return &multiError{&multierror.Error{}}
}

func (e *multiError) Append(err error) error {
	if err == nil {
		return e.Error.ErrorOrNil()
	}
	_ = multierror.Append(e.Error, err)
	return e.Error.ErrorOrNil()
}

func (e multiError) AppendErrFunc(errFunc func() error) error {
	return e.Append(errFunc())
}

var (
	ErrNotSupportedRemoteScript = errors.New("not supported RemoteScript type")
	ErrStdoutAlreadySet         = errors.New("stdout already set")
	ErrStderrAlreadySet         = errors.New("stderr already set")
)
