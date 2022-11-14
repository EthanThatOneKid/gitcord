package markdown

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v47/github"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

type link struct {
	title string
	href  string
}

type links []link

var mdMaxSize = 4096
var defaultTruncateMd = "..."
var mdParser = goldmark.New(goldmark.WithExtensions(extension.GFM)).Parser()

func Convert(githubMD, readMoreURL string) string {
	bytes, buf := []byte(githubMD), strings.Builder{}
	node := mdParser.Parse(text.NewReader(bytes))
	if err := DefaultRenderer.Render(&buf, bytes, node); err != nil {
		return githubMD
	}

	s := buf.String()
	s = strings.TrimRight(s, "\n")

	if len(s) <= mdMaxSize {
		return s
	}

	// If the message is too long, we need to truncate it.
	suffix := defaultTruncateMd
	if readMoreURL != "" {
		suffix = fmt.Sprintf("... (%s)", ConvertHyperlink("_read more_", readMoreURL))
	}

	return s[:strings.LastIndex(s[:mdMaxSize-len(suffix)], " ")] + suffix
}

func ConvertThreadURL(t *github.PullRequestThread) (url string) {
	comment := t.Comments[0]
	if comment != nil {
		url = comment.GetHTMLURL()
	}
	return
}

func ConvertUsers(users []*github.User) string {
	var names links
	for _, user := range users {
		names = append(names, link{title: user.GetLogin(), href: user.GetHTMLURL()})
	}
	return names.commaSeparatedList()
}

func ConvertTeams(teams []*github.Team) string {
	var names links
	for _, team := range teams {
		names = append(names, link{title: team.GetName(), href: team.GetHTMLURL()})
	}
	return names.commaSeparatedList()
}

func ConvertHyperlink(content, href string) string {
	if href == "" {
		return content
	}
	return fmt.Sprintf("[%s](%s)", content, href)
}

func ConvertLabels[T *github.Label | github.Label](labelsV []T) string {
	var labels []github.Label
	switch labelsV := any(labelsV).(type) {
	case []github.Label:
		labels = labelsV
	case []*github.Label:
		labels = convertSlicePtr(labelsV)
	}

	var labelNames links
	for _, label := range labels {
		labelNames = append(labelNames, link{
			title: label.GetName(),
			href:  label.GetURL(),
		})
	}

	return labelNames.commaSeparatedList()
}

func (l links) commaSeparatedList() string {
	var b []string
	for _, link := range l {
		b = append(b, ConvertHyperlink(link.title, link.href))
	}

	return strings.Join(b, ", ")
}

func convertSlicePtr[T any](s []*T) []T {
	r := make([]T, 0, len(s))
	for _, v := range s {
		r = append(r, *v)
	}
	return r
}
