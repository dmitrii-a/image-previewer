package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsErr(t *testing.T) {
	t.Parallel()
	assert.True(t, IsErr(os.ErrInvalid))
	assert.False(t, IsErr(nil))
}

func TestGetServerAddr(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "localhost:8080", GetServerAddr("localhost", 8080))
	assert.Equal(t, "192.168.1.1:80", GetServerAddr("192.168.1.1", 80))
}

func TestGetNotifyCancelCtx(t *testing.T) {
	t.Parallel()

	ctx, cancel := GetNotifyCancelCtx()
	defer cancel()

	assert.NotNil(t, ctx)
	assert.NotNil(t, cancel)
}
