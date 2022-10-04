package gitcord

import (
	"log"

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

	// ForceCreateThread will force create threads for existing threads.
	ForceCreateThread bool

	// EventID is the ID of the event that is being processed.
	EventID int64

	// Logger is the logger to use. If nil, the default logger will be used.
	Logger *log.Logger
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
	IssueClosed    StatusColors
	PROpened       StatusColors
	PRCommented    StatusColors
	PRClosed       StatusColors
}

var DefaultColorScheme = ColorSchemeConfig{
	IssueOpened: DefaultStatusColors,
	PROpened:    DefaultStatusColors,
}
