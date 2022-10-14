package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/ethanthatonekid/gitcord/gitcord"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

func main() {
	godotenv.Load(".env")
	app := NewApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type App struct {
	*cli.App
	client *gitcord.Client
}

func NewApp() *App {
	app := &App{}

	app.App = &cli.App{
		Name:     "gitcord",
		HelpName: "expand GitHub into Discord",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "force open threads",
			},
		},
		Action: func(ctx *cli.Context) error {
			channelID, err := discord.ParseSnowflake(os.Getenv("DISCORD_CHANNEL_ID"))
			if err != nil {
				return errors.Wrap(err, "failed to parse Discord channel ID")
			}

			colors, err := parseEnvColors()
			if err != nil {
				return err
			}

			config := gitcord.Config{
				GitHubOAuth: oauth2.StaticTokenSource(&oauth2.Token{
					AccessToken: os.Getenv("GITHUB_TOKEN"),
				}),
				DiscordToken:     "Bot " + os.Getenv("DISCORD_TOKEN"),
				DiscordChannelID: discord.ChannelID(channelID),
				ColorScheme:      colors,
				ForceOpen:        ctx.Bool("force"),
				Logger:           log.Default(),
			}

			app.client = gitcord.NewClient(config).WithContext(ctx.Context)

			eventIDStr := ctx.Args().First()
			switch eventIDStr {
			case "":
				eventName := os.Getenv("GITHUB_EVENT_NAME")
				if eventName == "" {
					return errors.New("no event ID provided")
				}

				eventPayload := os.Getenv("GITHUB_EVENT_PAYLOAD")
				if eventPayload == "" {
					return errors.New("no github event payload provided")
				}

				return app.client.DoEventPayload(eventName, eventPayload)

			default:
				eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
				if err != nil {
					return errors.Wrap(err, "failed to parse GitHub event ID")
				}

				return app.client.DoEventID(eventID)
			}
		},
	}

	return app
}

// colorEnvMap maps environment variable prefixes to their respective color
// scheme key.
var colorEnvMap = map[string]gitcord.ColorSchemeKey{
	"GITCORD_COLOR_ISSUE_OPENED": gitcord.IssueOpened,
}

func parseEnvColors() (gitcord.ColorScheme, error) {
	newScheme := gitcord.ColorScheme{}

	for env, schemeKey := range colorEnvMap {
		colors := gitcord.DefaultStatusColors
		if err := parseColorEnv(env+"_SUCCESS", &colors.Success); err != nil {
			return nil, err
		}
		if err := parseColorEnv(env+"_ERROR", &colors.Error); err != nil {
			return nil, err
		}
		newScheme[schemeKey] = colors
	}

	return newScheme, nil
}

// parseColorEnv parses a color from an environment variable into dst. If the
// environment variable is not set, then dst is not modified.
func parseColorEnv(env string, dst *discord.Color) error {
	val := os.Getenv(env)
	if val == "" {
		return nil
	}

	if !strings.HasPrefix(val, "#") {
		return fmt.Errorf("$%s: invalid color must be of format #XXXXXX", env)
	}

	c, err := strconv.ParseInt(strings.TrimPrefix(val, "#"), 10, 32)
	if err != nil {
		return errors.Wrapf(err, "$%s: invalid color", env)
	}

	*dst = discord.Color(c)
	return nil

}
