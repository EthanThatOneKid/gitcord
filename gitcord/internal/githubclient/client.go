package githubclient

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/google/go-github/v47/github"
	"golang.org/x/oauth2"
)

// Client is a wrapped Discord client.
type Client struct {
	*github.Client
	config Config
	ctx    context.Context
}

func (c *Client) WithContext(ctx context.Context) *Client {
	cpy := *c
	cpy.ctx = ctx
	return &cpy
}

type Config struct {
	GitHubOAuth oauth2.TokenSource
	GitHubRepo  string
	Logger      *log.Logger
}

// SplitGitHubRepo splits the GitHub repository path into its owner and name.
func (c Config) SplitGitHubRepo() (owner, repo string) {
	var ok bool
	owner, repo, ok = strings.Cut(c.GitHubRepo, "/")
	if !ok {
		panic("invalid GitHubRepo, must be in owner/repo form")
	}
	return
}

func New(cfg Config) *Client {
	return &Client{
		Client: github.NewClient(oauth2.NewClient(context.Background(), cfg.GitHubOAuth)),
		config: cfg,
	}
}

func (c *Client) logln(v ...any) {
	if c.config.Logger != nil {
		prefixed := []any{"github:"}
		prefixed = append(prefixed, v...)
		c.config.Logger.Println(prefixed...)
	}
}

func (c *Client) FindEvent(eventID int64) (*github.Event, error) {
	owner, repo := c.config.SplitGitHubRepo()
	evs, resp, err := c.Activity.ListRepositoryEvents(c.ctx, owner, repo, nil)
	if err != nil {
		return nil, err
	}

	c.logln("found", len(evs), "events, last page:", resp.LastPage)

	eventIDStr := strconv.FormatInt(eventID, 10)
	for _, ev := range evs {
		if ev.GetID() == eventIDStr {
			return ev, nil
		}
	}

	c.logln("event not found")
	return nil, nil
}
