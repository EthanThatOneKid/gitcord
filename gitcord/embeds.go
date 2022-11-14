package gitcord

import (
	"fmt"
	"strconv"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/ethanthatonekid/gitcord/gitcord/internal/discordclient"
	"github.com/ethanthatonekid/gitcord/gitcord/internal/markdown"
	"github.com/google/go-github/v47/github"
)

/// START IssuesEvent Discord embeds

func (c *Config) makeIssueEmbed(issue *github.Issue) discord.Embed {
	fields := []discord.EmbedField{
		{
			Name:  "Status",
			Value: issue.GetState(),
		},
	}

	if len(issue.Labels) > 0 {
		fields = append(fields, discord.EmbedField{
			Name:  "Labels",
			Value: markdown.ConvertLabels(issue.Labels),
		})
	}

	if len(issue.Assignees) > 0 {
		fields = append(fields, discord.EmbedField{
			Name:  "Assignees",
			Value: markdown.ConvertUsers(issue.Assignees),
		})
	}

	if issue.Milestone != nil {
		fields = append(fields, discord.EmbedField{
			Name:  "Milestone",
			Value: markdown.ConvertHyperlink(issue.Milestone.GetTitle(), issue.Milestone.GetHTMLURL()),
		})
	}

	if issue.GetLocked() {
		fields = append(fields, discord.EmbedField{
			Name:  "Locked",
			Value: "🔒",
		})
	}

	return discord.Embed{
		Title: discordclient.IssueMsgPrefix + fmt.Sprintf("%d %s", issue.GetNumber(), issue.GetTitle()),
		URL:   issue.GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			URL:  issue.GetUser().GetHTMLURL(),
			Name: issue.GetUser().GetLogin(),
			Icon: issue.GetUser().GetAvatarURL(),
		},
		Description: markdown.ConvertDiscord(issue.GetBody(), issue.GetHTMLURL()),
		Color:       c.ColorScheme.Color(IssueOpened, true),
		Fields:      fields,
	}
}

