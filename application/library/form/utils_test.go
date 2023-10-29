package form

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDuration(t *testing.T) {
	matches := durationRegex.FindStringSubmatch(`1y2m3d4h5i6s`)
	assert.Equal(t, []string{
		"1y2m3d4h5i6s",
		"1y",
		"2m",
		"3d",
		"4h",
		"5i",
		"6s",
	}, matches)
	v := ParseDuration(`1h`)
	assert.Equal(t, time.Hour, v)
}
