package gitcord

import (
	"context"

	"github.com/ethanthatonekid/gitcord/gitcord/internal/discordclient"
	"github.com/ethanthatonekid/gitcord/gitcord/internal/githubclient"

	"github.com/google/go-github/v47/github"
)

// client for the GitHub Discord bot
type client struct {
	github  *githubclient.Client
	discord *discordclient.Client
	config  Config
}

type Client struct {
	*client

	Issues         *IssuesClient
	Comments       *IssueCommentClient
	PRs            *PRsClient
	Reviews        *ReviewsClient
	ReviewComments *ReviewCommentsClient
	ReviewThreads  *ReviewThreadsClient
}

func NewClient(cfg Config) *Client {
	c := newClient(cfg)
	return &Client{
		Issues:         &IssuesClient{c, context.Background()},
		Comments:       &IssueCommentClient{c, context.Background()},
		PRs:            &PRsClient{c, context.Background()},
		Reviews:        &ReviewsClient{c, context.Background()},
		ReviewComments: &ReviewCommentsClient{c, context.Background()},
		ReviewThreads:  &ReviewThreadsClient{c, context.Background()},
	}
}

func newClient(cfg Config) *client {
	return &client{
		github: githubclient.New(githubclient.Config{
			GitHubOAuth: cfg.GitHubOAuth,
			Logger:      cfg.Logger,
		}),
		discord: discordclient.New(discordclient.Config{
			DiscordToken:     cfg.DiscordToken,
			DiscordGuildID:   cfg.DiscordGuildID,
			DiscordChannelID: cfg.DiscordChannelID,
			Logger:           cfg.Logger,
		}),
		config: cfg,
	}
}

func (c *client) logln(v ...any) {
	if c.config.Logger != nil {
		c.config.Logger.Println(v...)
	}
}

func (c *client) WithContext(ctx context.Context) *client {
	return &client{
		github:  c.github.WithContext(ctx),
		discord: c.discord.WithContext(ctx),
		config:  c.config,
	}
}

func (c *Client) WithContext(ctx context.Context) *Client {
	return &Client{
		Issues:         c.Issues.WithContext(ctx),
		Comments:       c.Comments.WithContext(ctx),
		PRs:            c.PRs.WithContext(ctx),
		Reviews:        c.Reviews.WithContext(ctx),
		ReviewComments: c.ReviewComments.WithContext(ctx),
		ReviewThreads:  c.ReviewThreads.WithContext(ctx),
	}
}

// DoEvent handles a GitHub event
//
// https://docs.github.com/en/developers/webhooks-and-events/events/github-event-types
func (c Client) DoEvent() error {
	ev, err := c.github.FindEvent(c.config.EventID)
	if err != nil {
		return err
	}

	data, err := ev.ParsePayload()
	if err != nil {
		return err
	}

	switch *ev.Type {
	case "IssuesEvent":
		return c.handleIssuesEvent(data.(*github.IssuesEvent))
	case "IssueCommentEvent":
		return c.handleIssueCommentEvent(data.(*github.IssueCommentEvent))
	case "PullRequestEvent":
		return c.handlePREvent(data.(*github.PullRequestEvent))
	case "PullRequestReviewEvent":
		return c.handlePullRequestReviewEvent(data.(*github.PullRequestReviewEvent))
	case "PullRequestReviewCommentEvent":
		return c.handlePullRequestReviewCommentEvent(data.(*github.PullRequestReviewCommentEvent))
	case "PullRequestReviewThreadEvent":
		return c.handlePullRequestReviewThreadEvent(data.(*github.PullRequestReviewThreadEvent))
	default:
		return nil
	}
}

// handleIssuesEvent handles an IssuesEvent
//
// https://docs.github.com/en/developers/webhooks-and-events/events/github-event-types#issuesevent
func (c Client) handleIssuesEvent(ev *github.IssuesEvent) error {
	switch *ev.Action {
	case "opened":
		return c.Issues.OpenAndEmbedInitialMsg(ev)
	case "edited":
		return c.Issues.EditInitialMsg(ev)
	case "deleted":
		return c.Issues.EmbedDeletedMsg(ev)
	case "transferred":
		return c.Issues.EmbedTransferredMsg(ev)

	case "closed", "reopened", "assigned", "unassigned", "labeled", "unlabeled", "locked", "unlocked", "milestoned", "demilestoned":
		var err error
		switch *ev.Action {
		case "closed":
			err = c.Issues.EmbedClosedMsg(ev)
		case "reopened":
			err = c.Issues.EmbedReopenedMsg(ev)
		case "assigned":
			err = c.Issues.EmbedAssignedMsg(ev)
		case "unassigned":
			err = c.Issues.EmbedUnassignedMsg(ev)
		case "labeled":
			err = c.Issues.EmbedLabeledMsg(ev)
		case "unlabeled":
			err = c.Issues.EmbedUnlabeledMsg(ev)
		case "locked":
			err = c.Issues.EmbedLockedMsg(ev)
		case "unlocked":
			err = c.Issues.EmbedUnlockedMsg(ev)
		case "milestoned":
			err = c.Issues.EmbedMilestonedMsg(ev)
		case "demilestoned":
			err = c.Issues.EmbedDemilestonedMsg(ev)
		}

		if err != nil {
			return err
		}

		return c.Issues.EditInitialMsg(ev)

	default:
		return nil
	}
}

