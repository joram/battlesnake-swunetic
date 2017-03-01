package swu

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var testTime = time.Date(2017, 2, 28, 0, 0, 0, 0, time.UTC)

func TestUnixMilli(t *testing.T) {
	val := TimeToUnixMilli(testTime)
	assert.Equal(t, int64(1488240000000), val)
}

func TestUnixMilliToNano(t *testing.T) {
	val := UnixMilliToNano(1488240000000)
	assert.Equal(t, testTime.UnixNano(), val)
}

func TestUnixNanoToMilli(t *testing.T) {
	val := UnixNanoToMilli(testTime.UnixNano())
	assert.Equal(t, int64(1488240000000), val)
}
