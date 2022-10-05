package gitcord

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/google/go-github/v47/github"
	smore "github.com/pythonian23/SMoRe"
)

/// START IssuesEvent Discord embeds

func (c Config) makeIssueEmbed(issue *github.Issue) discord.Embed {
	fields := &[]discord.EmbedField{
		{
			Name:   "Status",
			Value:  issue.GetState(),
			Inline: true,
		},
	}

	if len(issue.Labels) > 0 {
		*fields = append(*fields, discord.EmbedField{
			Name:   "Labels",
			Value:  convertLabels(issue.Labels),
			Inline: true,
		})
	}

	if len(issue.Assignees) > 0 {
		*fields = append(*fields, discord.EmbedField{
			Name:   "Assignees",
			Value:  convertUsers(issue.Assignees),
			Inline: true,
		})
	}

	if issue.Milestone != nil {
		*fields = append(*fields, discord.EmbedField{
			Name:   "Milestone",
			Value:  issue.Milestone.GetTitle(),
			Inline: true,
		})
	}

	if issue.GetLocked() {
		*fields = append(*fields, discord.EmbedField{
			Name:   "Locked",
			Value:  "🔒",
			Inline: true,
		})
	}

	return discord.Embed{
		Title: fmt.Sprintf("Issue opened: #%d %s", issue.GetNumber(), issue.GetTitle()),
		URL:   issue.GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			Name: issue.GetUser().GetLogin(),
			Icon: issue.GetUser().GetAvatarURL(),
		},
		Description: convertMarkdown(issue.GetBody()),
		Color:       discord.Color(c.Colors.IssueOpened.Success),
		Fields:      *fields,
	}
}

