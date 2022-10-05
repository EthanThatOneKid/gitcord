package gitcord

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/google/go-github/v47/github"
	"github.com/pkg/errors"
)

type IssuesClient client

func (c *IssuesClient) OpenAndEmbedInitialMsg(ev *github.IssuesEvent) error {
	issue := ev.GetIssue()

	ch := c.discord.FindThreadByNumber(issue.GetNumber())
	if ch != nil {
		if ch.Type != discord.GuildPublicThread {
			if !c.config.ForceOpen {
				return fmt.Errorf("channel %d is not a public thread", ch.ID)
			}
			// c.logln(fmt.Sprintf("ignoring channel %d is not a public thread", ch.ID))
		}

		if !c.config.ForceOpen {
			return fmt.Errorf("issue %d already has a thread %d", issue.GetNumber(), ch.ID)
		}
		// c.logln(fmt.Sprintf("ignoring existing thread %d", ch.ID))
	}

	ch, err := c.discord.StartThreadWithoutMessage(c.config.DiscordChannelID, api.StartThreadData{
		Name:                fmt.Sprintf("#%d: %s", issue.GetNumber(), issue.GetTitle()),
		Type:                discord.GuildPublicThread,
		AutoArchiveDuration: discord.SevenDaysArchive,
	})
	if err != nil {
		return errors.Wrap(err, "failed to open thread")
	}

	_, err = c.discord.SendEmbeds(ch.ID, c.config.makeIssueEmbed(issue))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c *IssuesClient) EditInitialMsg(ev *github.IssuesEvent) error {
	issue := ev.GetIssue()

	ch := c.discord.FindThreadByNumber(issue.GetNumber())
	if ch == nil {
		return fmt.Errorf("issue %d does not have a thread", issue.GetNumber())
	}

	msg := c.discord.FindMsgByIssue(ch, issue.GetNumber())
	if msg == nil {
		return fmt.Errorf("issue %d does not have an initial message", issue.GetNumber())
	}

	_, err := c.discord.EditEmbeds(ch.ID, msg.ID, c.config.makeIssueEmbed(issue))
	if err != nil {
		return errors.Wrap(err, "failed to edit message")
	}

	return nil
}

func (c *IssuesClient) EmbedClosedMsg(ev *github.IssuesEvent) error {
	issue := ev.GetIssue()

	ch := c.discord.FindThreadByNumber(issue.GetNumber())
	if ch == nil {
		return fmt.Errorf("issue %d does not have a thread", issue.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makeIssueClosedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c *IssuesClient) EmbedReopenedMsg(ev *github.IssuesEvent) error {
	issue := ev.GetIssue()

	ch := c.discord.FindThreadByNumber(issue.GetNumber())
	if ch == nil {
		return fmt.Errorf("issue %d does not have a thread", issue.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makeIssueReopenedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c IssuesClient) EmbedDeletedMsg(ev *github.IssuesEvent) error {
	issue := ev.GetIssue()

	ch := c.discord.FindThreadByNumber(issue.GetNumber())
	if ch == nil {
		return fmt.Errorf("issue %d does not have a thread", issue.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makeIssueDeletedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c IssuesClient) EmbedTransferredMsg(ev *github.IssuesEvent) error {
	issue := ev.GetIssue()

	ch := c.discord.FindThreadByNumber(issue.GetNumber())
	if ch == nil {
		return fmt.Errorf("issue %d does not have a thread", issue.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makeIssueTransferredEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c IssuesClient) EmbedAssignedMsg(ev *github.IssuesEvent) error {
	issue := ev.GetIssue()

	ch := c.discord.FindThreadByNumber(issue.GetNumber())
	if ch == nil {
		return fmt.Errorf("issue %d does not have a thread", issue.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makeIssueAssignedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c IssuesClient) EmbedUnassignedMsg(ev *github.IssuesEvent) error {
	issue := ev.GetIssue()

	ch := c.discord.FindThreadByNumber(issue.GetNumber())
	if ch == nil {
		return fmt.Errorf("issue %d does not have a thread", issue.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makeIssueUnassignedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c IssuesClient) EmbedLabeledMsg(ev *github.IssuesEvent) error {
	issue := ev.GetIssue()

	ch := c.discord.FindThreadByNumber(issue.GetNumber())
	if ch == nil {
		return fmt.Errorf("issue %d does not have a thread", issue.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makeIssueLabeledEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c IssuesClient) EmbedUnlabeledMsg(ev *github.IssuesEvent) error {
	issue := ev.GetIssue()

	ch := c.discord.FindThreadByNumber(issue.GetNumber())
	if ch == nil {
		return fmt.Errorf("issue %d does not have a thread", issue.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makeIssueUnlabeledEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c IssuesClient) EmbedLockedMsg(ev *github.IssuesEvent) error {
	issue := ev.GetIssue()

	ch := c.discord.FindThreadByNumber(issue.GetNumber())
	if ch == nil {
		return fmt.Errorf("issue %d does not have a thread", issue.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makeIssueLockedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c IssuesClient) EmbedUnlockedMsg(ev *github.IssuesEvent) error {
	issue := ev.GetIssue()

	ch := c.discord.FindThreadByNumber(issue.GetNumber())
	if ch == nil {
		return fmt.Errorf("issue %d does not have a thread", issue.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makeIssueUnlockedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c IssuesClient) EmbedMilestonedMsg(ev *github.IssuesEvent) error {
	issue := ev.GetIssue()

	ch := c.discord.FindThreadByNumber(issue.GetNumber())
	if ch == nil {
		return fmt.Errorf("issue %d does not have a thread", issue.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makeIssueMilestonedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c IssuesClient) EmbedDemilestonedMsg(ev *github.IssuesEvent) error {
	issue := ev.GetIssue()

	ch := c.discord.FindThreadByNumber(issue.GetNumber())
	if ch == nil {
		return fmt.Errorf("issue %d does not have a thread", issue.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makeIssueDemilestonedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}
