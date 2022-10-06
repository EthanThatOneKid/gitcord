package gitcord

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"golang.org/x/oauth2"
)

// Config is the configuration for the Client
type Config struct {
	// GitHubOAuth is the GitHub OAuth token
	GitHubOAuth oauth2.TokenSource
	// GitHubRepo is the owner and name of the repository, formatted as: <owner>/<name>
	GitHubRepo string
	// DiscordToken is the Discord bot token
	DiscordToken string
	// DiscordChannelID is the ID of the parent channel in which all threads
	// will be created under
	DiscordChannelID discord.ChannelID
	DiscordGuildID   discord.GuildID
	Colors           ColorSchemeConfig
	// ForceOpen will force create a new thread for existing threads
	ForceOpen bool
	// Logger is the logger to use. If nil, the default logger will be used
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
	IssueOpened            StatusColors
	IssueClosed            StatusColors
	IssueReopened          StatusColors
	IssueLabeled           StatusColors
	IssueUnlabeled         StatusColors
	IssueAssigned          StatusColors
	IssueUnassigned        StatusColors
	IssueMilestoned        StatusColors
	IssueDemilestoned      StatusColors
	IssueDeleted           StatusColors
	IssueLocked            StatusColors
	IssueUnlocked          StatusColors
	IssueTransferred       StatusColors
	IssueCommented         StatusColors
	IssueCommentDeleted    StatusColors
	PROpened               StatusColors
	PRReopened             StatusColors
	PRCommented            StatusColors
	PRClosed               StatusColors
	PRAssigned             StatusColors
	PRUnassigned           StatusColors
	PRDeleted              StatusColors
	PRTransferred          StatusColors
	PRLabeled              StatusColors
	PRUnlabeled            StatusColors
	PRMilestoned           StatusColors
	PRDemilestoned         StatusColors
	PRLocked               StatusColors
	PRUnlocked             StatusColors
	PRReviewRequested      StatusColors
	PRReviewRequestRemoved StatusColors
	PRReadyForReview       StatusColors
	Reviewed               StatusColors
	ReviewDismissed        StatusColors
	ReviewCommented        StatusColors
	ReviewCommentDeleted   StatusColors
	ReviewThreaded         StatusColors
	ReviewThreadResolved   StatusColors
	ReviewThreadUnresolved StatusColors
}

var DefaultColorScheme = ColorSchemeConfig{
	IssueOpened: DefaultStatusColors,
	PROpened:    DefaultStatusColors,
}
