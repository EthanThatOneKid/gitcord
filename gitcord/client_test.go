package gitcord

import (
	"context"
	"errors"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
	"github.com/diamondburned/arikawa/v3/utils/httputil/httpdriver"
)

type fakeDiscordClient struct {
	*api.Client
	fakeHTTP *fakeHTTPClient
}

func newFakeDiscordClient(records []fakeHTTPRequest) fakeDiscordClient {
	c := fakeDiscordClient{}
	c.fakeHTTP = &fakeHTTPClient{
		records: records,
		current: 0,
	}

	httpClient := httputil.NewClient()
	httpClient.Client = c.fakeHTTP

	c.Client = api.NewCustomClient("", httpClient)
	return c
}

type fakeHTTPRequest struct {
	Request  httpdriver.MockRequest
	Response httpdriver.MockResponse
}

type fakeHTTPClient struct {
	records []fakeHTTPRequest
	excess  []httpdriver.Request
	current int
}

var _ httpdriver.Client = (*fakeHTTPClient)(nil)

func (c *fakeHTTPClient) NewRequest(ctx context.Context, method, url string) (httpdriver.Request, error) {
	return httpdriver.NewMockRequestWithContext(ctx, method, url, nil, nil), nil
}

func (c *fakeHTTPClient) Do(req httpdriver.Request) (httpdriver.Response, error) {
	if c.current >= len(c.records) {
		c.excess = append(c.excess, req)
		return nil, errors.New("no more requests in fake client")
	}

	current := &c.records[c.current]
	c.current++

	if err := httpdriver.ExpectMockRequest(&current.Request, req); err != nil {
		c.excess = append(c.excess, req)
		return nil, err
	}

	return &current.Response, nil
}
