package gitcord

import (
	"log"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"golang.org/x/oauth2"
)

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

type StatusColors struct {
	Success discord.Color
	Error   discord.Color
}

var DefaultStatusColors = StatusColors{
	Success: 0x00FF00,
	Error:   0xFF0000,
}

type ColorSchemeConfig struct {
	IssueOpened    StatusColors
	IssueCommented StatusColors
	PROpened       StatusColors
	PRCommented    StatusColors
}

var DefaultColorScheme = ColorSchemeConfig{
	IssueOpened: DefaultStatusColors,
	PROpened:    DefaultStatusColors,
}
