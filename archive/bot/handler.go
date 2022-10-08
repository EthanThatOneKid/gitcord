package bot

// See: https://github.com/diamondburned/acmregister/blob/c72c9311d3/acmregister/bot/handler.go
import (
	"context"
	"fmt"
	"logger"
	"sync"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

type Config struct {
	// DiscordToken is the Discord bot token
	DiscordToken string
	// DiscordChannelID is the ID of the parent channel in which all threads
	// will be created under
	DiscordChannelID discord.ChannelID
}

type Handler struct {
	s      *state.state
	ctx    context.Context
	cancel context.CancelFunc
	c      Config
	wg     sync.WaitGroup
}

// NewHandler creates a new Handler instance bound to the given State.
func NewHandler(s *state.State, c Config) *Handler {
	ctx, cancel := context.WithCancel(s.Context())
	return &Handler{
		s:      s.WithContext(ctx),
		ctx:    ctx,
		cancel: cancel,
		c:      c,
	}
}

// Wait waits for all background jobs to finish. This is useful for closing
// database connections.
func (h *Handler) Wait() {
	h.wg.Wait()
}

// Close waits for everything to be done, then closes up everything that it
// needs to.
func (h *Handler) Close() error {
	h.cancel()
	h.wg.Wait()
	return nil
}

func (h *Handler) Intents() gateway.Intents {
	return 0 |
		gateway.IntentGuilds |
		gateway.IntentDirectMessages
}

func (h *Handler) HandleInteraction(ev *discord.InteractionEvent) *api.InteractionResponse {
	defer func() {
		if panicked := recover(); panicked != nil {
			h.privateWarning(ev, fmt.Errorf("bug: panic occured: %v", panicked))
		}
	}()
}

// privateWarning is like privateErr, except the user does not get a reply back
// saying things have gone wrong. Use this if we don't intend to return after
// the error.
func (h *Handler) privateWarning(ev *discord.InteractionEvent, sendErr error) {
	h.logErr(ev.GuildID, sendErr)
	h.sendDMErr(ev, sendErr)
}

func (h *Handler) sendDMErr(ev *discord.InteractionEvent, sendErr error) {
	guild, err := h.store.GuildInfo(ev.GuildID)
	if err != nil {
		h.logErr(ev.GuildID, err)
		return
	}

	dm, err := h.s.CreatePrivateChannel(guild.InitUserID)
	if err != nil {
		h.logErr(ev.GuildID, err)
		return
	}

	if _, err = h.s.SendMessage(dm.ID, "⚠️ Error: "+sendErr.Error()); err != nil {
		h.logErr(ev.GuildID, errors.Wrap(err, "cannot send error to DM"))
		return
	}
}

func (h *Handler) logErr(guildID discord.GuildID, err error) {
	var guildInfo string
	if guild, err := h.s.Guild(guildID); err == nil {
		guildInfo = fmt.Sprintf("%q (%d)", guild.Name, guild.ID)
	} else {
		guildInfo = fmt.Sprintf("%d", guildID)
	}

	logger := logger.FromContext(h.ctx)
	logger.Println("guild "+guildInfo+":", "command error:", err)
}
