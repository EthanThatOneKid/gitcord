package discordclient

import (
	"context"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/ethanthatonekid/gitcord/gitcord/internal/slices"
	"github.com/pkg/errors"
)

// Client is a wrapped Discord client.
type Client struct {
	*api.Client
	config Config
}

// Config is the configuration for the Discord client.
type Config struct {
	Token     string
	ChannelID discord.ChannelID
	// Logger is optional. By default, it will log to the standard logger.
	Logger *log.Logger
}

// New creates a new Discord client.
func New(cfg Config) *Client {
	if cfg.Logger == nil {
		cfg.Logger = log.Default()
	}

	return &Client{
		Client: api.NewClient(cfg.Token),
		config: cfg,
	}
}

func (c *Client) logln(v ...any) {
	prefixed := []any{"discord:"}
	prefixed = append(prefixed, v...)
	c.config.Logger.Println(prefixed...)
}

func (c *Client) WithContext(ctx context.Context) *Client {
	return &Client{
		Client: c.Client.WithContext(ctx),
		config: c.config,
	}
}

func (c *Client) guildID() (discord.GuildID, error) {
	ch, err := c.Channel(c.config.ChannelID)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get guild ID for channel %d", c.config.ChannelID)
	}
	return ch.GuildID, nil
}

func (c *Client) activeThreads() ([]discord.Channel, error) {
	guildID, err := c.guildID()
	if err != nil {
		return nil, err
	}

	active, err := c.Client.ActiveThreads(guildID)
	if err != nil {
		return nil, err
	}

	// Filter for relevant threads only.
	relevant := slices.FilterReuse(active.Threads, func(ch *discord.Channel) bool {
		return ch.ParentID == c.config.ChannelID
	})
	return relevant, nil
}

func (c *Client) FindThreadByNumber(id int) *discord.Channel {
	chs, err := c.activeThreads()
	if err != nil {
		return nil
	}

	return findChannelByNumber(chs, id)
}

func findChannelByNumber(channels []discord.Channel, targetID int) *discord.Channel {
	return slices.Find(channels, func(ch *discord.Channel) bool {
		var n int
		_, err := fmt.Sscanf(ch.Name, "%d", &n)
		return err == nil && n == targetID
	})
}

func (c *Client) FindMsgByComment(ch *discord.Channel, commentID int64) *discord.Message {
	return c.findMsg(ch, false, func(msg *discord.Message) bool {
		var id int64

		if len(msg.Embeds) != 1 {
			return false
		}

		_, err := fmt.Sscanf(msg.Embeds[0].Footer.Text, "0x%x", &id)
		if err != nil {
			return false
		}

		return err == nil && id == commentID
	})
}

func (c *Client) FindMsgByIssue(ch *discord.Channel, issueID int) *discord.Message {
	return c.findMsg(ch, true, func(msg *discord.Message) bool {
		var id int

		if len(msg.Embeds) != 1 {
			return false
		}

		_, err := fmt.Sscanf(msg.Embeds[0].Title, "%d", &id)
		if err != nil {
			return false
		}

		return err == nil && id == issueID
	})
}

func (c *Client) findMsg(ch *discord.Channel, fromTop bool, f func(msg *discord.Message) bool) *discord.Message {
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

		for i := range msgs {
			if f(&msgs[i]) {
				return &msgs[i]
			}
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
