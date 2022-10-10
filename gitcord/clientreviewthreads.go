package gitcord

import (
	"fmt"

	"github.com/google/go-github/v47/github"
	"github.com/pkg/errors"
)

type ReviewThreadsClient client

func (c ReviewThreadsClient) EmbedReviewThreadMsg(ev *github.PullRequestReviewThreadEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRReviewThreadEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c ReviewThreadsClient) EditReviewThreadMsg(ev *github.PullRequestReviewThreadEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	msg := c.discord.FindMsgByComment(ch, ev.GetThread().GetID())
	if msg == nil {
		return fmt.Errorf("failed to find message")
	}

	_, err := c.discord.EditEmbeds(ch.ID, msg.ID, c.config.makePRReviewThreadEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to edit message")
	}

	return nil
}

func (c ReviewThreadsClient) EmbedReviewThreadResolvedMsg(ev *github.PullRequestReviewThreadEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRReviewThreadResolvedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c ReviewThreadsClient) EmbedReviewThreadUnresolvedMsg(ev *github.PullRequestReviewThreadEvent) error {
	pr := ev.GetPullRequest()

	ch := c.discord.FindThreadByNumber(pr.GetNumber())
	if ch == nil {
		return fmt.Errorf("pull request %d does not have a thread", pr.GetNumber())
	}

	_, err := c.discord.SendEmbeds(ch.ID, c.config.makePRReviewThreadUnresolvedEmbed(ev))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}
