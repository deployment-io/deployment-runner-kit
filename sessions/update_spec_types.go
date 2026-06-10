package sessions

import "github.com/deployment-io/deployment-runner-kit/types"

// UpdateSpecDtoV1 is the latest structured task-spec the planning agent emitted
// during a session (agentbox's task-spec.json), forwarded by the runner to
// deployment-server, which persists it to Session.StructuredSpec. Sent whenever
// the spec file changes (the agent refines it across turns).
type UpdateSpecDtoV1 struct {
	JobID       string
	Title       string
	Goal        string
	Context     string
	Acceptance  []string
	Assumptions []string
	OutOfScope  []string
	Complexity  string
	Readiness   string
	Notes       string
}

type UpdateSpecArgsV1 struct {
	types.AuthArgsV1
	Spec UpdateSpecDtoV1
}

type UpdateSpecReplyV1 struct {
	Done bool
}
