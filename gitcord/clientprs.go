package gitcord

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/google/go-github/v47/github"
	"github.com/pkg/errors"
)

type PRsClient struct {
	*client
	ctx context.Context
}

func (c *PRsClient) WithContext(ctx context.Context) *PRsClient {
	cpy := *c
	cpy.client = cpy.client.WithContext(ctx)
	cpy.ctx = ctx
	return &cpy
}

func (c *PRsClient) OpenAndEmbedInitialMsg(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch != nil {
		if ch.Type != discord.GuildPublicThread {
			if !c.config.ForceOpen {
				return fmt.Errorf("channel %d is not a public thread", ch.ID)
			}
			c.logln(fmt.Sprintf("ignoring channel %d is not a public thread", ch.ID))
		}

		if !c.config.ForceOpen {
			return fmt.Errorf("pull request %d already has a thread %d", pr.GetNumber(), ch.ID)
		}
		c.logln(fmt.Sprintf("ignoring existing thread %d", ch.ID))
	}

	ch, err := c.discord.StartThreadWithoutMessage(c.config.DiscordChannelID, api.StartThreadData{
		Name:                fmt.Sprintf("#%d: %s", pr.GetNumber(), pr.GetTitle()),
		Type:                discord.GuildPublicThread,
		AutoArchiveDuration: discord.SevenDaysArchive,
	})
	if err != nil {
		return errors.Wrap(err, "failed to open thread")
	}

	_, err = c.discord.SendEmbeds(ch.ID, c.config.makePREmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c *PRsClient) EditInitialMsg(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	msg := c.discord.FindMsgByIssue(ch, pr.GetNumber())
	if msg == nil {
		return fmt.Errorf("pull request %d does not have an initial message", pr.GetNumber())
	}

	_, err := c.discord.EditEmbeds(ch.ID, msg.ID, c.config.makePREmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to edit message")
	}

	return nil
}

func (c *PRsClient) EmbedClosedMsg(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRClosedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c *PRsClient) EmbedReopenedMsg(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRReopenedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c PRsClient) EmbedDeletedMsg(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRDeletedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c PRsClient) EmbedTransferredMsg(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRTransferredEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c PRsClient) EmbedAssignedMsg(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRAssignedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c PRsClient) EmbedUnassignedMsg(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRUnassignedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c PRsClient) EmbedLabeledMsg(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRLabeledEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c PRsClient) EmbedUnlabeledMsg(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRUnlabeledEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c PRsClient) EmbedLockedMsg(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRLockedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c PRsClient) EmbedUnlockedMsg(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRUnlockedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c PRsClient) EmbedMilestonedMsg(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRMilestonedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c PRsClient) EmbedDemilestonedMsg(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRDemilestonedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c PRsClient) EmbedReviewRequestedMsg(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRReviewRequestedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c PRsClient) EmbedReviewRequestRemovedMsg(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRReviewRequestRemovedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c PRsClient) EmbedReadyForReviewMsg(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRReadyForReviewEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}
