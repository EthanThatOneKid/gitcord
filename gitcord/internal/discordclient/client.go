package discordclient

import (
	"context"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
)

// Client is a wrapped Discord client.
type Client struct {
	*api.Client
	config Config
}

type Config struct {
	DiscordToken     string
	DiscordGuildID   discord.GuildID
	DiscordChannelID discord.ChannelID
	Logger           *log.Logger
}

func New(cfg Config) *Client {
	return &Client{
		Client: api.NewClient(cfg.DiscordToken),
		config: cfg,
	}
}

func (c *Client) logln(v ...any) {
	if c.config.Logger != nil {
		prefixed := []any{"discord:"}
		prefixed = append(prefixed, v...)
		c.config.Logger.Println(prefixed...)
	}
}

func (c *Client) WithContext(ctx context.Context) *Client {
	return &Client{
		Client: c.Client.WithContext(ctx),
		config: c.config,
	}
}

func (c *Client) ActiveThreads() ([]discord.Channel, error) {
	active, err := c.Client.ActiveThreads(c.config.DiscordGuildID)
	if err != nil {
		return nil, err
	}

	relevantThreads := active.Threads[:0]
	for _, thread := range active.Threads {
		if thread.ParentID == c.config.DiscordChannelID {
			relevantThreads = append(relevantThreads, thread)
		}
	}

	return relevantThreads, nil
}

func (c *Client) FindChannelByIssue(id int) *discord.Channel {
	chs, err := c.ActiveThreads()

	if err != nil {
		c.logln("failed to load channels:", err)
		return nil
	}

	return findChannelByID(chs, id)
}

func findChannel(channels []discord.Channel, f func(ch *discord.Channel) bool) *discord.Channel {
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

func findMsg(msgs []discord.Message, f func(msg *discord.Message) bool) *discord.Message {
	for i := range msgs {
		if f(&msgs[i]) {
			return &msgs[i]
		}
	}
	return nil
}

func findMsgByID(msgs []discord.Message, targetID int64) *discord.Message {
	return findMsg(msgs, func(msg *discord.Message) bool {
		var id int64

		if len(msg.Embeds) != 1 {
			return false
		}

		_, err := fmt.Sscanf(msg.Embeds[0].Footer.Text, "0x%x", &id)
		if err != nil {
			return false
		}

		return err == nil && id == targetID
	})
}

func (c *Client) FindMsgByID(ch *discord.Channel, targetID int64) *discord.Message {
	return c.findMsgByID(ch, targetID, false)
}

func (c *Client) FindMsgByIDFromTop(ch *discord.Channel, targetID int64) *discord.Message {
	return c.findMsgByID(ch, targetID, true)
}

func (c *Client) findMsgByID(ch *discord.Channel, targetID int64, fromTop bool) *discord.Message {
	var lastID discord.MessageID
	msgs := make([]discord.Message, 0, 100)

	for len(msgs) > 0 {
		var err error
		switch {
		case fromTop:
			msgs, err = c.MessagesAfter(ch.ID, lastID, 100)
		default:
			msgs, err = c.MessagesBefore(ch.ID, lastID, 100)
		}

		if err != nil {
			c.logln("failed to load messages:", err)
			return nil
		}

		msg := findMsgByID(msgs, targetID)
		if msg != nil {
			return msg
		}

		switch {
		case fromTop:
			lastID = msgs[len(msgs)-1].ID
		default:
			lastID = msgs[0].ID
		}
	}

	return nil
}
