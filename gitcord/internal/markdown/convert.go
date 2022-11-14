package markdown

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/google/go-github/v47/github"
	"github.com/jaytaylor/html2text"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

type link struct {
	title string
	href  string
}

type links []link

var defaultMessageSize = 4096
var defaultTruncateSuffix = "…"

var mdConverter = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithXHTML(),
	),
)

// ConvertDiscord converts the given GitHub Flavored Markdown to plain text,
// truncating the result to the given maximum length.
func ConvertDiscord(md, readMoreURL string) string {
	var suffix = defaultTruncateSuffix
	if readMoreURL != "" {
		suffix += fmt.Sprintf(" %s", ConvertHyperlink("Read more", readMoreURL))
	}

	return ConvertTruncated(md, suffix, defaultMessageSize)
}

// ConvertTruncated converts the given GitHub Flavored Markdown to plain text,
// truncating the result to the given maximum length.
func ConvertTruncated(md, suffix string, maxLen int) string {
	var text string
	err := Convert(md, &text)
	if err != nil {
		errors.Wrap(err, "failed to convert GitHub Flavored Markdown to plain text")
	}

	return Truncate(text, suffix, maxLen)
}

// Truncate truncates the given string to the given maximum length.
func Truncate(s, suffix string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	// If the message is too long, we need to truncate it.
	return s[:strings.LastIndex(s[:maxLen-len(suffix)], " ")] + suffix
}

// Convert converts GitHub Flavored Markdown to plain text.
func Convert(md string, dst *string) error {
	var buf bytes.Buffer
	err := mdConverter.Convert([]byte(md), &buf)
	if err != nil {
		return err
	}

	*dst, err = html2text.FromString(buf.String(), html2text.Options{PrettyTables: true})
	if err != nil {
		return err
	}

	return nil
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
