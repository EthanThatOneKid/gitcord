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
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
						Name:   "opened",
						Usage:  "create a new issue channel",
						Action: app.issuesOpened,
					},
					{
						Name:   "reopened",
						Usage:  "forward issues reopened event",
						Action: app.issuesReopened,
					},
					{
						Name:   "edited",
						Usage:  "forward issues edit event",
						Action: app.issuesEdited,
					},
					{
						Name:   "closed",
						Usage:  "forward issues closed event",
						Action: app.issuesClosed,
					},
					{
						Name:   "deleted",
						Usage:  "forward issues deleted event",
						Action: app.issuesDeleted,
					},
					{
						Name:   "transferred",
						Usage:  "forward issues transferred event",
						Action: app.issuesTransferred,
					},
					{
						Name:   "assigned",
						Usage:  "forward issues assigned event",
						Action: app.issuesAssigned,
					},
					{
						Name:   "unassigned",
						Usage:  "forward issues unassigned event",
						Action: app.issuesUnassigned,
					},
					{
						Name:   "labeled",
						Usage:  "forward issues labeled event",
						Action: app.issuesLabeled,
					},
					{
						Name:   "unlabeled",
						Usage:  "forward issues unlabeled event",
						Action: app.issuesUnlabeled,
					},
				},
			},
			{
				Name:  "pull_request",
				Usage: "manage pull requests",
				Subcommands: []*cli.Command{
					{
						Name:   "opened",
						Usage:  "create a new pull request channel",
						Action: app.pull_requestOpened,
					}, {
						Name:   "reopened",
						Usage:  "forward pull_request reopened event",
						Action: app.pull_requestReopened,
					},
					{
						Name:   "closed",
						Usage:  "forward pull_request closed event",
						Action: app.pull_requestClosed,
					},
					{
						Name:   "deleted",
						Usage:  "forward pull_request deleted event",
						Action: app.pull_requestDeleted,
					},
					{
						Name:   "review_requested",
						Usage:  "forward pull_request review_requested event",
						Action: app.pull_requestReviewRequested,
					},
					{
						Name:   "review_request_removed",
						Usage:  "forward pull_request review_request_removed event",
						Action: app.pull_requestReviewRequestRemoved,
					},
					{
						Name:   "assigned",
						Usage:  "forward pull_request assigned event",
						Action: app.pull_requestAssigned,
					},
					{
						Name:   "unassigned",
						Usage:  "forward pull_request unassigned event",
						Action: app.pull_requestUnassigned,
					},
					{
						Name:   "labeled",
						Usage:  "forward pull_request labeled event",
						Action: app.pull_requestLabeled,
					},
					{
						Name:   "unlabeled",
						Usage:  "forward pull_request unlabeled event",
						Action: app.pull_requestUnlabeled,
					},
					{
						Name:   "converted_to_draft",
						Usage:  "forward pull_request converted_to_draft event",
						Action: app.pull_requestConvertedToDraft,
					},
					{
						Name:   "converted_to_draft",
						Usage:  "forward pull_request converted_to_draft event",
						Action: app.pull_requestConvertedToDraft,
					},
					{
						Name:   "ready_for_review",
						Usage:  "forward pull_request ready_for_review event",
						Action: app.pull_requestReadyForReview,
					},
					{
						Name:   "converted_to_draft",
						Usage:  "forward pull_request converted_to_draft event",
						Action: app.pull_requestConvertedToDraft,
					},
					{
						Name:   "auto_merge_enabled",
						Usage:  "forward pull_request auto_merge_enabled event",
						Action: app.pull_requestAutoMergeEnabled,
					},
					{
						Name:   "auto_merge_disabled",
						Usage:  "forward pull_request auto_merge_disabled event",
						Action: app.pull_requestAutoMergeDisabled,
					},
				},
			},
			{
				Name:  "issue_comment",
				Usage: "manage issue comments",
				Subcommands: []*cli.Command{
					{
						Name:   "created",
						Usage:  "forward a comment to the issue channel",
						Action: app.issue_commentCreated,
					},
					{
						Name:   "edited",
						Usage:  "forward a comment edit to the issue channel",
						Action: app.issue_commentEdited,
					},
				},
			},
		},
	}

	return app
}

func (app *App) issuesOpened(ctx *cli.Context) error {
	issueID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse issue number")
	}

	return app.client.Issues.OpenChannel(issueID)
}

func (app *App) issuesReopened(ctx *cli.Context) error {
	issueID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse issue number")
	}

	return app.client.Issues.ForwardReopened(issueID)
}

func (app *App) issuesEdited(ctx *cli.Context) error {
	issueID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse issue number")
	}

	return app.client.Issues.ForwardEdited(issueID)
}

func (app *App) issuesClosed(ctx *cli.Context) error {
	issueID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse PR number")
	}

	commentID, err := strconv.ParseInt(ctx.Args().Get(1), 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse comment ID")
	}
	return app.client.Issues.ForwardClosed(issueID, commentID)
}

