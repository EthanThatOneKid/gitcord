package gitcord

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/google/go-github/v47/github"
	"golang.org/x/oauth2"
)

// client for the GitHub Discord bot
type client struct {
	github  *github.Client
	discord *api.Client
	config  Config
}

type Client struct {
	Issues   *IssuesClient
	PRs      *PRsClient
	Comments *IssueCommentClient
}

func NewClient(cfg Config) *Client {
	c := newClient(cfg)
	return &Client{
		Issues:   &IssuesClient{c, context.Background()},
		PRs:      &PRsClient{c, context.Background()},
		Comments: &IssueCommentClient{c, context.Background()},
	}
}

func newClient(cfg Config) *client {
	return &client{
		github:  github.NewClient(oauth2.NewClient(context.Background(), cfg.GitHubOAuth)),
		discord: api.NewClient(cfg.DiscordToken),
		config:  cfg,
	}
}

func (c *client) logln(v ...any) {
	if c.config.Logger != nil {
		c.config.Logger.Println(v...)
	}
}

func (c *Client) WithContext(ctx context.Context) *Client {
	return &Client{
		Issues:   c.Issues.WithContext(ctx),
		PRs:      c.PRs.WithContext(ctx),
		Comments: c.Comments.WithContext(ctx),
	}
}
