package gitcord

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/google/go-github/v47/github"
	"github.com/pkg/errors"
)

type IssuesClient struct {
	*client
	ctx context.Context
}

func (c *IssuesClient) WithContext(ctx context.Context) *IssuesClient {
	cpy := *c
	cpy.client = cpy.client.WithContext(ctx)
	cpy.ctx = ctx
	return &cpy
}

func (c *IssuesClient) OpenThread(ev *github.IssuesEvent) error {
	issue := ev.GetIssue()

	if issue.GetState() != "open" {
		if !c.config.ForceCreateThread {
			return fmt.Errorf("issue %d is not open", issue.GetNumber())
		}
		c.logln("ignoring closed issue", issue.GetNumber())
	}

	ch := c.discord.FindThreadByIssue(issue.GetNumber())
	if ch != nil {
		if ch.Type != discord.GuildPublicThread {
			if !c.config.ForceCreateThread {
				return fmt.Errorf("channel %d is not a public thread", ch.ID)
			}
			c.logln("ignoring channel %d is not a public thread", ch.ID)
		}

		if !c.config.ForceCreateThread {
			return fmt.Errorf("issue %d already has a thread %d", issue.GetNumber(), ch.ID)
		}
		c.logln("ignoring existing thread %d", ch.ID)
	}

	ch, err := c.discord.StartThreadWithoutMessage(c.config.DiscordChannelID, api.StartThreadData{
		Name:                fmt.Sprintf("#%d: %s", issue.GetNumber(), issue.GetTitle()),
		Type:                discord.GuildPublicThread,
		AutoArchiveDuration: discord.SevenDaysArchive,
	})
	if err != nil {
		return errors.Wrap(err, "failed to open thread")
	}

	var embed = c.config.makeIssueEmbed(issue)

	_, err = c.discord.SendEmbeds(ch.ID, embed)
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c *IssuesClient) EmbedClosed(ev *github.IssuesEvent) error {
	issue := ev.GetIssue()

	ch := c.discord.FindThreadByIssue(issue.GetNumber())
	if ch == nil {
		return fmt.Errorf("issue %d does not have a thread", issue.GetNumber())
	}

	var embed = c.config.makeIssueClosedEmbed(ev)

	_, err := c.discord.SendEmbeds(ch.ID, embed)
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c *IssuesClient) SendReopened(issueID int) error {
	return errors.New("not implemented")
}

func (c *IssuesClient) SendEdited(issueID int) error {
	return errors.New("not implemented")
}

func (c *IssuesClient) SendCreated(issueID int) error {
	return errors.New("not implemented")
}

func (c *IssuesClient) SendDeleted(issueID int) error {
	return errors.New("not implemented")
}

func (c *IssuesClient) SendTransferred(issueID int) error {
	return errors.New("not implemented")
}

func (c *IssuesClient) SendAssigned(issueID int) error {
	return errors.New("not implemented")
}

func (c *IssuesClient) SendUnassigned(issueID int) error {
	return errors.New("not implemented")
}

func (c *IssuesClient) SendLabeled(issueID int) error {
	return errors.New("not implemented")
}

func (c *IssuesClient) SendUnlabeled(issueID int) error {
	return errors.New("not implemented")
}
