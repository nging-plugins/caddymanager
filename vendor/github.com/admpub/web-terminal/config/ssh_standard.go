package config

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"golang.org/x/crypto/ssh"
)

func NewSSHStandardWithAccount(ws io.ReadWriter, account *AccountConfig) (*ssh.ClientConfig, error) {
	return NewSSHStandard(ws, ws, account)
}

func NewHostConfigWithAccount(ws io.ReadWriter, account *AccountConfig, host string, port int) (*HostConfig, error) {
	clientConfig, err := NewSSHStandardWithAccount(ws, account)
	if err != nil {
		return nil, err
	}
	return &HostConfig{
		ClientConfig: clientConfig,
		Host:         host,
		Port:         port,
	}, err
}

func NewSSHStandard(reader io.Reader, writer io.Writer, account *AccountConfig) (*ssh.ClientConfig, error) {
	// Dial code is taken from the ssh package example
	sshConfig := &ssh.ClientConfig{
		Config:          ssh.Config{Ciphers: supportedCiphers},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		User:            account.User,
		Auth:            []ssh.AuthMethod{},
	}
	account.SetDefault()
	if account.PrivateKey != nil {
		var signer ssh.Signer
		var err error
		pemBytes := account.PrivateKey
		if account.Passphrase != nil {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, account.Passphrase)
		} else {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		}
		if err != nil {
			return sshConfig, err
		}
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(signer))
	}

	if len(account.Password) > 0 {
		sshConfig.Auth = append(sshConfig.Auth, ssh.Password(account.Password))
		if reader != nil && writer != nil {
			bufReader := bufio.NewReader(reader)
			sshConfig.Auth = append(sshConfig.Auth, ssh.KeyboardInteractive(KeyboardInteractivefunc(bufReader, writer, account)))
		}
	}

	sshConfig.SetDefaults()
	return sshConfig, nil
}

func KeyboardInteractivefunc(reader *bufio.Reader, writer io.Writer, account *AccountConfig) func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
	var (
		passwordCount         int
		emptyInteractiveCount int
		password              = account.Password
		charset               = account.Charset
	)
	return func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
		if len(questions) == 0 {
			emptyInteractiveCount++
			if emptyInteractiveCount++; emptyInteractiveCount > 50 {
				return nil, errors.New("interactive count is too much")
			}
			return nil, nil
		}
		for _, question := range questions {
			io.WriteString(DecodeBy(charset, writer), question)
			//panic(question)
			question = strings.TrimSpace(question)
			if len(question) > 12 {
				question = question[0:12]
			}
			question = strings.ToLower(question)
			switch question {
			case "password:", "password as ", "password for":
				passwordCount++
				if passwordCount == 1 {
					answers = append(answers, password)
					break
				}
				fallthrough
			default:
				line, _, e := reader.ReadLine()
				if nil != e {
					return nil, e
				}
				answers = append(answers, string(line))
			}
		}
		return answers, nil
	}
}
