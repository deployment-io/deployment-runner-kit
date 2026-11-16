package tasks

import "github.com/deployment-io/deployment-runner-kit/types"

// QACriterionDtoV1 is one acceptance criterion the QA run checked, on the
// runner→server hop. State is one of the QACriterion* names below; like the
// review DTOs it is a string because the runner cannot import kit's
// qa_enums, and deployment-server parses it back and drops what does not
// parse.
//
// EvidenceIDs are the hex ids UploadQAEvidenceV1 returned for the
// screenshots cited for this criterion; empty for a runner that uploads
// none. A gob-added field: an older deployment-server ignores it, and the
// server keeps only ids stored for this Task and Job.
type QACriterionDtoV1 struct {
	Criterion   string
	State       string
	Reason      string
	EvidenceIDs []string
}

// QAResultDtoV1 is one QA stage run as the runner reports it.
//
// Status is one of the QAStatus* names below. StatusReason is the
// user-facing reason when the run was skipped or could not run. NoCriteria
// marks the smoke check run when the Task has no acceptance criteria.
// Findings reuse the Review stage's finding shape.
type QAResultDtoV1 struct {
	Status       string
	StatusReason string
	NoCriteria   bool
	Criteria     []QACriterionDtoV1
	Findings     []ReviewFindingDtoV1
	AgentType    string
	Model        string
	TokenUsage   ReviewTokenUsageDtoV1
	CostUSD      float64
}

// UpdateTaskStepQADtoV1 carries the final QA stage result for one Task Step
// to deployment-server (Jobs.UpdateTaskStepQAV1), which bounds it and stores
// it on the Step. Mirrors UpdateTaskStepReviewDtoV1.
type UpdateTaskStepQADtoV1 struct {
	TaskID    string
	StepIndex int
	JobID     string
	Result    QAResultDtoV1
}

type UpdateTaskStepQAArgsV1 struct {
	types.AuthArgsV1
	Update UpdateTaskStepQADtoV1
}

type UpdateTaskStepQAReplyV1 struct {
	Done bool
}

// UploadQAEvidenceDtoV1 carries one QA screenshot to deployment-server
// (Jobs.UploadQAEvidenceV1), which stores it and replies its id for
// QACriterionDtoV1.EvidenceIDs. One image per call.
//
// Image is a JPEG of at most 1 MiB (kit's qa_models.MaxEvidenceBytes); the
// server checks both and refuses anything else. CheckID names the check
// that took it ([a-z0-9_-], at most 16 characters); Criteria are indexes
// into the round's criteria it is evidence for. Width and Height are the
// image's pixel size.
type UploadQAEvidenceDtoV1 struct {
	TaskID    string
	StepIndex int
	JobID     string
	Round     int
	CheckID   string
	Criteria  []int
	Image     []byte
	Width     int
	Height    int
}

type UploadQAEvidenceArgsV1 struct {
	types.AuthArgsV1
	Upload UploadQAEvidenceDtoV1
}

type UploadQAEvidenceReplyV1 struct {
	EvidenceID string
}

// The WIRE MIRROR of kit's qa_enums: the names the runner sends in
// QACriterionDtoV1.State and QAResultDtoV1.Status. kit owns a test pinning
// each name to the enum value it must parse to.
const (
	QACriterionPassed     = "passed"
	QACriterionFailed     = "failed"
	QACriterionNotChecked = "not_checked"
)

const (
	QAStatusCompleted   = "completed"
	QAStatusSkipped     = "skipped"
	QAStatusCouldNotRun = "couldnt_run"
)
