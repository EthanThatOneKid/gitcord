package gitcord

import (
	"context"
	"fmt"

	"github.com/google/go-github/v47/github"
	"github.com/pkg/errors"
)

type ReviewCommentsClient struct {
	*client
	ctx context.Context
}

func (c *ReviewCommentsClient) WithContext(ctx context.Context) *ReviewCommentsClient {
	cpy := *c
	cpy.client = cpy.client.WithContext(ctx)
	cpy.ctx = ctx
	return &cpy
}

func (c ReviewCommentsClient) EmbedReviewCommentMsg(ev *github.PullRequestReviewCommentEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRReviewCommentEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c ReviewCommentsClient) EditReviewCommentMsg(ev *github.PullRequestReviewCommentEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	msg := c.discord.FindMsgByComment(ch, ev.GetComment().GetID())
	if msg == nil {
		return fmt.Errorf("failed to find message")
	}

	_, err := c.discord.EditEmbeds(ch.ID, msg.ID, c.config.makePRReviewCommentEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to edit message")
	}

	return nil
}

func (c ReviewCommentsClient) EmbedReviewCommentDeletedMsg(ev *github.PullRequestReviewCommentEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRReviewCommentDeletedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}
