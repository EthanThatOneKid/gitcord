# Gitcord 😎

_Open the floor on Discord._

## Usage

### Set up your GitHub repository

In order to set up Gitcord for a GitHub repository, you must first gather some information.
Collect the information listed below in any order.

- GitHub personal access token: <https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token>

### Set up your Discord server

- Discord bot token: <https://discord.com/developers/docs/topics/oauth2>

## Dev 👩‍💻

### Using the CLI tool

```sh
cp .env.example .env # populate .env
source .env

go run . <event_id>

# ex: go run . 24292424235
```

### Testing 👷‍♂️

_TODO: Create gitcord package tests_