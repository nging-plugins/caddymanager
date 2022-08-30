package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallSQL(t *testing.T) {
	assert.NotEmpty(t, installSQL)
}