func (c *Config) makeIssueClosedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d closed as completed", ev.GetIssue().GetNumber()),
		Color: c.ColorScheme.Color(IssueClosed, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makeIssueReopenedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d reopened", ev.GetIssue().GetNumber()),
		Color: c.ColorScheme.Color(IssueReopened, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makeIssueLabeledEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		// TODO: Add label color
		Title: fmt.Sprintf("Issue #%d: added %s label", ev.GetIssue().GetNumber(), ev.GetLabel().GetName()),
		Color: c.ColorScheme.Color(IssueLabeled, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makeIssueUnlabeledEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		// TODO: Add label color
		Title: fmt.Sprintf("Issue #%d: removed %s label", ev.GetIssue().GetNumber(), ev.GetLabel().GetName()),
		Color: c.ColorScheme.Color(IssueUnlabeled, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makeIssueAssignedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d assigned to %s", ev.GetIssue().GetNumber(), ev.GetAssignee().GetLogin()),
		Color: c.ColorScheme.Color(IssueAssigned, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makeIssueUnassignedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d unassigned to %s", ev.GetIssue().GetNumber(), ev.GetAssignee().GetLogin()),
		Color: c.ColorScheme.Color(IssueUnassigned, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makeIssueMilestonedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		// TODO: Add milestone title
		Title: fmt.Sprintf("Issue #%d: %s milestone added", ev.GetIssue().GetNumber(), ev.GetIssue().GetMilestone().GetTitle()),
		Color: c.ColorScheme.Color(IssueMilestoned, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makeIssueDemilestonedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d: %s milestone removed", ev.GetIssue().GetNumber(), ev.GetIssue().GetMilestone().GetTitle()),
		Color: c.ColorScheme.Color(IssueDemilestoned, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makeIssueDeletedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d deleted", ev.GetIssue().GetNumber()),
		Color: c.ColorScheme.Color(IssueDeleted, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makeIssueLockedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d locked", ev.GetIssue().GetNumber()),
		Color: c.ColorScheme.Color(IssueLocked, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makeIssueUnlockedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d unlocked", ev.GetIssue().GetNumber()),
		Color: c.ColorScheme.Color(IssueUnlocked, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makeIssueTransferredEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d transferred to %s", ev.GetIssue().GetNumber(), ev.GetRepo().GetFullName()),
		Color: c.ColorScheme.Color(IssueTransferred, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

/// END IssuesEvent Discord embeds
/// START IssueCommentEvent Discord embeds

func (c *Config) makeIssueCommentEmbed(ev *github.IssueCommentEvent) discord.Embed {
	issue, comment := ev.GetIssue(), ev.GetComment()

	var fields []discord.EmbedField
	for react, count := range parseReactions(comment.Reactions) {
		if count == 0 {
			continue
		}
		fields = append(fields, discord.EmbedField{
			Name:  react,
			Value: fmt.Sprintf("%d", count),
		})
	}

	title := fmt.Sprintf("Comment on issue #%d", issue.GetNumber())
	if checkPR(issue) {
		title = fmt.Sprintf("Comment on pull request #%d", issue.GetNumber())
	}

	return discord.Embed{
		Title:       title,
		Description: markdown.ConvertDiscord(comment.GetBody(), comment.GetHTMLURL()),
		URL:         comment.GetHTMLURL(),
		Color:       c.ColorScheme.Color(IssueCommented, true),
		Fields:      fields,
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: comment.GetUser().GetLogin(),
			Icon: comment.GetUser().GetAvatarURL(),
		},
		// Footer is used to store the comment ID
		Footer: &discord.EmbedFooter{Text: strconv.FormatInt(comment.GetID(), 10)},
	}
}

func (c *Config) makeIssueCommentDeletedEmbed(ev *github.IssueCommentEvent) discord.Embed {
	return discord.Embed{
		Title:       fmt.Sprintf("Deleted comment on issue #%d", ev.GetIssue().GetNumber()),
		Description: markdown.ConvertDiscord(ev.GetComment().GetBody(), ""),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
		Color:  c.ColorScheme.Color(IssueCommentDeleted, true),
		Footer: &discord.EmbedFooter{Text: strconv.FormatInt(ev.GetComment().GetID(), 10)},
	}
}

/// END IssueCommentEvent Discord embeds
/// START PullRequestEvent Discord embeds

func (c *Config) makePREmbed(ev *github.PullRequestEvent) discord.Embed {
	pr := ev.GetPullRequest()

	var fields []discord.EmbedField

	if len(pr.Labels) > 0 {
		fields = append(fields, discord.EmbedField{
			Name:  "Labels",
			Value: markdown.ConvertLabels(pr.Labels),
		})
	}

	if len(pr.Assignees) > 0 {
		fields = append(fields, discord.EmbedField{
			Name:  "Assignees",
			Value: markdown.ConvertUsers(pr.Assignees),
		})
	}

	if len(pr.RequestedReviewers) > 0 {
		fields = append(fields, discord.EmbedField{
			Name:  "Requested reviewers",
			Value: markdown.ConvertUsers(pr.RequestedReviewers),
		})
	}

	if len(pr.RequestedTeams) > 0 {
		fields = append(fields, discord.EmbedField{
			Name:  "Requested teams",
			Value: markdown.ConvertTeams(pr.RequestedTeams),
		})
	}

	if pr.Milestone != nil {
		fields = append(fields, discord.EmbedField{
			Name: "Milestone",
			// TODO: Provide link to milestone
			Value: pr.Milestone.GetTitle(),
		})
	}

	return discord.Embed{
		Title: discordclient.PRMsgPrefix + fmt.Sprintf("%d %s", pr.GetNumber(), pr.GetTitle()),
		URL:   pr.GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: pr.GetUser().GetLogin(),
			Icon: pr.GetUser().GetAvatarURL(),
		},
		Description: markdown.ConvertDiscord(pr.GetBody(), pr.GetHTMLURL()),
		Color:       c.ColorScheme.Color(IssueOpened, true),
		Fields:      fields,
	}
}

func (c *Config) makePRClosedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Pull request #%d closed", ev.GetPullRequest().GetNumber()),
		Color: c.ColorScheme.Color(PRClosed, true),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makePRReopenedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Pull request #%d reopened", ev.GetPullRequest().GetNumber()),
		Color: c.ColorScheme.Color(PRReopened, true),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makePRAssignedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Pull request #%d assigned to %s", ev.GetPullRequest().GetNumber(), ev.GetAssignee().GetLogin()),
		Color: c.ColorScheme.Color(PRAssigned, true),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makePRUnassignedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Pull request #%d unassigned from %s", ev.GetPullRequest().GetNumber(), ev.GetAssignee().GetLogin()),
		Color: c.ColorScheme.Color(PRUnassigned, true),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makePRDeletedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Deleted pull request #%d", ev.GetPullRequest().GetNumber()),
		Color: c.ColorScheme.Color(PRDeleted, true),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makePRTransferredEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Transferred pull request #%d", ev.GetPullRequest().GetNumber()),
		Color: c.ColorScheme.Color(PRTransferred, true),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makePRLabeledEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Labeled pull request #%d", ev.GetPullRequest().GetNumber()),
		Color: c.ColorScheme.Color(PRLabeled, true),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makePRUnlabeledEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Unlabeled pull request #%d", ev.GetPullRequest().GetNumber()),
		Color: c.ColorScheme.Color(PRUnlabeled, true),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makePRMilestonedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Milestoned pull request #%d", ev.GetPullRequest().GetNumber()),
		Color: c.ColorScheme.Color(PRMilestoned, true),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makePRDemilestonedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Demilestoned pull request #%d", ev.GetPullRequest().GetNumber()),
		Color: c.ColorScheme.Color(PRDemilestoned, true),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makePRLockedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Locked pull request #%d", ev.GetPullRequest().GetNumber()),
		Color: c.ColorScheme.Color(PRLocked, true),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makePRUnlockedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Unlocked pull request #%d", ev.GetPullRequest().GetNumber()),
		Color: c.ColorScheme.Color(PRUnlocked, true),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makePRReviewRequestedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Review requested for pull request #%d", ev.GetPullRequest().GetNumber()),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Color: c.ColorScheme.Color(PRReviewRequested, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makePRReviewRequestRemovedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Review request removed for pull request #%d", ev.GetPullRequest().GetNumber()),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Color: c.ColorScheme.Color(PRReviewRequestRemoved, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makePRReadyForReviewEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Pull request #%d is ready for review", ev.GetPullRequest().GetNumber()),
		Color: c.ColorScheme.Color(PRReadyForReview, true),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

/// END PullRequestEvent Discord embeds
/// START PullRequestReviewEvent Discord embeds

func (c *Config) makePRReviewEmbed(ev *github.PullRequestReviewEvent) discord.Embed {
	pr, review := ev.GetPullRequest(), ev.GetReview()

	return discord.Embed{
		Title:       fmt.Sprintf("Review submitted on pull request #%d", pr.GetNumber()),
		Description: markdown.ConvertDiscord(review.GetBody(), review.GetHTMLURL()),
		URL:         review.GetHTMLURL(),
		Color:       c.ColorScheme.Color(Reviewed, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: review.GetUser().GetLogin(),
			Icon: review.GetUser().GetAvatarURL(),
		},
		// Footer is used to store the review ID, similar to makeIssueCommentEmbed
		Footer: &discord.EmbedFooter{Text: strconv.FormatInt(review.GetID(), 10)},
	}
}

func (c *Config) makePRReviewDismissedEmbed(ev *github.PullRequestReviewEvent) discord.Embed {
	pr, review := ev.GetPullRequest(), ev.GetReview()

	return discord.Embed{
		Title:       fmt.Sprintf("Review dismissed on pull request #%d", pr.GetNumber()),
		Description: markdown.ConvertDiscord(review.GetBody(), review.GetHTMLURL()),
		URL:         review.GetHTMLURL(),
		Color:       c.ColorScheme.Color(ReviewDismissed, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

/// END PullRequestReviewEvent Discord embeds
/// START PullRequestReviewCommentEvent Discord embeds

func (c *Config) makePRReviewCommentEmbed(ev *github.PullRequestReviewCommentEvent) discord.Embed {
	pr, comment := ev.GetPullRequest(), ev.GetComment()

	var fields []discord.EmbedField
	for react, count := range parseReactions(comment.Reactions) {
		if count == 0 {
			continue
		}
		fields = append(fields, discord.EmbedField{
			Name:  react,
			Value: fmt.Sprintf("%d", count),
		})
	}

	return discord.Embed{
		Title:       fmt.Sprintf("Review comment on pull request #%d", pr.GetNumber()),
		Description: markdown.ConvertDiscord(comment.GetBody(), comment.GetHTMLURL()),
		URL:         comment.GetHTMLURL(),
		Color:       c.ColorScheme.Color(ReviewCommented, true),
		Fields:      fields,
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: comment.GetUser().GetLogin(),
			Icon: comment.GetUser().GetAvatarURL(),
		},
		// Footer is used to store the comment ID, similar to makeIssueCommentEmbed
		Footer: &discord.EmbedFooter{Text: strconv.FormatInt(comment.GetID(), 10)},
	}
}

func (c *Config) makePRReviewCommentDeletedEmbed(ev *github.PullRequestReviewCommentEvent) discord.Embed {
	pr, comment := ev.GetPullRequest(), ev.GetComment()

	return discord.Embed{
		Title: fmt.Sprintf("Review comment deleted on pull request #%d", pr.GetNumber()),
		URL:   comment.GetHTMLURL(),
		Color: c.ColorScheme.Color(ReviewCommentDeleted, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

/// END PullRequestReviewCommentEvent Discord embeds
/// START PullRequestReviewThreadEvent Discord embeds

func (c *Config) makePRReviewThreadEmbed(ev *github.PullRequestReviewThreadEvent) discord.Embed {
	pr, t := ev.GetPullRequest(), ev.GetThread()

	var fields []discord.EmbedField
	for _, comment := range t.Comments {
		fields = append(fields, discord.EmbedField{
			Name:   fmt.Sprintf("Review comment %d", comment.GetID()),
			Value:  markdown.ConvertDiscord(comment.GetBody(), comment.GetHTMLURL()),
			Inline: false,
		})
	}

	return discord.Embed{
		Title:  fmt.Sprintf("Review thread on pull request #%d", pr.GetNumber()),
		URL:    markdown.ConvertThreadURL(t),
		Fields: fields,
		Color:  c.ColorScheme.Color(ReviewThreaded, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makePRReviewThreadResolvedEmbed(ev *github.PullRequestReviewThreadEvent) discord.Embed {
	pr, t := ev.GetPullRequest(), ev.GetThread()

	return discord.Embed{
		Title: fmt.Sprintf("Review thread resolved on pull request #%d", pr.GetNumber()),
		URL:   markdown.ConvertThreadURL(t),
		Color: c.ColorScheme.Color(ReviewThreadResolved, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c *Config) makePRReviewThreadUnresolvedEmbed(ev *github.PullRequestReviewThreadEvent) discord.Embed {
	pr, t := ev.GetPullRequest(), ev.GetThread()

	return discord.Embed{
		Title: fmt.Sprintf("Review thread unresolved on pull request #%d", pr.GetNumber()),
		URL:   markdown.ConvertThreadURL(t),
		Color: c.ColorScheme.Color(ReviewThreadUnresolved, true),
		Author: &discord.EmbedAuthor{
			URL:  ev.GetSender().GetHTMLURL(),
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

/// END PullRequestReviewThreadEvent Discord embeds

// checkPR checks if an issue happens to be a pull request
func checkPR(issue *github.Issue) bool {
	return issue.GetPullRequestLinks() != nil
}

// parseReactions parses the reactions from a comment
func parseReactions(r *github.Reactions) map[string]int {
	reacts := make(map[string]int, 8)
	reacts["👍"] = r.GetPlusOne()
	reacts["👎"] = r.GetMinusOne()
	reacts["😆"] = r.GetLaugh()
	reacts["🎉"] = r.GetHooray()
	reacts["😕"] = r.GetConfused()
	reacts["❤️"] = r.GetHeart()
	reacts["🚀"] = r.GetRocket()
	reacts["👀"] = r.GetEyes()
	return reacts
}
