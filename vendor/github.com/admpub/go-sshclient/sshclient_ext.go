package sshclient

import (
	"bytes"

	"golang.org/x/crypto/ssh"
)

func New(sshClient *ssh.Client) *Client {
	return &Client{sshClient: sshClient}
}

func NewRemoteScript(sshClient *ssh.Client, opts ...func(*RemoteScript)) *RemoteScript {
	rs := &RemoteScript{
		_type:  rawScript,
		client: sshClient,
	}
	for _, opt := range opts {
		opt(rs)
	}
	return rs
}

func RSCmd(cmd string) func(*RemoteScript) {
	return func(rs *RemoteScript) {
		rs._type = cmdLine
		rs.script = bytes.NewBufferString(cmd + "\n")
	}
}

func RSScript(script string) func(*RemoteScript) {
	return func(rs *RemoteScript) {
		rs._type = rawScript
		rs.script = bytes.NewBufferString(script + "\n")
	}
}

func RSScriptFile(file string) func(*RemoteScript) {
	return func(rs *RemoteScript) {
		rs._type = scriptFile
		rs.scriptFile = file
	}
}

func NewRemoteShell(sshClient *ssh.Client, opts ...func(*RemoteShell)) *RemoteShell {
	rs := &RemoteShell{
		client:     sshClient,
		requestPty: false,
	}
	for _, opt := range opts {
		opt(rs)
	}
	return rs
}

func RShell() func(*RemoteShell) {
	return func(rs *RemoteShell) {
		rs.requestPty = false
		rs.terminalConfig = nil
	}
}

func RTerminal(config *TerminalConfig) func(*RemoteShell) {
	return func(rs *RemoteShell) {
		rs.requestPty = true
		rs.terminalConfig = config
	}
}
