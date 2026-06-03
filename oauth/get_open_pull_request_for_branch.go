package oauth

import "github.com/deployment-io/deployment-runner-kit/types"

// GetOpenPullRequestForBranchArgsV1 is the runner→deployment-server RPC
// payload for the read-only counterpart of OpenPullRequestV1. The server
// looks up the installation, dispatches on provider, and returns the
// most recent PR matching {RepoName, HeadBranch} or absence. Used by
// the runner's re-run flow (see Q15 in PLAN_tasks_verification.md) to
// decide whether to branch from base (no PR / closed-unmerged) or from
// the existing PR head (open), and to fail loudly when re-running
// against a merged PR.
//
// Same identifier conventions as OpenPullRequestArgsV1 — RepoName is
// "owner/repo" on GitHub / BitBucket, numeric ID or "namespace/path"
// on GitLab. Server interprets via installation.OauthProvider.
type GetOpenPullRequestForBranchArgsV1 struct {
	types.AuthArgsV1
	InstallationID string
	RepoName       string
	HeadBranch     string
}

// GetOpenPullRequestForBranchDtoV1 is the normalized response. Found
// discriminates "no PR has ever existed for this head" from "PR exists,
// here's its state." When Found is true, State and Merged together let
// the caller decide: open → iterate on this branch; closed+merged →
// the work is already in base, fail the re-run loudly; closed+!merged →
// treat as clean slate and re-attempt from base.
//
// HeadBranch / BaseBranch / HeadSHA are surfaced so the caller can
// fetch + checkout precisely without a second round-trip.
type GetOpenPullRequestForBranchDtoV1 struct {
	Found      bool
	PRNumber   int
	URL        string
	State      string // "open" or "closed" — provider-normalized.
	Merged     bool   // True only when State == "closed" AND the PR was merged.
	HeadBranch string
	BaseBranch string
	HeadSHA    string
}
