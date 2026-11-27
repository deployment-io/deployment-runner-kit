package sessions

import "github.com/deployment-io/deployment-runner-kit/types"

// MaxPlanTasks is how many Tasks one plan may have. One constant: the
// runner's prompt, agentbox (told by the runner), deployment-server and
// app-server all read it, so raising it is a change here plus releases.
const MaxPlanTasks = 5

// PlanTaskDtoV1 is one Task of a plan the planning agent wrote: its spec (the
// fields of UpdateSpecDtoV1) and the session repositories it changes.
type PlanTaskDtoV1 struct {
	Key          string
	Title        string
	Goal         string
	Context      string
	Acceptance   []string
	Assumptions  []string
	OutOfScope   []string
	Complexity   string
	Readiness    string
	Notes        string
	Repositories []string // session repository names, "owner/name"
}

// PlanEdgeDtoV1: the Task To starts after the Task From's When
// ("merge") happened. From and To are keys of the plan's Tasks.
type PlanEdgeDtoV1 struct {
	From   string
	To     string
	When   string
	Reason string
}

// UpdatePlanDtoV1 is the latest plan of Tasks the planning agent wrote during
// a session, forwarded by the runner to deployment-server, which validates it
// and persists it to Session.StructuredPlan. Over gob an empty list arrives as
// nil; nil and empty are treated alike.
type UpdatePlanDtoV1 struct {
	JobID string
	Tasks []PlanTaskDtoV1
	Edges []PlanEdgeDtoV1
}

type UpdatePlanArgsV1 struct {
	types.AuthArgsV1
	Plan UpdatePlanDtoV1
}

// Refused is why deployment-server refused the plan ("" when stored).
// A refusal is not an error: the runner logs it once and doesn't retry.
type UpdatePlanReplyV1 struct {
	Done    bool
	Refused string
}
