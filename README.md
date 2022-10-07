# Gitcord 😎

_expand GitHub into Discord_

## Background

Gitcord is a Go package and command line tool for automating Discord management via GitHub Events.

By default, Gitcord supports an opinionated set of Discord operations that are conducive to collaborative development on GitHub.
This opinionated set of Discord operations include opening a new Discord thread for every newly opened GitHub issue and pull request, forwarding issue and pull request comments to Discord, and more.

## Usage

### Set up your Discord server

#### Requirements

- One Discord guild (a.k.a. Discord server)
  - **How to obtain**: Either find an existing server or create a new Discord server
  - **Why?**: Gitcord manages a dedicated Discord text channel within a Discord server
- One Discord text channel ID
  - **How to obtain**: Set value of `$DISCORD_CHANNEL_ID` to the dedicated text channel's ID
  - **Why?**: A Discord text channel is dedicated to being managed by Gitcord
- One Discord bot token
  - **How to obtain**: Create a new Discord bot in the Discord developer settings ([documentation](https://discord.com/developers/docs/topics/oauth2)) and set the value of `$
  - **Why?**: A Discord bot is used as an agent to manage the given Discord text channel

#### Inviting the Discord bot to your server

_TODO: Provide instructions for inviting the Discord bot to the desired server._

### Setting up your GitHub repository

In order to set up Gitcord for a GitHub repository, you must first gather some information.
Collect the information listed below in any order.

#### Requirements

- One GitHub repository
  - **How to obtain**: A GitHub repository may be created on [`github.com/new`](https://github.com/new) and set value of `$GITHUB_REPO` in the format owner/repo
  - **Why?**: Gitcord is intended to execute on GitHub repository events
- One GitHub access token
  - **How to obtain**: Generate a personal access token in the GitHub developer settings ([documentation](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)) and set the value of `$GITHUB_TOKEN` to it
  - **Why?**: An access token is used to authorize the reading of GitHub repository event data
- (Recommended) One entrypoint GitHub Workflow file
  - **How to obtain**: Create a new file `.github/workflows/[your_workflow_filename].yaml` (example: [`.github/workflows/gitcord.yaml`](.github/workflows/gitcord.yaml))
  - **Why?**: Execute the `gitcord` tool via GitHub Workflow event triggers (pass GitHub event payload via stdin)

## Dev 👩‍💻

### Using the tool

```sh
cp .env.example .env # TODO(newdev): populate .env
source .env

go run . <event_id>

# ex: go run . 24292424235
```

#### Passing GitHub event by ID

A real GitHub repository event ID may be passed (as the first argument) to the `gitcord` tool.

The program attempts to fetch the GitHub event by the ID passed via the tool, then execute the expected behavior.

#### Passing GitHub event via Stdin

Arbitrary GitHub event data may be passed to the `gitcord` tool via stdin.

There are a couple of cool reasons why I wanted to see this feature implemented:

1. This is a slight optimization (skips the initial payload fetch for the GitHub event)
1. This allows us to pass arbitrary/imaginary GitHub events (helpful for debugging in production)

### Testing 👷‍♂️

_TODO: Create gitcord package tests_

---

Created with 😎 by [**ACM at CSUF**](https://acmcsuf.com/about)
