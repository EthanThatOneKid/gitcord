// Package github is a focused wrapper around `github.com/google/go-github`
// with methods needed for `ethanthatonekid/gitcord`.
//
// Adapted from
// https://github.com/blend/require-conditional-status-checks/blob/f71534cb20/pkg/github/client.go

package github

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/blend/go-sdk/ex"
	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

// NewClient creates a new client and determines if it's needed for the public
// GitHub API or GitHub Enterprise. The `rootURL` is expected to be the value
// of the `GITHUB_API_URL` environment variable / the `${{ github.api_url }}`
// context value.
func NewClient(ctx context.Context, rootURL, token string) (*github.Client, error) {
	sts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, sts)
	if rootURL == "https://api.github.com" {
		return github.NewClient(tc), nil
	}
}

// GetFile downloads the contents of a file.
func GetFile(ctx context.Context, c *github.Client, owner, repo, ref, path string) ([]byte, error) {
	opts := &github.RepositoryContentGetOptions{Ref: ref}
	rc, response, err := c.Repositories.DownloadContents(ctx, owner, repo, path, opts)
	if err != nil {
		return nil, ex.New("Failed to download file", ex.OptMessagef("Repository: %s/%s, Ref: %s, Path: %s", owner, repo, ref, path), ex.OptInner(err))
	}
	defer rc.Close()

	if response.StatusCode != http.StatusOK {
		return nil, ex.New("Raw download HTTP failure", ex.OptMessagef("Status Code: %d, Repository: %s/%s, Ref: %s, Path: %s", response.StatusCode, owner, repo, ref, path), ex.OptInner(err))
	}

	body, err := io.ReadAll(rc)
	if err != nil {
		return nil, ex.New("Failed to read body of raw download", ex.OptMessagef("Repository: %s/%s, Ref: %s, Path: %s", owner, repo, ref, path), ex.OptInner(err))
	}

	return body, nil
}
