package sessions

import "github.com/deployment-io/deployment-runner-kit/types"

// SuggestedRepositoryDtoV1 is one repository the planning agent suggested
// adding to the session. Name is a lookup key ("org/repo") the dashboard
// resolves against the org's repository list — the payload is display data, so
// no clone URL / provider / installation rides along. Confidence is a plain
// string ("high" | "medium" | "low"), normalised by agentbox's extractor and
// re-normalised server-side.
type SuggestedRepositoryDtoV1 struct {
	Name       string
	Reason     string
	Confidence string
}

// SetRepoSuggestionDtoV1 is the latest repository suggestion the planning agent
// emitted during a session (agentbox's repo-suggestion.json), forwarded by the
// runner to deployment-server, which persists it to Session.RepoSuggestion.
// Sent whenever the file changes; latest wins and is never cleared.
type SetRepoSuggestionDtoV1 struct {
	JobID        string
	Repositories []SuggestedRepositoryDtoV1
}

type SetRepoSuggestionArgsV1 struct {
	types.AuthArgsV1
	Suggestion SetRepoSuggestionDtoV1
}

type SetRepoSuggestionReplyV1 struct {
	Done bool
}
