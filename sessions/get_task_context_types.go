package sessions

import "github.com/deployment-io/deployment-runner-kit/types"

// GetTaskContextArgsV1 asks for the record of the Task an Ask session is about
// (a session with a TaskID). The session is found through the session Job's
// thread, like UpdateSpecV1.
type GetTaskContextArgsV1 struct {
	types.AuthArgsV1
	JobID string
}

// GetTaskContextReplyV1 is the Task's record. Found is false, and nothing else
// is set, when the session has no TaskID (every non-Ask session) or its Task no
// longer exists. The runner writes the record under /work/context/task.
type GetTaskContextReplyV1 struct {
	Found bool
	Task  TaskContextV1
	// Runs are the agent runs of the Task's Steps, oldest first, at most the
	// newest 10. A run whose Job no longer exists is left out.
	Runs []TaskRunContextV1
}

// TaskContextV1 is what the Task is: what was asked and where it landed.
type TaskContextV1 struct {
	ID                 string
	Title              string
	Description        string
	AcceptanceCriteria []string
	// SpecJSON is the Task's structured spec as JSON; "" when it has none.
	SpecJSON     string
	Repositories []TaskRepositoryContextV1
}

// TaskRepositoryContextV1 is one of the Task's repositories.
type TaskRepositoryContextV1 struct {
	Name       string
	BaseBranch string
	TaskBranch string
	// CheckedOutBranch is the branch the Ask session checked this repository
	// out at (the session repository at the same index).
	CheckedOutBranch string
	PRURL            string
	PRState          string
}

// TaskRunContextV1 is one agent run (Job) of the Task.
type TaskRunContextV1 struct {
	JobID     string
	CreatedTs int64
	State     string
	// OutputJSON is the Job's Output as JSON; "" when it has none.
	OutputJSON    string
	RerunFeedback string
	// EventsJSONL is the Job's events, one JSON object per line.
	EventsJSONL string
	// Log is the run's log, oldest line first, one line per job log entry; at
	// most its newest 1 MiB. LogTruncated is true when older lines were left
	// out for that cap.
	Log          string
	LogTruncated bool
}
