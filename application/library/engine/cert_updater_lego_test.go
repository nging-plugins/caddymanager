package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLegoInstall(t *testing.T) {
	return
	tarFile := `/Users/hank/Downloads/lego_v4.22.2_linux_amd64.tar.gz`
	err := legoInstall(tarFile)
	assert.NoError(t, err)
}
