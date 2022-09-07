package gitcord

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/go-github/v42/github"
	githubactions "github.com/sethvargo/go-githubactions"

	githubshim "github.com/ethanthatonekid/gitcord/github"
)

// TODO: Base procedure choice on event triggered e.g.
// https://docs.github.com/en/developers/webhooks-and-events/events/github-event-types#createevent
// https://github.com/blend/require-conditional-status-checks/blob/f71534cb20327333691f98ada07c1eabcf662cc2/vendor/github.com/google/go-github/v42/github/event_types.go#L103
func ChooseProcedure() {}

// Run executes the Gitcord GitHub Action; which invokes a procedure
// depending on the triggered event type.
func Run(action *githubactions.Action) error {
	cfg, err := NewFromInputs(action)
	if err != nil {
		return err
	}

	action.Infof("%s %s", ansi.Color256Gold3Alt2.Apply("[CONFIG]  Timeout:"), cfg.Timeout)
	action.Infof("%s %s", ansi.Color256Gold3Alt2.Apply("[CONFIG] Interval:"), cfg.Interval)
	action.Infof("%s %s", ansi.Color256Gold3Alt2.Apply("[CONFIG]     Base:"), cfg.BaseSHA)
	action.Infof("%s %s", ansi.Color256Gold3Alt2.Apply("[CONFIG]     Head:"), cfg.HeadSHA)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	client, err := githubshim.NewClient(ctx, cfg.GitHubRootURL, cfg.GitHubToken)
	if err != nil {
		return err
	}

	cc, _, err := client.Repositories.CompareCommits(ctx, cfg.GitHubOrg, cfg.GitHubRepo, cfg.BaseSHA, cfg.HeadSHA, nil)
	if err != nil {
		return err
	}

	checks, err := cfg.GetChecks(ctx, client)
	if err != nil {
		return err
	}

	incomplete := []string{}
	for _, check := range checks {
		required, err := check.Required(cc)
		if err != nil {
			return err
		}

		if required {
			action.Infof("%s %s", ansi.Blue("Required Job:"), check.Job)
			incomplete = append(incomplete, check.Job)
		}
	}

	if len(incomplete) == 0 {
		return nil
	}

	return Wait(ctx, action, client, cfg, incomplete)
}

// Wait polls the checks associated with the the `HEAD` SHA for the current
// pull request until either
// - **all** required runs have completed successfully
// - **at least one** required run has completed unsuccessfully
// - the context is done (due to timeout)
func Wait(ctx context.Context, action *githubactions.Action, client *github.Client, cfg *Config, incomplete []string) error {
	var err error
	for {
		select {
		case <-ctx.Done():
			action.Errorf("Timed out waiting for checks to complete successfully")
			return context.Canceled
		default:
		}

		incomplete, err = CheckSatisfied(ctx, action, client, cfg, incomplete)
		if err != nil {
			return err
		}
		if len(incomplete) == 0 {
			return nil
		}

		action.Infof("Sleeping for %s...", cfg.Interval)
		time.Sleep(cfg.Interval)
	}
}

// CheckSatisfied determines which `required` checks have successfully completed
// and which are still running. If **any** checks have completed unsuccessfully,
// an error will be returned.
//
// This uses the `GET /repos/:org/:repo/commits/:ref/check-runs` API to get
// the status and conclusion of all checks associated with the `HEAD` SHA for
// the current pull request and checks the status of the `required` checks.
func CheckSatisfied(ctx context.Context, action *githubactions.Action, client *github.Client, cfg *Config, required []string) ([]string, error) {
	checkRuns, err := depaginateListCheckRunsForRef(ctx, client.Checks, cfg);
	if err != nil {
		return required, err
	}

	sortCheckRuns(checkRuns)
	known := map[string]github.CheckRun{}
	for _, run := range checkRuns {
		if run == nil || run.Name == nil {
			continue
		}
		name := *run.Name
		// NOTE: **If** there are multiple jobs with the same name, writing
		//       to this map in this way means the "last one wins". It is
		//       the responsibility of the repository owner to ensure that
		//       the check names are unique.
		known[name] = *run
	}

	failed := false
	incomplete := []string{}
	for _, name := range required {
		run := known[name]
		status := safeDereference(run.Status, "unknown")
		conclusion := safeDereference(run.Conclusion, "unknown")
		runID := safeInt64String(run.ID, "unknown")
		if status == "completed" {
			if conclusion == "success" {
				action.Infof(
					"%[1]s %[2]s (%[3]s %[4]s)",
					ansi.Green("Check was successful:"), // 1
					name,                                // 2
					ansi.Green("Run ID:"),               // 3
					runID,                               // 4
				)
				continue
			}

			failed = true
			action.Errorf(
				"Check was not successful: %[1]s (Conclusion: %[2]s, Run ID: %[3]s)",
				name,       // 1
				conclusion, // 2
				runID,      // 3
			)
			continue
		}

		incomplete = append(incomplete, name)
		action.Infof(
			"%[1]s %[2]s (%[3]s %[4]s, %[5]s %[6]s)",
			ansi.Color256Gold3Alt2.Apply("Check is not complete:"), // 1
			name,                                    // 2
			ansi.Color256Gold3Alt2.Apply("Status:"), // 3
			status,                                  // 4
			ansi.Color256Gold3Alt2.Apply("Run ID:"), // 5
			runID,                                   // 6
		)
	}

	if failed {
		return incomplete, ex.New("Some required checks were not successful")
	}

	return incomplete, nil
}

func depaginateListCheckRunsForRef(ctx context.Context, checksAPI *github.ChecksService, cfg *Config) ([]*github.CheckRun, error) {
	listOpts := github.ListCheckRunsOptions{}
	listOpts.PerPage = 100
	var checkRuns []*github.CheckRun
	for {
		listResult, resp, err := checksAPI.ListCheckRunsForRef(ctx, cfg.GitHubOrg, cfg.GitHubRepo, cfg.HeadSHA, &listOpts)
		if err != nil {
			return checkRuns, err
		}
		checkRuns = append(checkRuns, listResult.CheckRuns...);
		if resp.NextPage == 0 {
			break
		}
		listOpts.Page = resp.NextPage
	}
	return checkRuns, nil
}

func sortCheckRuns(checks []*github.CheckRun) {
	// left is the lowest priority run, right is the highest priority
	sort.SliceStable(checks, func(i, j int) bool {
		left := checks[i]
		right := checks[j]
		if left == nil || left.ID == nil {
			return true
		}
		if right == nil || right.ID == nil {
			return false
		}
		return *left.ID < *right.ID
	})
}

func safeDereference(s *string, fallback string) string {
	if s == nil {
		return fallback
	}

	return *s
}

func safeInt64String(i *int64, fallback string) string {
	if i == nil {
		return fallback
	}

	return fmt.Sprintf("%d", *i)
}
