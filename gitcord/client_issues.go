package gitcord

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
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

func (c *IssuesClient) OpenChannel(issueID int) error {
	owner, repo := c.config.SplitGitHubRepo()

	issue, _, err := c.github.Issues.Get(c.ctx, owner, repo, issueID)
	if err != nil {
		return errors.Wrap(err, "failed to get issue")
	}

	ch, err := c.discord.StartThreadWithoutMessage(c.config.DiscordChannelID, api.StartThreadData{
		Name:                fmt.Sprintf("#%d: %s", issue.GetNumber(), issue.GetTitle()),
		Type:                discord.GuildPublicThread,
		AutoArchiveDuration: discord.OneDayArchive,
	})
	if err != nil {
		return errors.Wrap(err, "failed to open channel")
	}

	_, err = c.discord.SendEmbeds(ch.ID, embedOpeningIssueChannelMsg(c.config, issue))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c *IssuesClient) ForwardClosed(issueID int, commentID int64) error {
	return errors.New("not implemented")
}

func (c *IssuesClient) ForwardReopened(issueID int) error {
	return errors.New("not implemented")
}

func (c *IssuesClient) ForwardEdited(issueID int) error {
	return errors.New("not implemented")
}

func (c *IssuesClient) ForwardCreated(issueID int) error {
	return errors.New("not implemented")
}

func (c *IssuesClient) ForwardDeleted(issueID int) error {
	return errors.New("not implemented")
}

func (c *IssuesClient) ForwardTransferred(issueID int) error {
	return errors.New("not implemented")
}

func (c *IssuesClient) ForwardAssigned(issueID int) error {
	return errors.New("not implemented")
}

func (c *IssuesClient) ForwardUnassigned(issueID int) error {
	return errors.New("not implemented")
}

func (c *IssuesClient) ForwardLabeled(issueID int) error {
	return errors.New("not implemented")
}

func (c *IssuesClient) ForwardUnlabeled(issueID int) error {
	return errors.New("not implemented")
}