func (app *App) issuesDeleted(ctx *cli.Context) error {
	issueID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse issue number")
	}

	return app.client.Issues.ForwardDeleted(issueID)
}

func (app *App) issuesTransferred(ctx *cli.Context) error {
	issueID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse issue number")
	}

	return app.client.Issues.ForwardTransferred(issueID)
}

func (app *App) issuesAssigned(ctx *cli.Context) error {
	issueID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse issue number")
	}

	return app.client.Issues.ForwardAssigned(issueID)
}

func (app *App) issuesUnassigned(ctx *cli.Context) error {
	issueID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse issue number")
	}

	return app.client.Issues.ForwardUnassigned(issueID)
}

func (app *App) issuesLabeled(ctx *cli.Context) error {
	issueID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse issue number")
	}

	return app.client.Issues.ForwardLabeled(issueID)
}

func (app *App) issuesUnlabeled(ctx *cli.Context) error {
	issueID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse issue number")
	}

	return app.client.Issues.ForwardUnlabeled(issueID)
}

func (app *App) issue_commentCreated(ctx *cli.Context) error {
	issueID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse issue number")
	}

	commentID, err := strconv.ParseInt(ctx.Args().Get(2), 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse comment ID")
	}
	return app.client.Comments.ForwardCreated(issueID, commentID)
}

func (app *App) issue_commentEdited(ctx *cli.Context) error {
	issueID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse issue number")
	}

	commentID, err := strconv.ParseInt(ctx.Args().Get(1), 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse comment ID")
	}
	return app.client.Comments.ForwardEdited(issueID, commentID)
}

func (app *App) pull_requestOpened(ctx *cli.Context) error {
	prID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse PR number")
	}

	return app.client.PRs.OpenChannel(prID)
}

func (app *App) pull_requestReopened(ctx *cli.Context) error {
	prID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse PR number")
	}

	return app.client.PRs.ForwardReopened(prID)
}

func (app *App) pull_requestClosed(ctx *cli.Context) error {
	prID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse PR number")
	}

	commentID, err := strconv.ParseInt(ctx.Args().Get(1), 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse comment ID")
	}

	return app.client.PRs.ForwardClosed(prID, commentID)
}

func (app *App) pull_requestDeleted(ctx *cli.Context) error {
	prID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse PR number")
	}

	return app.client.PRs.ForwardDeleted(prID)
}

func (app *App) pull_requestAssigned(ctx *cli.Context) error {
	prID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse PR number")
	}

	return app.client.PRs.ForwardAssigned(prID)
}

func (app *App) pull_requestUnassigned(ctx *cli.Context) error {
	prID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse PR number")
	}

	return app.client.PRs.ForwardUnassigned(prID)
}

func (app *App) pull_requestLabeled(ctx *cli.Context) error {
	prID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse PR number")
	}

	return app.client.PRs.ForwardLabeled(prID)
}

func (app *App) pull_requestUnlabeled(ctx *cli.Context) error {
	prID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse PR number")
	}

	return app.client.PRs.ForwardUnlabeled(prID)
}

func (app *App) pull_requestEdited(ctx *cli.Context) error {
	prID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse PR number")
	}

	return app.client.PRs.ForwardEdited(prID)
	// return errors.New("not implemented")
}

func (app *App) pull_requestReadyForReview(ctx *cli.Context) error {
	prID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse PR number")
	}

	return app.client.PRs.ForwardReadyForReview(prID)
	// return errors.New("not implemented")
}

func (app *App) pull_requestReviewRequested(ctx *cli.Context) error {
	prID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse PR number")
	}

	return app.client.PRs.ForwardReviewRequested(prID)
	// return errors.New("not implemented")
}

func (app *App) pull_requestReviewRequestRemoved(ctx *cli.Context) error {
	prID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse PR number")
	}

	return app.client.PRs.ForwardReviewRequestRemoved(prID)
	// return errors.New("not implemented")
}

func (app *App) pull_requestConvertedToDraft(ctx *cli.Context) error {
	prID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse PR number")
	}

	return app.client.PRs.ForwardConvertedToDraft(prID)
	// return errors.New("not implemented")
}

func (app *App) pull_requestAutoMergeEnabled(ctx *cli.Context) error {
	prID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse PR number")
	}

	return app.client.PRs.ForwardAutoMergeEnabled(prID)
	// return errors.New("not implemented")
}

func (app *App) pull_requestAutoMergeDisabled(ctx *cli.Context) error {
	prID, err := strconv.Atoi(ctx.Args().Get(0))
	if err != nil {
		return errors.Wrap(err, "failed to parse PR number")
	}

	return app.client.PRs.ForwardAutoMergeDisabled(prID)
	// return errors.New("not implemented")
}

func parseColors(envMap map[string]*gitcord.StatusColors) error {
	for env, colors := range envMap {
		if err := parseColorEnv(env+"_SUCCESS", &colors.Success); err != nil {
			return err
		}

		if err := parseColorEnv(env+"_ERROR", &colors.Error); err != nil {
			return err
		}
	}

	return nil
}

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
