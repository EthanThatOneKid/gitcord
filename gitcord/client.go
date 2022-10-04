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
		Issues:   c.Issues.WithContext(ctx),
		PRs:      c.PRs.WithContext(ctx),
		Comments: c.Comments.WithContext(ctx),
	}
}

func (c Client) DoEvent() error {
	ev, err := c.github.FindEvent(c.config.EventID)
	if err != nil {
		return err
	}

	payload, err := ev.ParsePayload()
	if err != nil {
		return err
	}

	switch *ev.Type {
	case "IssuesEvent":
		return c.handleIssuesEvent(ev, payload.(*github.IssuesEvent))
	case "IssueCommentEvent":
		return c.handleIssueCommentEvent(ev, payload.(*github.IssueCommentEvent))
	case "PullRequestEvent":
		return c.handlePullRequestEvent(ev, payload.(*github.PullRequestEvent))
	case "PullRequestReviewEvent":
		return c.handlePullRequestReviewEvent(ev, payload.(*github.PullRequestReviewEvent))
	case "PullRequestReviewCommentEvent":
		return c.handlePullRequestReviewCommentEvent(ev, payload.(*github.PullRequestReviewCommentEvent))
	case "PullRequestReviewThreadEvent":
		return c.handlePullRequestReviewThreadEvent(ev, payload.(*github.PullRequestReviewThreadEvent))
	default:
		return nil
	}
}

// handleIssuesEvent handles an IssuesEvent
func (c *Client) handleIssuesEvent(ev *github.Event, payload *github.IssuesEvent) error {
	switch *payload.Action {
	case "opened":
		return c.Issues.OpenThread(payload)

	case "closed":
		err := c.Issues.EmbedClosedMsg(payload)
		if err != nil {
			return err
		}
		return c.Issues.EditInitialMsg(payload)

	case "reopened":
		err := c.Issues.EmbedReopenedMsg(payload)
		if err != nil {
			return err
		}
		return c.Issues.EditInitialMsg(payload)

	case "edited":
		return c.Issues.EditInitialMsg(payload)

	case "assigned":
		err := c.Issues.EmbedAssignedMsg(payload)
		if err != nil {
			return err
		}
		return c.Issues.EditInitialMsg(payload)

	case "unassigned":
		err := c.Issues.EmbedUnassignedMsg(payload)
		if err != nil {
			return err
		}
		return c.Issues.EditInitialMsg(payload)

	case "labeled":
		err := c.Issues.EmbedLabeledMsg(payload)
		if err != nil {
			return err
		}
		return c.Issues.EditInitialMsg(payload)

	case "unlabeled":
		err := c.Issues.EmbedUnlabeledMsg(payload)
		if err != nil {
			return err
		}
		return c.Issues.EditInitialMsg(payload)

	case "locked":
		err := c.Issues.EmbedLockedMsg(payload)
		if err != nil {
			return err
		}
		return c.Issues.EditInitialMsg(payload)

	case "unlocked":
		err := c.Issues.EmbedUnlockedMsg(payload)
		if err != nil {
			return err
		}
		return c.Issues.EditInitialMsg(payload)

	case "transferred":
		return c.Issues.EmbedTransferred(payload)

	case "deleted":
		return c.Issues.EmbedDeleted(payload)

	case "milestoned":
		err := c.Issues.EmbedMilestoned(payload)
		if err != nil {
			return err
		}
		return c.Issues.EditInitialMsg(payload)

	case "demilestoned":
		err := c.Issues.EmbedDemilestoned(payload)
		if err != nil {
			return err
		}
		return c.Issues.EditInitialMsg(payload)

	default:
		return nil
	}
}

// handleIssueCommentEvent handles an IssueCommentEvent
func (c *client) handleIssueCommentEvent(ev *github.Event, payload *github.IssueCommentEvent) error {
	switch *payload.Action {
	case "created":
		return c.handleIssueCommentCreated(ev, payload)
	case "edited":
		return c.handleIssueCommentEdited(ev, payload)
	case "deleted":
		return c.handleIssueCommentDeleted(ev, payload)
	default:
		return nil
	}
}

