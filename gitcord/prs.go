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

func (c *PRsClient) OpenThread(ev *github.PullRequestEvent) error {
	pr := ev.GetPullRequest()

	if pr.GetState() != "open" {
		if !c.config.ForceCreateThread {
			return fmt.Errorf("pull request %d is not open", pr.GetNumber())
		}
		c.logln("ignoring closed pull request", pr.GetNumber())
	}

	ch := c.discord.FindThreadByIssue(pr.GetNumber())
	if ch != nil {
		if ch.Type != discord.GuildPublicThread {
			if !c.config.ForceCreateThread {
				return fmt.Errorf("channel %d is not a public thread", ch.ID)
			}
			c.logln("ignoring channel %d is not a public thread", ch.ID)
		}

		if !c.config.ForceCreateThread {
			return fmt.Errorf("pull request %d already has a thread %d", pr.GetNumber(), ch.ID)
		}
		c.logln("ignoring existing thread %d", ch.ID)
	}

	ch, err := c.discord.StartThreadWithoutMessage(c.config.DiscordChannelID, api.StartThreadData{
		Name:                fmt.Sprintf("#%d: %s", pr.GetNumber(), pr.GetTitle()),
		Type:                discord.GuildPublicThread,
		AutoArchiveDuration: discord.SevenDaysArchive,
	})
	if err != nil {
		return errors.Wrap(err, "failed to open thread")
	}

	_, err = c.discord.SendEmbeds(ch.ID, *c.config.makePREmbed(pr))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (c *PRsClient) ForwardClosed(issueID int, commentID int64) error {
	return errors.New("not implemented")
}

func (c *PRsClient) ForwardReopened(issueID int) error {
	return errors.New("not implemented")
}

// func (c *PRsClient) ForwardReviewComment(prID int, commentID int64) error {
// 	owner, repo := c.config.SplitGitHubRepo()

// 	comment, _, err := c.github.PullRequests.GetComment(c.ctx, owner, repo, commentID)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to get comment")
// 	}

// 	pr, _, err := c.github.PullRequests.Get(c.ctx, owner, repo, prID)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to get PR")
// 	}

// 	chs, err := c.activeThreads()
// 	if err != nil {
// 		return errors.Wrap(err, "failed to get active threads")
// 	}

// 	ch := findChannelByID(chs, prID)
// 	if ch == nil {
// 		c.logln("skipping unknown PR with ID", prID)
// 		return nil
// 	}

// 	_, err = c.discord.SendEmbeds(ch.ID, embedPRReview(c.config, pr, comment))
// 	if err != nil {
// 		return errors.Wrap(err, "failed to send message")
// 	}

// 	return nil
// }

func (c *PRsClient) ForwardDeleted(prID int) error {
	return errors.New("not implemented")
}

func (c *PRsClient) ForwardMerged(prID int) error {
	return errors.New("not implemented")
}

func (c *PRsClient) ForwardUnmerged(prID int) error {
	return errors.New("not implemented")
}

func (c *PRsClient) ForwardEdited(prID int) error {
	return errors.New("not implemented")
}

func (c *PRsClient) ForwardAssigned(prID int) error {
	return errors.New("not implemented")
}

func (c *PRsClient) ForwardUnassigned(prID int) error {
	return errors.New("not implemented")
}

func (c *PRsClient) ForwardLabeled(prID int) error {
	return errors.New("not implemented")
}

func (c *PRsClient) ForwardUnlabeled(prID int) error {
	return errors.New("not implemented")
}

func (c *PRsClient) ForwardReviewRequested(prID int) error {
	return errors.New("not implemented")
}

func (c *PRsClient) ForwardReviewRequestRemoved(prID int) error {
	return errors.New("not implemented")
}

func (c *PRsClient) ForwardReadyForReview(prID int) error {
	return errors.New("not implemented")
}

func (c *PRsClient) ForwardConvertedToDraft(prID int) error {
	return errors.New("not implemented")
}

func (c *PRsClient) ForwardAutoMergeEnabled(prID int) error {
	return errors.New("not implemented")
}

func (c *PRsClient) ForwardAutoMergeDisabled(prID int) error {
	return errors.New("not implemented")
}
