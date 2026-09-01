package jobs

import "github.com/deployment-io/deployment-runner-kit/types"

// ReleasingJobsArgsV1 asks the server to put jobs the runner accepted but
// could not start back into the pending queue.
//
// This exists because jobs are marked Running the moment GetPending hands
// them out. Before this, a runner that could not start a job had two
// options, both bad: block a worker waiting for capacity, or let the job
// sit until the stuck-job cron marked it TimedOut — a FAILED deployment
// caused by nothing but the runner being busy.
//
// Releasing is not a failure. The job returns to Pending and is offered
// again on a later poll, so a busy runner delays work instead of losing
// it.
type ReleasingJobsArgsV1 struct {
	types.AuthArgsV1
	// JobIDs are hex ObjectIDs of jobs to return to the pending queue.
	JobIDs []string
	// Reason is recorded in the server log so an operator can tell a
	// capacity release apart from any future release cause. Not shown to
	// the user; the job simply reappears as pending.
	Reason string
}

type ReleasingJobsReplyV1 struct {
	Done bool
}
