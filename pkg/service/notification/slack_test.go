package notification

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSlackService_SendNotiToDevice(t *testing.T) {
	t.Skip()

	assertT := assert.New(t)
	httpClient := &http.Client{}

	slackService,err  := NewSlackService(httpClient)
	assertT.NoError(err)

	err = slackService.SendNotiToDevice("send slack to channel general")
	assertT.NoError(err)
}
