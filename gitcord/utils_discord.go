package gitcord

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
)

func convertMarkdown(githubMD string) string {
	// See https://github.com/pythonian23/SMoRe.
	return githubMD
}

func (c *client) activeThreads() ([]discord.Channel, error) {
	active, err := c.discord.ActiveThreads(c.config.DiscordGuildID)
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

func (c *client) findMsgByID(ch *discord.Channel, targetID int64, top bool) *discord.Message {
	var lastID discord.MessageID
	msgs := make([]discord.Message, 0, 100)

	for len(msgs) > 0 {
		var err error
		if top {
			msgs, err = c.discord.MessagesAfter(ch.ID, lastID, 100)
		} else {
			msgs, err = c.discord.MessagesBefore(ch.ID, lastID, 100)
		}

		if err != nil {
			c.logln("failed to load messages:", err)
			return nil
		}

		msg := findMsgByID(msgs, targetID)
		if msg != nil {
			return msg
		}

		if top {
			lastID = msgs[len(msgs)-1].ID
		} else {
			lastID = msgs[0].ID
		}
	}

	return nil
}
