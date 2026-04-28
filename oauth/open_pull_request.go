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
}

type OpenPullRequestDtoV1 struct {
	URL    string
	Number int
}
