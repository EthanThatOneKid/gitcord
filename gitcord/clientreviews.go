package gitcord

import (
	"fmt"

	"github.com/google/go-github/v47/github"
	"github.com/pkg/errors"
)

type ReviewsClient client

func (c *ReviewsClient) EmbedReviewMsg(ev *github.PullRequestReviewEvent) error {
	pr := ev.GetPullRequest()

	ch, err := c.discord.FindThreadByNumber(pr.GetNumber())
	if err != nil {
		return &PRThreadError{PR: pr.GetNumber(), Err: err}
	}

	_, err = c.discord.SendEmbeds(ch.ID, c.config.makePRReviewEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c ReviewsClient) EmbedReviewDismissedMsg(ev *github.PullRequestReviewEvent) error {
	pr := ev.GetPullRequest()

	ch, err := c.discord.FindThreadByNumber(pr.GetNumber())
	if err != nil {
		return &PRThreadError{PR: pr.GetNumber(), Err: err}
	}

	_, err = c.discord.SendEmbeds(ch.ID, c.config.makePRReviewDismissedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c ReviewsClient) EditReviewMsg(ev *github.PullRequestReviewEvent) error {
	pr := ev.GetPullRequest()

	ch, err := c.discord.FindThreadByNumber(pr.GetNumber())
	if err != nil {
		return &PRThreadError{PR: pr.GetNumber(), Err: err}
	}

	msg := c.discord.FindMsgByComment(ch, ev.GetReview().GetID())
	if msg == nil {
		return fmt.Errorf("failed to find message")
	}

	_, err = c.discord.EditEmbeds(ch.ID, msg.ID, c.config.makePRReviewEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to edit message")
	}

	return nil
}
