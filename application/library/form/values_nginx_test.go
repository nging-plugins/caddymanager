package form

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsNginxLogFormat(t *testing.T) {
	value := `{combined}`
	logFormat := AsNginxLogFormat(value)
	assert.Equal(t, `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent`, logFormat)
}
