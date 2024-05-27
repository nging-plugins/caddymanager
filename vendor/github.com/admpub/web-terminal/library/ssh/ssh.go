package ssh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/admpub/errors"
	"github.com/admpub/web-terminal/config"
	websocketx "github.com/admpub/web-terminal/library/websocket"
	"github.com/admpub/websocket"
	"golang.org/x/crypto/ssh"
)

var (
	defaultTimeout = 10 * time.Second
	defaultPort    = 22
)

func New(config *config.SSHConfig) *SSH {
	return &SSH{
		Config: config,
	}
}

type SSH struct {
	Config  *config.SSHConfig
	Client  *ssh.Client
	Session *ssh.Session
	Timeout time.Duration
	stdout  io.Reader
	stderr  io.Reader
	stdin   io.WriteCloser
}

func (s *SSH) Connect() error {
	client, err := s.buildClient()
	if err != nil {
		return errors.New("Failed to dial: " + err.Error())
	}
	s.Client = client

	session, err := client.NewSession()
	if err != nil {
		return errors.New("Failed to create session: " + err.Error())
	}
	s.Session = session

	return nil
}

// buildClient builds the *ssh.Client connection via the jump
// host to the end host. ATM only one jump host is supported
func (s *SSH) buildClient() (*ssh.Client, error) {
	port := defaultPort
	if s.Config.End.Port > 0 {
		port = s.Config.End.Port
	}
	endHostAddress := fmt.Sprintf(`%s:%d`, s.Config.End.Host, port)
	if len(s.Config.Jumps) > 0 { //TODO atm only one jump host is supported
		jumpConfg := s.Config.Jumps[0]
		port := defaultPort
		if jumpConfg.Port > 0 {
			port = jumpConfg.Port
		}
		jumpHostAddress := fmt.Sprintf(`%s:%d`, jumpConfg.Host, port)
		jumpHostClient, err := ssh.Dial("tcp", jumpHostAddress, jumpConfg.ClientConfig)
		if err != nil {
			return nil, fmt.Errorf("ssh.Dial to jump host failed: %w", err)
		}
		jumpHostConn, err := s.dialNextJump(jumpHostClient, endHostAddress)
		if err != nil {
			return nil, fmt.Errorf("ssh.Dial from jump to jump host(%q) failed: %w", endHostAddress, err)
		}

		ncc, chans, reqs, err := ssh.NewClientConn(jumpHostConn, endHostAddress, s.Config.End.ClientConfig)
		if err != nil {
			jumpHostConn.Close()
			return nil, fmt.Errorf("failed to create ssh client to end host(%q): %w", endHostAddress, err)
		}

		return ssh.NewClient(ncc, chans, reqs), nil
	}

	endHostClient, err := ssh.Dial("tcp", endHostAddress, s.Config.End.ClientConfig)
	if err != nil {
		return nil, fmt.Errorf("ssh.Dial directly to end host(%q) failed: %w", endHostAddress, err)
	}

	return endHostClient, nil
}

// dialNextJump dials the next jump host in the chain
func (s *SSH) dialNextJump(jumpHostClient *ssh.Client, nextJumpAddress string) (net.Conn, error) {
	// NOTE: no timeout param in ssh.Dial: https://github.com/golang/go/issues/20288
	// implement it by hand
	var (
		jumpHostConn net.Conn
		err          error
	)
	connChan := make(chan net.Conn)
	go func() {
		jumpHostConn, err = jumpHostClient.Dial("tcp", nextJumpAddress)
		if err != nil {
			return
		}
		connChan <- jumpHostConn
	}()
	connectionTimeout := s.Timeout
	if connectionTimeout <= 0 {
		connectionTimeout = defaultTimeout
	}
	t := time.NewTimer(connectionTimeout)
	defer t.Stop()
	select {
	case jumpHostConnSel := <-connChan:
		jumpHostConn = jumpHostConnSel
	case <-t.C:
		return nil, fmt.Errorf("ssh.Dial from jump host to next jump failed after timeout")
	}

	return jumpHostConn, err
}

func (s *SSH) Close() error {
	if s.Session != nil {
		s.Session.Close()
	}
	if s.stdin != nil {
		s.stdin.Close()
	}
	if s.Client != nil {
		s.Client.Close()
	}
	return nil
}

func (s *SSH) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := websocketx.Connect(w, req,
		func(conn websocketx.Writer) error {
			return s.StartShell(conn, 80, 120)
		},
		s.HandleRecv,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func (s *SSH) StartShell(conn websocketx.Writer, rows, columns int) error {
	return s.StartShellWithCallback(func() error {
		return s.WithZModem(conn)
	}, rows, columns)
}

func (s *SSH) StartShellWithCallback(onInit func() error, rows, columns int) error {
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	err := s.Session.RequestPty("xterm", rows, columns, modes)
	if err != nil {
		return fmt.Errorf("request for pseudo terminal failed: %w", err)
	}
	if onInit != nil {
		if err = onInit(); err != nil {
			return err
		}
	}
	if err = s.Session.Shell(); nil != err {
		return fmt.Errorf("unable to execute command: %w", err)
	}
	if err = s.Session.Wait(); nil != err {
		err = fmt.Errorf("unable to execute command: %w", err)
	}
	return err
}

func (s *SSH) WithZModem(conn websocketx.Writer) error {
	if s.Config.Transform == nil {
		return errors.New(`config.Transform can't be nil`)
	}
	var err error
	s.stdout, s.stderr, s.stdin, err = TransformChannel(s.Session, conn, s.Config.Transform)
	return err
}

func (s *SSH) HandleRecv(conn websocketx.Writer, msgType int, data []byte) error {
	// BinaryMessage 是 zmodem 数据流，则直接发送给 ssh 服务端, 可以提高 rz 上传速率
	if msgType == websocket.BinaryMessage {
		_, err := s.stdin.Write(data)
		if err != nil {
			err = errors.Wrap(err, "write message to ssh error")
		}
		return err
	}
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return errors.Wrap(err, "error format input message")
	}
	switch msg.Type {
	case MessageTypeStdin:
		_, err = s.stdin.Write(msg.Data)
		if err != nil {
			_ = conn.WriteJSON(&Message{Type: MessageTypeStderr, Data: []byte("write to stdin error\r\n")})
			err = errors.Wrap(err, "write to stdin error")
		}
	case MessageTypeResize:
		err = s.Session.WindowChange(msg.Rows, msg.Cols)
		if err != nil {
			_ = conn.WriteJSON(&Message{Type: MessageTypeStderr, Data: []byte("resize error\r\n")})
			err = errors.Wrap(err, "resize error")
		}
	}
	return err
}

func (s *SSH) RunCmd(cmd string) error {
	session, err := s.Client.NewSession()
	if err != nil {
		return errors.New("Failed to create session: " + err.Error())
	}
	session.Stdout = s.Session.Stdout
	session.Stderr = s.Session.Stderr
	defer session.Close()
	return session.Run(cmd)
}

func (s *SSH) RunCmds(r *bytes.Buffer) error {
	for {
		statment, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := s.RunCmd(statment); err != nil {
			return err
		}
	}

	return nil
}
