package oauth

import "github.com/deployment-io/deployment-runner-kit/types"

// OpenPullRequestArgsV1 is the runner→deployment-server RPC payload for
// opening a PR/MR via the user's git provider installation. The server
// dispatches on the installation's provider (GitHub/GitLab/BitBucket),
// makes the provider REST call, and returns a normalized result.
//
// The provider-REST code lives server-side so:
//   - Three provider implementations live in one place (alongside the
//     existing webhook ingestion code in kit/dependencies/oauth/).
//   - PR-body templates / footer policy / error normalization can evolve
//     by deploying the server, without updating every user-side runner.
//   - Adding a fourth provider is a server-only change.
type OpenPullRequestArgsV1 struct {
	types.AuthArgsV1
	InstallationID string
	// RepoName is the provider-side repo identifier (e.g., "deployment-io/kit"
	// for GitHub, "deployment-io/kit" for GitLab, "deployment-io/kit" for
	// BitBucket workspace/slug). Server uses installation.OauthProvider to
	// interpret the format correctly.
	RepoName   string
	BaseBranch string
	HeadBranch string
	Title      string
	Body       string
	// Draft asks the provider to open the pull request in its draft state,
	// used by the Review stage when the loop exhausted with a must-fix
	// finding still open. A REQUEST, not a requirement: a provider with no
	// draft concept still creates an ordinary pull request and reports
	// DraftUnsupported. False behaves exactly as this call did before the
	// field existed.
	//
	// Honoured at CREATION only. Draft state cannot be set on a pull request
	// that already exists (GitHub's PATCH cannot change it), so a re-run
	// against an open PR carries the signal in its title and body instead.
	Draft bool
}

type OpenPullRequestDtoV1 struct {
	URL    string
	Number int
	// DraftUnsupported reports that a requested draft could not be honoured
	// by this provider. PURELY INFORMATIONAL AND NEVER A FAILURE: the pull
	// request was still created or updated, and URL and Number are populated
	// exactly as on any other successful call. The caller uses it to pick a
	// different way to signal "needs fixes" (a title prefix), not to decide
	// whether a pull request exists.
	DraftUnsupported bool
	// AlreadyExisted reports that the call resolved to a pull request that
	// was already open rather than creating one — the ordinary outcome of
	// re-running a Step against a branch that already has a PR. ALSO NOT A
	// FAILURE: the supplied title and body have been applied to that pull
	// request, and URL and Number name it. The caller uses it to know that
	// draft state could not be set (see Draft above).
	AlreadyExisted bool
}
