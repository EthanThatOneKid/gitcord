# Gitcord 😎

_Open the floor on Discord._

## Usage

### Using the CLI tool

TODO: Explain that for most cases, we do not install the CLI manually. We invoke this CLI from GitHub Actions.

TODO Look into `synchronize` action type.
issues:
    types: [opened, reopened, closed, deleted]
  # https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#pull_request
  pull_request:
    types: [opened, reopened, closed, deleted, ready_for_review]
  # https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#issue_comment
  issue_comment:
    types: [created, edited, deleted]
  # https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#pull_request_review
  pull_request_review:
    types: [submitted]

- Commands that will create a new thread on Discord
  - `go run . issues opened`
  - `go run . pull_request opened`
- Commands that will send a 'reopened' message to Discord  
  - `go run . issues reopened`
  - `go run . pull_request reopened`
- Commands that will send a 'closed' message to Discord and archive the thread
  - `go run . issues closed`
  - `go run . pull_request closed`
  

```sh
go run . issues opened
```

### Set up in your GitHub repository

In order to set up Gitcord for a GitHub repository, you must first gather some information.
Collect the information listed below in any order.

- GitHub personal access token: <https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token>
- 

## Dev 👩‍💻

```sh
```

### Testing 👷‍♂️

_TODO: Create gitcord package tests_