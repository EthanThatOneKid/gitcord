package gitcord

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/utils/httputil/httpdriver"
)

type fakeDiscordClient struct {
	api.Client
	records []fakeHTTPRequest
}

type fakeHTTPRequest struct {
	Request  httpdriver.MockRequest
	Response httpdriver.MockResponse
}

// func newFakeDiscordClient(records [])
