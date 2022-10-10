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
	// ColorScheme is the color scheme for use in embeds. Refer to ColorScheme
	// for more information.
	ColorScheme ColorScheme
	// ForceOpen will force create a new thread even if one already exists
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

// ColorSchemeKey is the key for a color within a color scheme.
type ColorSchemeKey uint

const (
	UnknownColorSchemeKey ColorSchemeKey = iota
	IssueOpened
	IssueClosed
	IssueReopened
	IssueLabeled
	IssueUnlabeled
	IssueAssigned
	IssueUnassigned
	IssueMilestoned
	IssueDemilestoned
	IssueDeleted
	IssueLocked
	IssueUnlocked
	IssueTransferred
	IssueCommented
	IssueCommentDeleted
	PROpened
	PRReopened
	PRCommented
	PRClosed
	PRAssigned
	PRUnassigned
	PRDeleted
	PRTransferred
	PRLabeled
	PRUnlabeled
	PRMilestoned
	PRDemilestoned
	PRLocked
	PRUnlocked
	PRReviewRequested
	PRReviewRequestRemoved
	PRReadyForReview
	Reviewed
	ReviewDismissed
	ReviewCommented
	ReviewCommentDeleted
	ReviewThreaded
	ReviewThreadResolved
	ReviewThreadUnresolved

	maxColorSchemeKey // internal use only
)

// ColorScheme describes the color scheme for all embed colors made by gitcord.
// It maps each color scheme key to a status color struct, which has two
// possible colors for two cases.
//
// By default, all color scheme keys ap to DefaultStatusColors.
type ColorScheme map[ColorSchemeKey]StatusColors

// DefaultColorScheme is the default color scheme.
var DefaultColorScheme = ColorScheme{}

func init() {
	for i := 0; i < int(maxColorSchemeKey); i++ {
		DefaultColorScheme[ColorSchemeKey(i)] = DefaultStatusColors
	}
}

// Override creates a new ColorScheme that overrides all color keys inside s
// with the provided ones in with.
func (s ColorScheme) Override(with ColorScheme) ColorScheme {
	newScheme := make(ColorScheme, len(s))
	for k, v := range s {
		newScheme[k] = v
	}

	for k, v := range with {
		newScheme[k] = v
	}

	return newScheme
}

// Color gets the color corresponding to the given key. If success is true, then
// the Success color is used, else Error is used.
func (s ColorScheme) Color(k ColorSchemeKey, success bool) discord.Color {
	colors, ok := s[k]
	if !ok {
		colors = DefaultStatusColors
	}

	if success {
		return colors.Success
	} else {
		return colors.Error
	}
}
