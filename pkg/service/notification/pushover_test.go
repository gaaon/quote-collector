package notification

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendNotiToDevice(t *testing.T) {
	t.Skip()
	err := SendNotiToDevice("test message", "test title")
	assert.NoError(t, err)
}