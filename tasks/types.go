// Package tasks holds shared types used by the Tasks feature across the
// control plane (kit, deployment-server) and the runner. The runner reads
// the JSON-encoded list under parameters_enums.Repositories and unmarshals
// it into []RepositoryEntry to iterate per-repo for CheckoutRepo,
// CommitAndPush, and OpenPullRequest commands.
package tasks

// RepositoryEntry is one repo's per-Job context. Tokens are intentionally
// absent — the runner mints fresh installation tokens on demand via the
// existing RefreshGitToken RPC keyed on InstallationID.
type RepositoryEntry struct {
	Name           string `json:"name"`
	CloneURL       string `json:"cloneUrl"`
	BaseBranch     string `json:"baseBranch"`
	Provider       string `json:"provider"` // matches oauth_enums.Provider.String()
	InstallationID string `json:"installationID"`
}
