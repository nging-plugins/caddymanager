package engine

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
)

func testContainerExec(ctx context.Context, containerID string, cmd []string, env []string, outWriter io.Writer, errWriter io.Writer) error {
	outWriter.Write([]byte(`OK`))
	return nil
}

func TestGetContainerExec(t *testing.T) {
	exec0 := getContainerExec()
	assert.True(t, nil == exec0)

	echo.Set(`DockerContainerExec`, testContainerExec)
	exec := getContainerExec()
	t.Logf(`%T`, exec)
	buf := bytes.NewBuffer(nil)
	exec(nil, ``, nil, nil, buf, nil)
	assert.Equal(t, `OK`, buf.String())
}

func TestParseUnixSocketAddr(t *testing.T) {
	addr := `unix://var/run/docker.sock`
	expected := `/var/run/docker.sock`
	addr = ParseSocketAddr(addr)
	assert.Equal(t, expected, addr)

	addr = `unix:/var/run/docker.sock`
	addr = ParseSocketAddr(addr)
	assert.Equal(t, expected, addr)

	addr = `unix:var/run/docker.sock`
	addr = ParseSocketAddr(addr)
	assert.Equal(t, expected, addr)

	addr = `unix://///var/run/docker.sock`
	addr = ParseSocketAddr(addr)
	assert.Equal(t, expected, addr)
}

func TestHTTPSocket(t *testing.T) {
	addr := `unix://var/run/docker.sock`
	sockAddr := ParseSocketAddr(addr)
	client := newRestyClient(NewSocketClient(sockAddr))
	resp, err := client.R().Get(`http://sock/images/json`)
	assert.NoError(t, err)
	t.Logf(resp.String())
}
