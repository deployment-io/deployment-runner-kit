package logs

import (
	"github.com/deployment-io/deployment-runner-kit/types"
)

// Log stages a job log line can be tagged with. The runner owns and emits
// these; an empty Stage means the line is untagged.
const (
	LogStageSetup     = "setup"
	LogStageImplement = "implement"
	LogStageReview    = "review"
	LogStageQA        = "qa"
	LogStagePR        = "pr"
)

// Log kinds a job log line can be tagged with. The runner owns and emits
// these; an empty Kind means the line is untagged.
const (
	LogKindSection = "section"
	LogKindText    = "text"
	LogKindTool    = "tool"
	LogKindResult  = "result"
	LogKindJSON    = "json"
	LogKindDone    = "done"
	LogKindError   = "error"
)

// JobLogGroupV1 marks a log line that stands for Count lines of the same kind
// of noise that the runner folded into it (e.g. docker pull progress, agent
// tool chatter). Name names the group ("docker progress"); the line's Message
// is a sample or summary. Count is at least 1.
type JobLogGroupV1 struct {
	Name  string
	Count int
}

type AddJobLogDtoV1 struct {
	ID           string
	Message      string
	ErrorMessage string
	Ts           int64
	// Stage is one of the LogStage* constants; empty means untagged.
	Stage string
	// Kind is one of the LogKind* constants; empty means untagged.
	Kind string
	// Group is set when the line stands for a folded group of lines; nil otherwise.
	Group *JobLogGroupV1
}

type AddJobLogsArgsV1 struct {
	types.AuthArgsV1
	JobLogs []AddJobLogDtoV1
}

type AddJobLogsReplyV1 struct {
	Done bool
}
