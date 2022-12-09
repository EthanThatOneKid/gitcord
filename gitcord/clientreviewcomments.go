package gitcord

import (
	"fmt"

	"github.com/google/go-github/v47/github"
	"github.com/pkg/errors"
)

type ReviewCommentsClient client

func (c ReviewCommentsClient) EmbedReviewCommentMsg(ev *github.PullRequestReviewCommentEvent) error {
	pr := ev.GetPullRequest()

	ch, err := c.discord.FindThreadByNumber(pr.GetNumber())
	if err != nil {
		return &PRThreadError{PR: pr.GetNumber(), Err: err}
	}

	_, err = c.discord.SendEmbeds(ch.ID, c.config.makePRReviewCommentEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c ReviewCommentsClient) EditReviewCommentMsg(ev *github.PullRequestReviewCommentEvent) error {
	pr := ev.GetPullRequest()

	ch, err := c.discord.FindThreadByNumber(pr.GetNumber())
	if err != nil {
		return &PRThreadError{PR: pr.GetNumber(), Err: err}
	}

	msg := c.discord.FindMsgByComment(ch, ev.GetComment().GetID())
	if msg == nil {
		return fmt.Errorf("failed to find message")
	}

	_, err = c.discord.EditEmbeds(ch.ID, msg.ID, c.config.makePRReviewCommentEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to edit message")
	}

	return nil
}

func (c ReviewCommentsClient) EmbedReviewCommentDeletedMsg(ev *github.PullRequestReviewCommentEvent) error {
	pr := ev.GetPullRequest()

	ch, err := c.discord.FindThreadByNumber(pr.GetNumber())
	if err != nil {
		return &PRThreadError{PR: pr.GetNumber(), Err: err}
	}

	_, err = c.discord.SendEmbeds(ch.ID, c.config.makePRReviewCommentDeletedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}
