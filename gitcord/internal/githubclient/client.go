package githubclient

import (
	"context"
	"fmt"
	"log"
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
	OAuth  oauth2.TokenSource
	Repo   string
	Logger *log.Logger // default log.Default()
}

// SplitGitHubRepo splits the GitHub repository path into its owner and name.
func (c *Config) SplitGitHubRepo() (owner, repo string) {
	var ok bool
	owner, repo, ok = strings.Cut(c.Repo, "/")
	if !ok {
		panic("invalid GitHubRepo, must be in owner/repo form")
	}
	return
}

func New(cfg Config) *Client {
	if cfg.Logger == nil {
		cfg.Logger = log.Default()
	}

	return &Client{
		Client: github.NewClient(oauth2.NewClient(context.Background(), cfg.OAuth)),
		config: cfg,
	}
}

func (c *Client) logln(v ...any) {
	prefixed := []any{"github:"}
	prefixed = append(prefixed, v...)
	c.config.Logger.Println(prefixed...)
}

func (c *Client) EventByID(eventID int64) (*github.Event, error) {
	owner, repo := c.config.SplitGitHubRepo()
	eventIDStr := fmt.Sprintf("%d", eventID)

	var err error
	var resp *github.Response
	var evs []*github.Event
	var nextPage int

	for {
		evs, resp, err = c.Activity.ListRepositoryEvents(c.ctx, owner, repo, &github.ListOptions{
			Page:    nextPage,
			PerPage: 3,
		})
		if err != nil {
			return nil, err
		}

		for _, ev := range evs {
			if ev.GetID() == eventIDStr {
				return ev, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		nextPage = resp.NextPage
	}

	return nil, fmt.Errorf("failed to get event %d", eventID)
}