// handlePullRequestEvent handles a PullRequestEvent
//
// https://docs.github.com/en/developers/webhooks-and-events/events/github-event-types#pullrequestevent
func (c *Client) handlePREvent(ev *github.PullRequestEvent) error {
	switch *ev.Action {
	case "opened":
		return c.PRs.OpenAndEmbedInitialMsg(ev)
	case "edited":
		return c.PRs.EditInitialMsg(ev)
	case "deleted":
		return c.PRs.EmbedDeletedMsg(ev)
	case "transferred":
		return c.PRs.EmbedTransferredMsg(ev)
	case "review_requested":
		return c.PRs.EmbedReviewRequestedMsg(ev)
	case "review_request_removed":
		return c.PRs.EmbedReviewRequestRemovedMsg(ev)
	case "ready_for_review":
		return c.PRs.EmbedReadyForReviewMsg(ev)

	case "closed", "reopened", "assigned", "unassigned", "labeled", "unlabeled", "locked", "unlocked", "milestoned", "demilestoned":
		var err error
		switch *ev.Action {
		case "closed":
			err = c.PRs.EmbedClosedMsg(ev)
		case "reopened":
			err = c.PRs.EmbedReopenedMsg(ev)
		case "assigned":
			err = c.PRs.EmbedAssignedMsg(ev)
		case "unassigned":
			err = c.PRs.EmbedUnassignedMsg(ev)
		case "labeled":
			err = c.PRs.EmbedLabeledMsg(ev)
		case "unlabeled":
			err = c.PRs.EmbedUnlabeledMsg(ev)
		case "locked":
			err = c.PRs.EmbedLockedMsg(ev)
		case "unlocked":
			err = c.PRs.EmbedUnlockedMsg(ev)
		case "milestoned":
			err = c.PRs.EmbedMilestonedMsg(ev)
		case "demilestoned":
			err = c.PRs.EmbedDemilestonedMsg(ev)
		}

		if err != nil {
			return err
		}

		return c.PRs.EditInitialMsg(ev)

	default:
		return nil
	}
}

// handleIssueCommentEvent handles an IssueCommentEvent
//
// https://docs.github.com/en/developers/webhooks-and-events/events/github-event-types#issuecommentevent
func (c Client) handleIssueCommentEvent(ev *github.IssueCommentEvent) error {
	switch *ev.Action {
	case "created":
		return c.Comments.EmbedIssueCommentMsg(ev)
	case "edited":
		return c.Comments.EditIssueCommentMsg(ev)
	case "deleted":
		return c.Comments.EmbedDeletedMsg(ev)
	default:
		return nil
	}
}

// handlePullRequestReviewEvent handles a PullRequestReviewEvent
//
// https://docs.github.com/en/developers/webhooks-and-events/events/github-event-types#pullrequestreviewevent
func (c *Client) handlePullRequestReviewEvent(ev *github.PullRequestReviewEvent) error {
	switch *ev.Action {
	case "submitted":
		return c.Reviews.EmbedReviewMsg(ev)
	case "edited":
		return c.Reviews.EditReviewMsg(ev)
	case "dismissed":
		return c.Reviews.EmbedReviewDismissedMsg(ev)
	default:
		return nil
	}
}

// handlePullRequestReviewCommentEvent handles a PullRequestReviewCommentEvent
//
// https://docs.github.com/en/developers/webhooks-and-events/events/github-event-types#pullrequestreviewcommentevent
func (c *Client) handlePullRequestReviewCommentEvent(ev *github.PullRequestReviewCommentEvent) error {
	switch *ev.Action {
	case "created":
		return c.ReviewComments.EmbedReviewCommentMsg(ev)
	case "edited":
		return c.ReviewComments.EditReviewCommentMsg(ev)
	case "deleted":
		return c.ReviewComments.EmbedReviewCommentDeletedMsg(ev)
	default:
		return nil
	}
}

// handlePullRequestReviewThreadEvent handles a PullRequestReviewThreadEvent
//
// https://docs.github.com/en/developers/webhooks-and-events/events/github-event-types#pullrequestreviewthreadevent
func (c *Client) handlePullRequestReviewThreadEvent(ev *github.PullRequestReviewThreadEvent) error {
	switch *ev.Action {
	case "created":
		return c.ReviewThreads.EmbedReviewThreadMsg(ev)
	case "updated":
		return c.ReviewThreads.EditReviewThreadMsg(ev)
	case "resolved":
		return c.ReviewThreads.EmbedReviewThreadResolvedMsg(ev)
	case "unresolved":
		return c.ReviewThreads.EmbedReviewThreadUnresolvedMsg(ev)
	default:
		return nil
	}
}
