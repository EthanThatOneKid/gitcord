package gitcord

import (
	"fmt"
	"strconv"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/google/go-github/v47/github"
)

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