func (c Config) makeIssueClosedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d closed as completed", ev.GetIssue().GetNumber()),
		Color: discord.Color(c.Colors.IssueClosed.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makeIssueReopenedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d reopened", ev.GetIssue().GetNumber()),
		Color: discord.Color(c.Colors.IssueReopened.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makeIssueLabeledEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		// TODO: Add label color
		Title: fmt.Sprintf("Issue #%d: added %s label", ev.GetIssue().GetNumber(), ev.GetLabel().GetName()),
		Color: discord.Color(c.Colors.IssueLabeled.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makeIssueUnlabeledEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		// TODO: Add label color
		Title: fmt.Sprintf("Issue #%d: removed %s label", ev.GetIssue().GetNumber(), ev.GetLabel().GetName()),
		Color: discord.Color(c.Colors.IssueUnlabeled.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makeIssueAssignedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d assigned to %s", ev.GetIssue().GetNumber(), ev.GetAssignee().GetLogin()),
		Color: discord.Color(c.Colors.IssueAssigned.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makeIssueUnassignedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d unassigned to %s", ev.GetIssue().GetNumber(), ev.GetAssignee().GetLogin()),
		Color: discord.Color(c.Colors.IssueUnassigned.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makeIssueMilestonedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		// TODO: Add milestone title
		Title: fmt.Sprintf("Issue #%d: milestone added", ev.GetIssue().GetNumber()),
		Color: discord.Color(c.Colors.IssueMilestoned.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makeIssueDemilestonedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d: milestone removed", ev.GetIssue().GetNumber()),
		Color: discord.Color(c.Colors.IssueDemilestoned.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makeIssueDeletedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d deleted", ev.GetIssue().GetNumber()),
		Color: discord.Color(c.Colors.IssueDeleted.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makeIssueLockedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d locked", ev.GetIssue().GetNumber()),
		Color: discord.Color(c.Colors.IssueLocked.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makeIssueUnlockedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d unlocked", ev.GetIssue().GetNumber()),
		Color: discord.Color(c.Colors.IssueUnlocked.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makeIssueTransferredEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Issue #%d transferred to %s", ev.GetIssue().GetNumber(), ev.GetRepo().GetFullName()),
		Color: discord.Color(c.Colors.IssueTransferred.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

/// END IssuesEvent Discord embeds
/// START IssueCommentEvent Discord embeds

func (c Config) makeIssueCommentEmbed(ev *github.IssueCommentEvent) discord.Embed {
	issue, comment := ev.GetIssue(), ev.GetComment()

	var fields *[]discord.EmbedField
	for react, count := range parseReactions(comment.Reactions) {
		*fields = append(*fields, discord.EmbedField{
			Name:   react,
			Value:  fmt.Sprintf("%d", count),
			Inline: true,
		})
	}

	title := fmt.Sprintf("Comment on issue #%d", issue.GetNumber())
	if checkPR(issue) {
		title = fmt.Sprintf("Comment on pull request #%d", issue.GetNumber())
	}

	return discord.Embed{
		Title:       title,
		Description: convertMarkdown(comment.GetBody()),
		URL:         comment.GetHTMLURL(),
		Color:       discord.Color(c.Colors.IssueCommented.Success),
		Fields:      *fields,
		Author: &discord.EmbedAuthor{
			Name: comment.GetUser().GetLogin(),
			Icon: comment.GetUser().GetAvatarURL(),
		},
		// Footer is used to store the comment ID
		Footer: &discord.EmbedFooter{Text: "0x" + strconv.FormatInt(comment.GetID(), 16)},
	}
}

func (c Config) makeIssueCommentDeletedEmbed(ev *github.IssueCommentEvent) discord.Embed {
	return discord.Embed{
		Title:       fmt.Sprintf("Deleted comment on issue #%d", ev.GetIssue().GetNumber()),
		Description: fmt.Sprintf("Comment ID: %s", strconv.FormatInt(*ev.GetComment().ID, 16)),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
		Color: discord.Color(c.Colors.IssueCommentDeleted.Success),
	}
}

/// END IssueCommentEvent Discord embeds
/// START PullRequestEvent Discord embeds

func (c Config) makePREmbed(ev *github.PullRequestEvent) discord.Embed {
	pr := ev.GetPullRequest()

	var fields *[]discord.EmbedField

	if len(pr.Labels) > 0 {
		*fields = append(*fields, discord.EmbedField{
			Name:   "Labels",
			Value:  convertLabels(pr.Labels),
			Inline: true,
		})
	}

	if len(pr.Assignees) > 0 {
		*fields = append(*fields, discord.EmbedField{
			Name:   "Assignees",
			Value:  convertUsers(pr.Assignees),
			Inline: true,
		})
	}

	if len(pr.RequestedReviewers) > 0 {
		*fields = append(*fields, discord.EmbedField{
			Name:   "Requested reviewers",
			Value:  convertUsers(pr.RequestedReviewers),
			Inline: true,
		})
	}

	if len(pr.RequestedTeams) > 0 {
		*fields = append(*fields, discord.EmbedField{
			Name:   "Requested teams",
			Value:  convertTeams(pr.RequestedTeams),
			Inline: true,
		})
	}

	if pr.Milestone != nil {
		*fields = append(*fields, discord.EmbedField{
			Name: "Milestone",
			// TODO: Provide link to milestone
			Value:  pr.Milestone.GetTitle(),
			Inline: true,
		})
	}

	return discord.Embed{
		Title: fmt.Sprintf("Pull request opened: #%d %s", pr.GetNumber(), pr.GetTitle()),
		URL:   pr.GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			Name: pr.GetUser().GetLogin(),
			Icon: pr.GetUser().GetAvatarURL(),
		},
		Description: convertMarkdown(pr.GetBody()),
		Color:       discord.Color(c.Colors.IssueOpened.Success),
		Fields:      *fields,
	}
}

func (c Config) makePRClosedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Pull request #%d closed", ev.GetPullRequest().GetNumber()),
		Color: discord.Color(c.Colors.PRClosed.Success),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makePRReopenedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Pull request #%d reopened", ev.GetPullRequest().GetNumber()),
		Color: discord.Color(c.Colors.PRReopened.Success),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makePRAssignedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Pull request #%d assigned to %s", ev.GetPullRequest().GetNumber(), ev.GetAssignee().GetLogin()),
		Color: discord.Color(c.Colors.PRAssigned.Success),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makePRUnassignedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Pull request #%d unassigned from %s", ev.GetPullRequest().GetNumber(), ev.GetAssignee().GetLogin()),
		Color: discord.Color(c.Colors.PRUnassigned.Success),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makePRDeletedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Deleted pull request #%d", ev.GetPullRequest().GetNumber()),
		Color: discord.Color(c.Colors.PRDeleted.Success),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makePRTransferredEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Transferred pull request #%d", ev.GetPullRequest().GetNumber()),
		Color: discord.Color(c.Colors.PRTransferred.Success),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makePRLabeledEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Labeled pull request #%d", ev.GetPullRequest().GetNumber()),
		Color: discord.Color(c.Colors.PRLabeled.Success),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makePRUnlabeledEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Unlabeled pull request #%d", ev.GetPullRequest().GetNumber()),
		Color: discord.Color(c.Colors.PRUnlabeled.Success),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makePRMilestonedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Milestoned pull request #%d", ev.GetPullRequest().GetNumber()),
		Color: discord.Color(c.Colors.PRMilestoned.Success),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makePRDemilestonedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Demilestoned pull request #%d", ev.GetPullRequest().GetNumber()),
		Color: discord.Color(c.Colors.PRDemilestoned.Success),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makePRLockedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Locked pull request #%d", ev.GetPullRequest().GetNumber()),
		Color: discord.Color(c.Colors.PRLocked.Success),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makePRUnlockedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Unlocked pull request #%d", ev.GetPullRequest().GetNumber()),
		Color: discord.Color(c.Colors.PRUnlocked.Success),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makePRReviewRequestedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Review requested for pull request #%d", ev.GetPullRequest().GetNumber()),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Color: discord.Color(c.Colors.PRReviewRequested.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makePRReviewRequestRemovedEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Review request removed for pull request #%d", ev.GetPullRequest().GetNumber()),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Color: discord.Color(c.Colors.PRReviewRequestRemoved.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makePRReadyForReviewEmbed(ev *github.PullRequestEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Pull request #%d is ready for review", ev.GetPullRequest().GetNumber()),
		Color: discord.Color(c.Colors.PRReadyForReview.Success),
		URL:   ev.GetPullRequest().GetHTMLURL(),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

/// END PullRequestEvent Discord embeds
/// START PullRequestReviewEvent Discord embeds

func (c Config) makePRReviewEmbed(ev *github.PullRequestReviewEvent) discord.Embed {
	pr, review := ev.GetPullRequest(), ev.GetReview()

	return discord.Embed{
		Title:       fmt.Sprintf("Review submitted on pull request #%d", pr.GetNumber()),
		Description: convertMarkdown(review.GetBody()),
		URL:         review.GetHTMLURL(),
		Color:       discord.Color(c.Colors.Reviewed.Success),
		Author: &discord.EmbedAuthor{
			Name: review.GetUser().GetLogin(),
			Icon: review.GetUser().GetAvatarURL(),
		},
		// Footer is used to store the review ID, similar to makeIssueCommentEmbed
		Footer: &discord.EmbedFooter{Text: "0x" + strconv.FormatInt(review.GetID(), 16)},
	}
}

func (c Config) makePRReviewDismissedEmbed(ev *github.PullRequestReviewEvent) discord.Embed {
	pr, review := ev.GetPullRequest(), ev.GetReview()

	return discord.Embed{
		Title:       fmt.Sprintf("Review dismissed on pull request #%d", pr.GetNumber()),
		Description: convertMarkdown(review.GetBody()),
		URL:         review.GetHTMLURL(),
		Color:       discord.Color(c.Colors.ReviewDismissed.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

/// END PullRequestReviewEvent Discord embeds
/// START PullRequestReviewCommentEvent Discord embeds

func (c Config) makePRReviewCommentEmbed(ev *github.PullRequestReviewCommentEvent) discord.Embed {
	pr, comment := ev.GetPullRequest(), ev.GetComment()

	var fields *[]discord.EmbedField
	for react, count := range parseReactions(comment.Reactions) {
		*fields = append(*fields, discord.EmbedField{
			Name:   react,
			Value:  fmt.Sprintf("%d", count),
			Inline: true,
		})
	}

	return discord.Embed{
		Title:       fmt.Sprintf("Review comment on pull request #%d", pr.GetNumber()),
		Description: convertMarkdown(comment.GetBody()),
		URL:         comment.GetHTMLURL(),
		Color:       discord.Color(c.Colors.ReviewCommented.Success),
		Fields:      *fields,
		Author: &discord.EmbedAuthor{
			Name: comment.GetUser().GetLogin(),
			Icon: comment.GetUser().GetAvatarURL(),
		},
		// Footer is used to store the comment ID, similar to makeIssueCommentEmbed
		Footer: &discord.EmbedFooter{Text: "0x" + strconv.FormatInt(comment.GetID(), 16)},
	}
}

func (c Config) makePRReviewCommentDeletedEmbed(ev *github.PullRequestReviewCommentEvent) discord.Embed {
	pr, comment := ev.GetPullRequest(), ev.GetComment()

	return discord.Embed{
		Title: fmt.Sprintf("Review comment deleted on pull request #%d", pr.GetNumber()),
		URL:   comment.GetHTMLURL(),
		Color: discord.Color(c.Colors.ReviewCommentDeleted.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

/// END PullRequestReviewCommentEvent Discord embeds
/// START PullRequestReviewThreadEvent Discord embeds

func (c Config) makePRReviewThreadEmbed(ev *github.PullRequestReviewThreadEvent) discord.Embed {
	pr, t := ev.GetPullRequest(), ev.GetThread()

	var fields *[]discord.EmbedField
	for _, comment := range t.Comments {
		*fields = append(*fields, discord.EmbedField{
			Name:   fmt.Sprintf("Review comment %d", comment.GetID()),
			Value:  convertMarkdown(comment.GetBody()),
			Inline: false,
		})
	}

	return discord.Embed{
		Title:  fmt.Sprintf("Review thread on pull request #%d", pr.GetNumber()),
		URL:    threadURL(t),
		Fields: *fields,
		Color:  discord.Color(c.Colors.ReviewThreaded.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makePRReviewThreadResolvedEmbed(ev *github.PullRequestReviewThreadEvent) discord.Embed {
	pr, t := ev.GetPullRequest(), ev.GetThread()

	return discord.Embed{
		Title: fmt.Sprintf("Review thread resolved on pull request #%d", pr.GetNumber()),
		URL:   threadURL(t),
		Color: discord.Color(c.Colors.ReviewThreadResolved.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

func (c Config) makePRReviewThreadUnresolvedEmbed(ev *github.PullRequestReviewThreadEvent) discord.Embed {
	pr, t := ev.GetPullRequest(), ev.GetThread()

	return discord.Embed{
		Title: fmt.Sprintf("Review thread unresolved on pull request #%d", pr.GetNumber()),
		URL:   threadURL(t),
		Color: discord.Color(c.Colors.ReviewThreadUnresolved.Success),
		Author: &discord.EmbedAuthor{
			Name: ev.GetSender().GetLogin(),
			Icon: ev.GetSender().GetAvatarURL(),
		},
	}
}

/// END PullRequestReviewThreadEvent Discord embeds

func convertSlicePtr[T any](s []*T) []T {
	r := make([]T, 0, len(s))
	for _, v := range s {
		r = append(r, *v)
	}
	return r
}

func convertLabels[T *github.Label | github.Label](labelsV []T) string {
	var labels []github.Label
	switch labelsV := any(labelsV).(type) {
	case []github.Label:
		labels = labelsV
	case []*github.Label:
		labels = convertSlicePtr(labelsV)
	}

	var labelNames []string
	for _, label := range labels {
		labelNames = append(labelNames, label.GetName())
	}

	return strings.Join(labelNames, ", ")
}

func threadURL(t *github.PullRequestThread) (url string) {
	comment := t.Comments[0]
	if comment != nil {
		url = comment.GetHTMLURL()
	}
	return
}

func convertUsers(users []*github.User) string {
	var names []string

	for _, user := range users {
		names = append(names, user.GetLogin())
	}

	return strings.Join(names, ", ")
}

func convertTeams(teams []*github.Team) string {
	var names []string

	for _, team := range teams {
		names = append(names, team.GetName())
	}

	return strings.Join(names, ", ")
}

// see https://github.com/pythonian23/SMoRe
func convertMarkdown(githubMD string) string {
	return smore.Render(githubMD)
}

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