// handlePullRequestEvent handles a PullRequestEvent
func (c *client) handlePullRequestEvent(ev *github.Event, payload *github.PullRequestEvent) error {
	switch *payload.Action {
	case "opened":
		return c.handlePullRequestOpened(ev, payload)
	case "closed":
		return c.handlePullRequestClosed(ev, payload)
	case "reopened":
		return c.handlePullRequestReopened(ev, payload)
	case "edited":
		return c.handlePullRequestEdited(ev, payload)
	case "assigned":
		return c.handlePullRequestAssigned(ev, payload)
	case "unassigned":
		return c.handlePullRequestUnassigned(ev, payload)
	case "review_requested":
		return c.handlePullRequestReviewRequested(ev, payload)
	case "review_request_removed":
		return c.handlePullRequestReviewRequestRemoved(ev, payload)
	case "labeled":
		return c.handlePullRequestLabeled(ev, payload)
	case "unlabeled":
		return c.handlePullRequestUnlabeled(ev, payload)
	case "synchronize":
		return c.handlePullRequestSynchronize(ev, payload)
	case "ready_for_review":
		return c.handlePullRequestReadyForReview(ev, payload)
	case "locked":
		return c.handlePullRequestLocked(ev, payload)
	case "unlocked":
		return c.handlePullRequestUnlocked(ev, payload)
	case "reopened":
		return c.handlePullRequestReopened(ev, payload)
	case "edited":
		return c.handlePullRequestEdited(ev, payload)
	case "assigned":
		return c.handlePullRequestAssigned(ev, payload)
	case "unassigned":
		return c.handlePullRequestUnassigned(ev, payload)
	case "review_requested":
		return c.handlePullRequestReviewRequested(ev, payload)
	case "review_request_removed":
		return c.handlePullRequestReviewRequestRemoved(ev, payload)
	case "labeled":
		return c.handlePullRequestLabeled(ev, payload)
	case "unlabeled":
		return c.handlePullRequestUnlabeled(ev, payload)
	case "synchronize":
		return c.handlePullRequestSynchronize(ev, payload)
	case "ready_for_review":
		return c.handlePullRequestReadyForReview(ev, payload)
	case "locked":
		return c.handlePullRequestLocked(ev, payload)
	case "unlocked":
		return c.handlePullRequestUnlocked(ev, payload)
	default:
		return nil
	}
}

// handlePullRequestReviewEvent handles a PullRequestReviewEvent
func (c *client) handlePullRequestReviewEvent(ev *github.Event, payload *github.PullRequestReviewEvent) error {
	switch *payload.Action {
	case "submitted":
		return c.handlePullRequestReviewSubmitted(ev, payload)
	case "edited":
		return c.handlePullRequestReviewEdited(ev, payload)
	case "dismissed":
		return c.handlePullRequestReviewDismissed(ev, payload)
	default:
		return nil
	}
}

// handlePullRequestReviewCommentEvent handles a PullRequestReviewCommentEvent
func (c *client) handlePullRequestReviewCommentEvent(ev *github.Event, payload *github.PullRequestReviewCommentEvent) error {
	switch *payload.Action {
	case "created":
		return c.handlePullRequestReviewCommentCreated(ev, payload)
	case "edited":
		return c.handlePullRequestReviewCommentEdited(ev, payload)
	case "deleted":
		return c.handlePullRequestReviewCommentDeleted(ev, payload)
	default:
		return nil
	}
}

// handlePullRequestReviewThreadEvent handles a PullRequestReviewThreadEvent
func (c *client) handlePullRequestReviewThreadEvent(ev *github.Event, payload *github.PullRequestReviewThreadEvent) error {
	switch *payload.Action {
	case "created":
		return c.handlePullRequestReviewThreadCreated(ev, payload)
	case "updated":
		return c.handlePullRequestReviewThreadUpdated(ev, payload)
	case "resolved":
		return c.handlePullRequestReviewThreadResolved(ev, payload)
	case "unresolved":
		return c.handlePullRequestReviewThreadUnresolved(ev, payload)
	default:
		return nil
	}
}
