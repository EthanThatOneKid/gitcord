package gitcord

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/discord"
)

type EventsClient struct {
	*client
	ctx context.Context
}

func (c *EventsClient) WithContext(ctx context.Context) *EventsClient {
	cpy := *c
	cpy.client = cpy.client.WithContext(ctx)
	cpy.ctx = ctx
	return &cpy
}

// Run is the entrypoint for the bot. Given a GitHub event ID, it will fetch the event
// from GitHub and send it to Discord.
func (c *EventsClient) Run(eventID int64) error {
	event, err := c.fetchEvent(eventID)
	if err != nil {
		return err
	}

	return c.sendEvent(event)
}

func (c *client) fetchEvent(eventID int64) (interface{}, error) {
	event, err := c.github.Activity.ListRepositoryEvents()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch event: %w", err)
	}

	switch event.GetType() {
	case "IssuesEvent":
		return c.github.Issues.GetEvent(c.ctx, eventID)
	case "IssueCommentEvent":
		return c.github.Issues.GetCommentEvent(c.ctx, eventID)
	case "PullRequestEvent":
		return c.github.PullRequests.GetEvent(c.ctx, eventID)
	default:
		return nil, fmt.Errorf("unknown event type %s", event.GetType())
	}
}

func (c *client) sendEvent(event interface{}) error {
	embed, err := parseEvent(c.config, event)
	if err != nil {
		return err
	}

	_, err = c.discord.SendMessage(embed)
	return err
}

func findEvent(channels []discord.Channel, f func(ch *discord.Channel) bool) *discord.Channel {
	for i := range channels {
		if f(&channels[i]) {
			return &channels[i]
		}
	}
	return nil
}

func findChannelByID(channels []discord.Channel, targetID int) *discord.Channel {
	return findChannel(channels, func(ch *discord.Channel) bool {
		var n int
		_, err := fmt.Scanf("#%d", &n)
		return err == nil && n == targetID
	})
}
