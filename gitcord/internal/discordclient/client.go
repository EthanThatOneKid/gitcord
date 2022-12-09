package discordclient

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

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

func (c *Client) threads() ([]discord.Channel, error) {
	guildID, err := c.guildID()
	if err != nil {
		return nil, err
	}

	filter := func(chs []discord.Channel) []discord.Channel {
		return slices.FilterReuse(chs, func(ch *discord.Channel) bool {
			return ch.ParentID == c.config.ChannelID
		})
	}

	active, err := c.Client.ActiveThreads(guildID)
	if err != nil {
		return nil, err
	}

	threads := filter(active.Threads)

	var prevArchivedThreadTime discord.Timestamp
	hasMore := true
	for hasMore {
		archive, err := c.Client.PublicArchivedThreads(c.config.ChannelID, prevArchivedThreadTime, 100)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load archived threads")
		}

		threads = append(threads, filter(archive.Threads)...)
		hasMore = archive.More
		if hasMore {
			earliest := archive.Threads[len(archive.Threads)-1]
			prevArchivedThreadTime = earliest.ThreadMetadata.ArchiveTimestamp
		}
	}

	return threads, nil
}

const (
	totalRetries  = 10
	retryWaitTime = 10 * time.Second
)

func (c *Client) FindThreadByNumber(id int) (*discord.Channel, error) {
	for i := 0; i < totalRetries; i++ {
		chs, err := c.threads()
		if err != nil {
			return nil, fmt.Errorf("failed to get threads: %w", err)
		}

		ch := findChannelByNumber(chs, id)
		if ch != nil {
			return ch, nil
		}

		ticker := time.NewTicker(retryWaitTime)
		defer ticker.Stop()

		select {
		case <-ticker.C:
			continue
		case <-c.Context().Done():
			return nil, c.Context().Err()
		}
	}

	return nil, fmt.Errorf("thread #%d not found", id)
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
		if len(msg.Embeds) != 1 {
			return false
		}

		// Parse int64 comment ID from the embed's footer.
		id, err := strconv.ParseInt(msg.Embeds[0].Footer.Text, 10, 64)
		if err != nil {
			return false
		}

		return err == nil && id == commentID
	})
}

var IssueMsgPrefix = "Issue opened: #"
var issueNumberRe = regexp.MustCompile(IssueMsgPrefix + `(\d+)`)

func (c *Client) FindMsgByIssue(ch *discord.Channel, issueID int) *discord.Message {
	return c.findMsg(ch, true, func(msg *discord.Message) bool {
		if len(msg.Embeds) != 1 {
			return false
		}

		matches := issueNumberRe.FindStringSubmatch(msg.Embeds[0].Title)
		if len(matches) != 2 {
			return false
		}

		n, err := strconv.Atoi(matches[1])
		if err != nil {
			return false
		}

		return n == issueID
	})
}

var PRMsgPrefix = "Pull request opened: #"
var prNumberRe = regexp.MustCompile(PRMsgPrefix + `(\d+)`)

func (c *Client) FindMsgByPR(ch *discord.Channel, issueID int) *discord.Message {
	return c.findMsg(ch, true, func(msg *discord.Message) bool {
		if len(msg.Embeds) != 1 {
			return false
		}

		matches := prNumberRe.FindStringSubmatch(msg.Embeds[0].Title)
		if len(matches) != 2 {
			return false
		}

		n, err := strconv.Atoi(matches[1])
		if err != nil {
			return false
		}

		return n == issueID
	})
}

func (c *Client) findMsg(ch *discord.Channel, fromTop bool, f func(msg *discord.Message) bool) *discord.Message {
	var lastID discord.MessageID
	msgs := make([]discord.Message, 0, 100)

	for {
		var err error
		if fromTop {
			msgs, err = c.MessagesAfter(ch.ID, lastID, 100)
		} else {
			msgs, err = c.MessagesBefore(ch.ID, lastID, 100)
		}

		if err != nil {
			c.logln("failed to load messages:", err)
			break
		}

		if len(msgs) == 0 {
			break
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
