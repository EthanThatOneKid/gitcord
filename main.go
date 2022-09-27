package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "embed"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/ethanthatonekid/gitcord/gitcord"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

// This file is expected to be invoked via GitHub Workflow.

// Current features:
// - https://etok.codes/acmcsuf.com/blob/main/scripts/close-issue-channel.js
// - https://github.com/EthanThatOneKid/acmcsuf.com/blob/main/scripts/create-issue-channel.js
// - https://github.com/EthanThatOneKid/acmcsuf.com/blob/main/scripts/create-message.js
// - https://github.com/EthanThatOneKid/acmcsuf.com/blob/main/.github/workflows/close_issue_channel.yaml
// - https://github.com/EthanThatOneKid/acmcsuf.com/blob/main/.github/workflows/create_issue_channel.yaml
// - https://github.com/EthanThatOneKid/acmcsuf.com/blob/main/.github/workflows/create_message.yaml

// See:
// https://stackoverflow.com/questions/62325286/run-github-actions-when-pull-requests-have-a-specific-label#comment122159108_62331521

// Features:
// - Create new issue channel+initial message on issues:opened,reopened
// - Close issue channel on issues:closed,deleted
// - Create new message on issues:opened,reopened

func main() {
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
		Before: func(ctx *cli.Context) error {
			channelID, err := discord.ParseSnowflake(os.Getenv("DISCORD_CHANNEL_ID"))
			if err != nil {
				return errors.Wrap(err, "failed to parse Discord channel ID")
			}

			guildID, err := discord.ParseSnowflake(os.Getenv("DISCORD_GUILD_ID"))
			if err != nil {
				return errors.Wrap(err, "failed to parse Discord guild ID")
			}

			config := gitcord.Config{
				GitHubRepo: os.Getenv("GITHUB_REPO"),
				GitHubOAuth: oauth2.StaticTokenSource(&oauth2.Token{
					AccessToken: os.Getenv("GITHUB_TOKEN"),
				}),
				DiscordToken:     os.Getenv("DISCORD_TOKEN"),
				DiscordChannelID: discord.ChannelID(channelID),
				DiscordGuildID:   discord.GuildID(guildID),
				Logger:           log.Default(),
			}

			if err := parseColors(map[string]*gitcord.StatusColors{
				"GITCORD_COLOR_ISSUE": &config.Colors.IssueOpened,
				"GITCORD_COLOR_PR":    &config.Colors.PROpened,
			}); err != nil {
				return err
			}

			app.client = gitcord.NewClient(config).WithContext(ctx.Context)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "issues",
				Usage: "manage issues",
				Subcommands: []*cli.Command{
					{
						Name:    "opened",
						Aliases: []string{"reopened"},
						Usage:   "create a new issue channel",
						Action:  app.issuesOpened,
					},
					{
						Name:  "closed",
						Usage: "remove an existing template",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("removed task template: ", cCtx.Args().First())
							return nil
						},
					},
				},
			},
		},
	}

	return app
}

func (app *App) issuesOpened(ctx *cli.Context) error {
	n, err := strconv.Atoi(ctx.Args().First())
	if err != nil {
		return errors.Wrap(err, "failed to parse issue number")
	}
	return app.client.CreateIssueChannel(n)
}

func (app *App) issuesClosed(ctx *cli.Context) error {
	n, err := strconv.Atoi(ctx.Args().First())
	if err != nil {
		return errors.Wrap(err, "failed to parse issue number")
	}
	return app.client.CloseIssueChannel(n)
}

func (app *App) pullRequestReviewSubmitted(ctx *cli.Context) error {
	prID, err := strconv.Atoi(ctx.Args().Get(1))
	if err != nil {
		return errors.Wrap(err, "failed to parse PR number")
	}

	commentID, err := strconv.ParseInt(ctx.Args().Get(1), 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse comment ID")
	}

	return app.client.SendPRReview(prID, commentID)
}

func parseColors(envMap map[string]*gitcord.StatusColors) error {
	for env, colors := range envMap {
		if err := parseColorsEnv(env+"_SUCCESS", &colors.Success); err != nil {
			return err
		}

		if err := parseColorsEnv(env+"_ERROR", &colors.Error); err != nil {
			return err
		}
	}

	return nil
}

func pareseColorEnv(env string, dst *discord.Color) error {
	val := os.Getenv(env)
	if val == "" {
		return nil
	}

	if !strings.HasPrefix(val, "#") {
		return fmt.Errorf("$%s: invalid color must be of format #XXXXXX", env)
	}

	c, err := strconv.ParseInt(strings.TrimPrefix(val,"#"), 10, 32)
	if err!==nil {
		return errors.Wrapf(err ,"$%s: invalid color", env)
	}

	*dst = discord.Color(c)
	return nil

}