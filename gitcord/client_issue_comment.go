package gitcord

import (
	"context"

	"github.com/pkg/errors"
)

type IssueCommentClient struct {
	*client
	ctx context.Context
}

func (c *IssueCommentClient) WithContext(ctx context.Context) *IssueCommentClient {
	cpy := *c
	cpy.client = cpy.client.WithContext(ctx)
	cpy.ctx = ctx
	return &cpy
}

func (c *IssueCommentClient) ForwardCreated(issueID int, commentID int64) error {
	owner, repo := c.config.SplitGitHubRepo()

	comment, _, err := c.github.Issues.GetComment(c.ctx, owner, repo, commentID)
	if err != nil {
		return errors.Wrap(err, "failed to get comment")
	}

	issue, _, err := c.github.Issues.Get(c.ctx, owner, repo, issueID)
	if err != nil {
		return errors.Wrap(err, "failed to get issue")
	}

	ch := c.discord.FindChannelByIssue(issueID)
	if ch == nil {
		c.logln("skipping unknown issue with ID", issueID)
		return nil
	}

	_, err = c.discord.SendEmbeds(ch.ID, embedIssueComment(c.config, issue, comment))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c *IssueCommentClient) ForwardEdited(issueID int, commentID int64) error {
	owner, repo := c.config.SplitGitHubRepo()

	comment, _, err := c.github.Issues.GetComment(c.ctx, owner, repo, commentID)
	if err != nil {
		return errors.Wrap(err, "failed to get comment")
	}

	issue, _, err := c.github.Issues.Get(c.ctx, owner, repo, issueID)
	if err != nil {
		return errors.Wrap(err, "failed to get issue")
	}

	ch := c.discord.FindChannelByIssue(issueID)
	if ch == nil {
		c.logln("skipping unknown issue with ID", issueID)
		return nil
	}

	msg := c.discord.FindMsgByID(ch, commentID)
	if msg == nil {
		c.logln("skipping unknown comment with ID", commentID)
		return nil
	}

	_, err = c.discord.EditEmbeds(ch.ID, msg.ID, embedIssueComment(c.config, issue, comment))
	if err != nil {
		return errors.Wrap(err, "failed to edit message")
	}

	return nil
}
