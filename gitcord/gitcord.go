package gitcord

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/google/go-github/v47/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type StatusColors struct {
	Success discord.Color
	Error   discord.Color
}

var DefaultStatusColors = StatusColors{
	Success: 0x00FF00,
	Error:   0xFF0000,
}

type ColorSchemeConfig struct {
	IssueOpened StatusColors
	PROpened    StatusColors
}

var DefaultColorScheme = ColorSchemeConfig{
	IssueOpened: DefaultStatusColors,
	PROpened:    DefaultStatusColors,
}

// Config is the configuration for the Client.
type Config struct {
	// GitHubOAuth is the GitHub OAuth token.
	GitHubOAuth oauth2.TokenSource
	// GitHubRepo is the owner and name of the repository, formatted as: <owner>/<name>
	GitHubRepo string
	// DiscordToken is the Discord bot token.
	DiscordToken string
	// DiscordChannelID is the ID of the parent channel in which all threads
	// will be created under.
	DiscordChannelID discord.ChannelID
	DiscordGuildID   discord.GuildID
	Colors           ColorSchemeConfig

	// Logger is the logger to use. If nil, the default logger will be used.
	Logger *log.Logger
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

// Client is the client for the GitHub Discord bot.
type Client struct {
	github  *github.Client
	discord *api.Client
	config  Config
	ctx     context.Context
}

func NewClient(cfg Config) *Client {
	return &Client{
		github:  github.NewClient(oauth2.NewClient(context.Background(), cfg.GitHubOAuth)),
		discord: api.NewClient(cfg.DiscordToken),
		config:  cfg,
		ctx:     context.Background(),
	}
}

func (c *Client) logln(v ...any) {
	if c.config.Logger != nil {
		c.config.Logger.Println(v...)
	}
}

func (c *Client) WithContext(ctx context.Context) *Client {
	return &Client{
		github:  c.github,
		discord: c.discord.WithContext(ctx),
		config:  c.config,
		ctx:     ctx,
	}
}

func (c *Client) CreateIssueThreadt(c.ctx, owner, repo, issueNumber)
	if err != nil {
		return errors.Wrap(err, "failed to get issue")
	}

	ch, err := c.discord.StartThreadWithoutMessage(c.config.DiscordChannelID, api.StartThreadData{
		Name:                fmt.Sprintf("#%d: %s", issue.GetNumber(), issue.GetTitle()),
		Type:                discord.GuildPublicThread,
		AutoArchiveDuration: discord.OneDayArchive,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create channel")
	}

	_, err = c.discord.SendEmbeds(ch.ID, discord.Embed{
		Title: fmt.Sprintf("#%d: %s", issue.GetNumber(), issue.GetTitle()),
		URL:   issue.GetURL(),
		Author: &discord.EmbedAuthor{
			Name: issue.GetUser().GetLogin(),
			Icon: issue.GetUser().GetAvatarURL(),
		},
		// TODO: convert GitHub Markdown to valid Discord Markdown.
		// See https://github.com/pythonian23/SMoRe.
		Description: issue.GetBody(),
	})
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

// TODO: If channel is empty, do not archive.
func (c *Client) CloseIssueChannel(issue int) error {
	return nil
}

func (c *Client) SendPRReview(prID int, commentID int64) error {
	owner, repo := c.config.SplitGitHubRepo()

	comment, _, err := c.github.PullRequests.GetComment(c.ctx, owner, repo, commentID)
	if err != nil {
		return errors.Wrap(err, "failed to get comment")
	}

	pr, _, err := c.github.PullRequests.Get(c.ctx, owner, repo, prID)
	if err != nil {
		return errors.Wrap(err, "failed to get PR")
	}

	threads, err := c.activeThreads()
	if err != nil {
		return errors.Wrap(err, "failed to get active threads")
	}

	thread := findChannelNo(threads, prID)
	if thread == nil {
		c.logln("skipping unknown PR with ID", prID)
		return nil
	}

	content, err := renderTmpl(prReviewMessageTmpl, map[string]any{
		"PR":      pr,
		"Comment": comment,
		"Config":  c.config,
	})
	if err != nil {
		return errors.Wrap(err, "failed to render template")
	}

	_, err = c.discord.SendEmbeds(thread.ID, discord.Embed{
		Title: fmt.Sprintf("#%d comment %d", pr.GetNumber(), comment.GetID()),
		URL:   comment.GetURL(),
		Author: &discord.EmbedAuthor{
			Name: comment.GetUser().GetLogin(),
			Icon: comment.GetUser().GetAvatarURL(),
		},
		Description: convertMarkdown(content),
	})
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func convertMarkdown(githubMD string) string {
	// See https://github.com/pythonian23/SMoRe.
	return githubMD
}

func (c *Client) activeThreads() ([]discord.Channel, error) {
	active, err := c.discord.ActiveThreads(c.config.DiscordGuildID)
	if err != nil {
		return nil, err
	}

	relevantThreads := active.Threads[:0]
	for _, thread := range active.Threads {
		if thread.ParentID == c.config.DiscordChannelID {
			relevantThreads = append(relevantThreads, thread)
		}
	}

	return relevantThreads, nil
}

func findChannel(channels []discord.Channel, f func(ch *discord.Channel) bool) *discord.Channel {
	for i := range channels {
		if f(&channels[i]) {
			return &channels[i]
		}
	}
	return nil
}

func findChannelNo(channels []discord.Channel, targetID int) *discord.Channel {
	return findChannel(channels, func(ch *discord.Channel) bool {
		var n int
		_, err := fmt.Scanf("#%d", &n)
		return err == nil && n == targetID
	})
}
