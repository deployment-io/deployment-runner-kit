package tasks

import "github.com/deployment-io/deployment-runner-kit/types"

// UpdateTaskStepRunningDtoV1 reports that the runner has begun executing a
// Task Step — sent from CheckoutRepository.runForTask at checkout, the moment
// real work starts. On receipt deployment-server flips, idempotently from
// Pending: the Task -> Running, the Step -> StepRunning, and the agent-run Job
// -> Running. Mirrors how builds.UpdateBuildDtoV1 carries the runner-reported
// build Running status at checkout.
//
// Carries only identifiers — no state enum. Task/Step states live in kit's
// task_enums, which this lower-layer package can't import; the server-side
// handler owns the enum mapping.
type UpdateTaskStepRunningDtoV1 struct {
	TaskID    string
	StepIndex int
	JobID     string
}

// UpdateTaskStepRunningArgsV1 is the RPC envelope — a batch of step-started
// reports, mirroring builds.UpdateBuildsArgsV1.
type UpdateTaskStepRunningArgsV1 struct {
	types.AuthArgsV1
	Updates []UpdateTaskStepRunningDtoV1
}

type UpdateTaskStepRunningReplyV1 struct {
	Done bool
}
