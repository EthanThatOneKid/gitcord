package gitcord

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/google/go-github/v47/github"
	smore "github.com/pythonian23/SMoRe"
)

// IssuesEvent embeds

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
		Title: fmt.Sprintf("%s closed this as completed", ev.GetSender().GetLogin()),
		Color: discord.Color(c.Colors.IssueClosed.Success),
	}
}

func (c Config) makeIssueReopenedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("%s reopened this", ev.GetSender().GetLogin()),
		Color: discord.Color(c.Colors.IssueReopened.Success),
	}
}

func (c Config) makeIssueLabeledEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("%s added label %s", ev.GetSender().GetLogin(), ev.GetLabel().GetName()),
		Color: discord.Color(c.Colors.IssueLabeled.Success),
	}
}

func (c Config) makeIssueUnlabeledEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("%s removed label %s", ev.GetSender().GetLogin(), ev.GetLabel().GetName()),
		Color: discord.Color(c.Colors.IssueUnlabeled.Success),
	}
}

func (c Config) makeIssueAssignedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("%s assigned %s", ev.GetSender().GetLogin(), ev.GetAssignee().GetLogin()),
		Color: discord.Color(c.Colors.IssueAssigned.Success),
	}
}

func (c Config) makeIssueUnassignedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("%s unassigned %s", ev.GetSender().GetLogin(), ev.GetAssignee().GetLogin()),
		Color: discord.Color(c.Colors.IssueUnassigned.Success),
	}
}

func (c Config) makeIssueMilestonedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("%s added milestone %s", ev.GetSender().GetLogin(), ev.GetMilestone().GetTitle()),
		Color: discord.Color(c.Colors.IssueMilestoned.Success),
	}
}

func (c Config) makeIssueDemilestonedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("%s removed milestone %s", ev.GetSender().GetLogin(), ev.GetMilestone().GetTitle()),
		Color: discord.Color(c.Colors.IssueDemilestoned.Success),
	}
}

func (c Config) makeIssueDeletedEmbed(ev *github.IssuesEvent) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("%s deleted this", ev.GetSender().GetLogin()),
		Color: discord.Color(c.Colors.IssueDeleted.Success),
	}
}

// IssueCommentEvent embeds

func (c Config) makeIssueCommentEmbed(issue *github.Issue, comment *github.IssueComment) discord.Embed {
	var fields *[]discord.EmbedField

	for react, count := range parseReactions(comment.Reactions) {
		*fields = append(*fields, discord.EmbedField{
			Name:   react,
			Value:  fmt.Sprintf("%d", count),
			Inline: true,
		})
	}

	return discord.Embed{
		Title:       fmt.Sprintf("Comment on issue #%d: %s", issue.GetNumber(), issue.GetTitle()),
		Description: convertMarkdown(comment.GetBody()),
		URL:         comment.GetHTMLURL(),
		Color:       discord.Color(c.Colors.IssueCommented.Success),
		Fields:      *fields,
		Author: &discord.EmbedAuthor{
			Name: comment.GetUser().GetLogin(),
			Icon: comment.GetUser().GetAvatarURL(),
		},
		Footer: &discord.EmbedFooter{Text: "0x" + strconv.FormatInt(comment.GetID(), 16)},
	}
}

// PullRequestEvent embeds

func (c Config) makePREmbed(pr *github.PullRequest) discord.Embed {
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

// PullRequestCommentEvent embeds

func (c Config) makePRCommentEmbed(pr *github.PullRequest, comment *github.PullRequestComment) discord.Embed {
	var fields *[]discord.EmbedField

	for react, count := range parseReactions(comment.Reactions) {
		*fields = append(*fields, discord.EmbedField{
			Name:   react,
			Value:  fmt.Sprintf("%d", count),
			Inline: true,
		})
	}

	return discord.Embed{
		Title:       fmt.Sprintf("Comment on pull request #%d: %s", pr.GetNumber(), pr.GetTitle()),
		Description: convertMarkdown(comment.GetBody()),
		URL:         comment.GetHTMLURL(),
		Color:       discord.Color(c.Colors.IssueCommented.Success),
		Fields:      *fields,
		Author: &discord.EmbedAuthor{
			Name: comment.GetUser().GetLogin(),
			Icon: comment.GetUser().GetAvatarURL(),
		},
		Footer: &discord.EmbedFooter{Text: "0x" + strconv.FormatInt(comment.GetID(), 16)},
	}
}

// PullRequestReviewEvent embeds

func (c Config) makePRReviewEmbed(pr *github.PullRequest, review *github.PullRequestReview) discord.Embed {
	return discord.Embed{
		Title:       pr.GetTitle(),
		Description: review.GetBody(),
		URL:         review.GetHTMLURL(),
	}
}

// PullRequestReviewCommentEvent embeds

func (c Config) makePRReviewCommentEmbed(pr *github.PullRequest, comment *github.PullRequestComment) discord.Embed {
	return discord.Embed{
		Title:       pr.GetTitle(),
		Description: comment.GetBody(),
		URL:         comment.GetHTMLURL(),
	}
}

// PullRequestReviewThreadEvent embeds

func (c Config) makePRReviewThreadEmbed(pr *github.PullRequest) discord.Embed {
	return discord.Embed{
		
}

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

func parseReactions(r *github.Reactions) (reacts map[string]int) {
	switch {
	case *r.PlusOne > 0:
		reacts["👍"] = *r.PlusOne
		fallthrough
	case *r.MinusOne > 0:
		reacts["👎"] = *r.MinusOne
		fallthrough
	case *r.Laugh > 0:
		reacts["😆"] = *r.Laugh
		fallthrough
	case *r.Hooray > 0:
		reacts["🎉"] = *r.Hooray
		fallthrough
	case *r.Confused > 0:
		reacts["😕"] = *r.Confused
		fallthrough
	case *r.Heart > 0:
		reacts["❤️"] = *r.Heart
		fallthrough
	case *r.Rocket > 0:
		reacts["🚀"] = *r.Rocket
		fallthrough
	case *r.Eyes > 0:
		reacts["👀"] = *r.Eyes
	}

	return
}
