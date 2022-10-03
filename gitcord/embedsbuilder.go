package gitcord

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/google/go-github/v47/github"
	smore "github.com/pythonian23/SMoRe"
)

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

func convertMarkdown(githubMD string) string {
	// See https://github.com/pythonian23/SMoRe.
	return smore.Render(githubMD)
}

func parseEvent(c Config, event interface{}) (embed discord.Embed, err error) {
	switch event := event.(type) {
	case *github.IssueEvent:
		embed = parseIssueEvent(c, event)
	case *github.IssueCommentEvent:
		embed = parseIssueCommentEvent(c, event)
	case *github.PullRequestEvent:
		embed = parsePullRequestEvent(c, event)
	case *github.PullRequestReviewEvent:
		embed = parsePullRequestReviewEvent(c, event)
	case *github.PullRequestReviewCommentEvent:
		embed = parsePullRequestReviewCommentEvent(c, event)
	case *github.PushEvent:
		embed = parsePushEvent(c, event)
	case *github.ReleaseEvent:
		embed = parseReleaseEvent(c, event)
	case *github.RepositoryEvent:
		embed = parseRepositoryEvent(c, event)
	case *github.CreateEvent:
		embed = parseCreateEvent(c, event)
	case *github.DeleteEvent:
		embed = parseDeleteEvent(c, event)
	case *github.ForkEvent:
		embed = parseForkEvent(c, event)
	case *github.WatchEvent:
		embed = parseWatchEvent(c, event)
	case *github.GollumEvent:
		embed = parseGollumEvent(c, event)
	case *github.CommitCommentEvent:
		embed = parseCommitCommentEvent(c, event)
	default:
		err = fmt.Errorf("unknown event type %T", event)
	}

	return
}

func BuildEmbed(c Config, event interface{}) discord.Embed {
	switch event := event.(type) {
	case *github.IssueEvent:
		switch event.GetAction() {
		case "opened":
			return embedOpeningIssueChannelMsg(c, event.GetIssue())
		case "closed":
			return discord.Embed{
				Title: fmt.Sprintf("Issue closed: #%d %s", event.GetIssue().GetNumber(), event.GetIssue().GetTitle()),
				URL:   event.GetIssue().GetURL(),





type embeds = map[string](map[string]func(Config, *github.Issue, *github.IssueComment) discord.Embed)

const EMBEDS = embeds{
	"issues": {
		"opened": func(c Config, issue *github.Issue) discord.Embed {
			return discord.Embed{
				Title: fmt.Sprintf("#%d: %s", issue.GetNumber(), issue.GetTitle()),
				URL:   issue.GetURL(),
				Author: &discord.EmbedAuthor{
					Name: issue.GetUser().GetLogin(),
					Icon: issue.GetUser().GetAvatarURL(),
				},
				Description: convertMarkdown(issue.GetBody()),
				Color:       discord.Color(c.Colors.IssueOpened.Success),
				Footer:      &discord.EmbedFooter{Text: fmt.Sprintf("Labels: %s", convertLabels(issue.Labels))},
			}
		},
	},
	// "issue_comment": {
	// 	// "created": func(c Config, a ...any) discord.Embed {
	// 	// 	issue := a[0].(*github.Issue)
	// 	// 	comment := a[1].(*github.IssueComment)
	// 	// 	return embedIssueComment(c, issue, comment)
	// 	// }},
	// },
	// "pull_request_review_comment": embedPRComment,
}

func embedOpeningIssueChannelMsg(c Config, issue *github.Issue) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("#%d: %s", issue.GetNumber(), issue.GetTitle()),
		URL:   issue.GetURL(),
		Author: &discord.EmbedAuthor{
			Name: issue.GetUser().GetLogin(),
			Icon: issue.GetUser().GetAvatarURL(),
		},
		Description: convertMarkdown(issue.GetBody()),
		Color:       discord.Color(c.Colors.IssueOpened.Success),
		Footer:      &discord.EmbedFooter{Text: fmt.Sprintf("Labels: %s", convertLabels(issue.Labels))},
	}
}

func embedOpeningPRChannelMsg(c Config, pr *github.PullRequest) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Pull request opened: #%d %s", pr.GetNumber(), pr.GetTitle()),
		URL:   pr.GetURL(),
		Author: &discord.EmbedAuthor{
			Name: pr.GetUser().GetLogin(),
			Icon: pr.GetUser().GetAvatarURL(),
		},
		Description: convertMarkdown(pr.GetBody()),
		Color:       discord.Color(c.Colors.IssueOpened.Success),
		Footer:      &discord.EmbedFooter{Text: fmt.Sprintf("Labels: %s\nAssignees: %s\nReviewers: %s", convertLabels(pr.Labels), convertUsers(pr.Assignees), convertUsers(pr.RequestedReviewers))},
	}
}

func embedIssueComment(c Config, issue *github.Issue, comment *github.IssueComment) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Comment on issue #%d: %s", issue.GetNumber(), issue.GetTitle()),
		URL:   comment.GetURL(),
		Author: &discord.EmbedAuthor{
			Name: comment.GetUser().GetLogin(),
			Icon: comment.GetUser().GetAvatarURL(),
		},
		Description: convertMarkdown(comment.GetBody()),
		Color:       discord.Color(c.Colors.IssueCommented.Success),
		Footer:      &discord.EmbedFooter{Text: "0x" + strconv.FormatInt(comment.GetID(), 16)},
	}
}

func embedPRComment(c Config, pr *github.PullRequest, comment *github.PullRequestComment) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Comment on issue #%d: %s", pr.GetNumber(), pr.GetTitle()),
		URL:   comment.GetURL(),
		Author: &discord.EmbedAuthor{
			Name: comment.GetUser().GetLogin(),
			Icon: comment.GetUser().GetAvatarURL(),
		},
		Description: convertMarkdown(pr.GetBody()),
		Color:       discord.Color(c.Colors.PRCommented.Success),
		Footer:      &discord.EmbedFooter{Text: "0x" + strconv.FormatInt(comment.GetID(), 16)},
	}
}

func embedPRReview(c Config, pr *github.PullRequest, comment *github.PullRequestComment) discord.Embed {
	return discord.Embed{
		Title: fmt.Sprintf("Pull request review submitted: #%d %s", pr.GetNumber(), pr.GetTitle()),
		URL:   comment.GetURL(),
		Author: &discord.EmbedAuthor{
			Name: comment.GetUser().GetLogin(),
			Icon: comment.GetUser().GetAvatarURL(),
		},
		Description: convertMarkdown(convertMarkdown(comment.GetBody())),
		Color:       discord.Color(c.Colors.PRCommented.Success),
		Footer:      &discord.EmbedFooter{Text: "0x" + strconv.FormatInt(comment.GetID(), 16)},
	}
}
