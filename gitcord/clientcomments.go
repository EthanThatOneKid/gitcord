package gitcord

import (
	"fmt"

	"github.com/google/go-github/v47/github"
	"github.com/pkg/errors"
)

type IssueCommentClient client

func (c *IssueCommentClient) EmbedIssueCommentMsg(ev *github.IssueCommentEvent) error {
	issue := ev.GetIssue()

	t, err := c.discord.FindThreadByNumber(issue.GetNumber())
	if err != nil {
		return fmt.Errorf("issue %d does not have a thread", issue.GetNumber())
	}

	_, err = c.discord.SendEmbeds(t.ID, c.config.makeIssueCommentEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c *IssueCommentClient) EditIssueCommentMsg(ev *github.IssueCommentEvent) error {
	issue := ev.GetIssue()

	t, err := c.discord.FindThreadByNumber(issue.GetNumber())
	if err != nil {
		return fmt.Errorf("issue %d does not have a thread", issue.GetNumber())
	}

	msg := c.discord.FindMsgByComment(t, ev.GetComment().GetID())
	if msg == nil {
		return fmt.Errorf("failed to find message")
	}

	_, err = c.discord.EditEmbeds(t.ID, msg.ID, c.config.makeIssueCommentEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to edit message")
	}

	return nil
}

func (c *IssueCommentClient) EmbedDeletedMsg(ev *github.IssueCommentEvent) error {
	issue := ev.GetIssue()

	t, err := c.discord.FindThreadByNumber(issue.GetNumber())
	if err != nil {
		return fmt.Errorf("issue %d does not have a thread", issue.GetNumber())
	}

	_, err = c.discord.SendEmbeds(t.ID, c.config.makeIssueCommentDeletedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}
