package tasks

import (
	"strings"

	"github.com/deployment-io/deployment-runner-kit/types"
)

// ReviewFindingDtoV1 is one Review stage finding on the runner→server hop.
//
// Parameter, Severity and Stage are STRINGS here even though kit stores them
// as numbered enums. The producer is an agent writing prose into agentbox's
// result.json, and the runner cannot import kit (see the wire-mirror comment
// on ReviewParameterValue), so the name is what actually crosses every hop.
// deployment-server maps each back through review_enums' Parse functions and
// drops what does not parse rather than persisting a zero Parameter.
type ReviewFindingDtoV1 struct {
	Key       string
	Parameter string
	Severity  string
	Location  string
	What      string
	Why       string
	Stage     string
	Pass      string
}

// ReviewCoverageDtoV1 records what one parameter's review actually did —
// checked, not checked, or deliberately skipped with a reason. It exists so
// "no findings" can be told apart from "nothing looked".
type ReviewCoverageDtoV1 struct {
	Parameter string
	State     string
	Reason    string
}

// ReviewResultDtoV1 is one Review stage run as the runner reports it.
//
// MustFixOpen is the RUNNER's decision, computed from the findings and the
// thresholds stamped into the Job (parameters_enums.ReviewMustFixThresholds),
// never read from the agent — an agent that could assert it could wave its own
// findings through.
type ReviewResultDtoV1 struct {
	Findings    []ReviewFindingDtoV1
	Coverage    []ReviewCoverageDtoV1
	MustFixOpen bool
	AgentType   string
	Model       string
	TokenUsage  ReviewTokenUsageDtoV1
	CostUSD     float64
}

// ReviewTokenUsageDtoV1 mirrors agentbox's token_usage object. Cache writes
// bill separately from cache reads, so they are counted separately.
type ReviewTokenUsageDtoV1 struct {
	InputTokens         int
	OutputTokens        int
	CacheReadTokens     int
	CacheCreationTokens int
}

// UpdateTaskStepReviewDtoV1 carries the final Review stage result for one Task
// Step to deployment-server, which bounds it and stores it on the Step.
//
// Lives under tasks while the RPC that carries it (Jobs.UpdateTaskStepReviewV1)
// sits on the Jobs receiver — the same split UpdateTaskStepRunningDtoV1 and
// Jobs.MarkTaskStepRunningV1 already use.
type UpdateTaskStepReviewDtoV1 struct {
	TaskID    string
	StepIndex int
	JobID     string
	Result    ReviewResultDtoV1
}

type UpdateTaskStepReviewArgsV1 struct {
	types.AuthArgsV1
	Update UpdateTaskStepReviewDtoV1
}

type UpdateTaskStepReviewReplyV1 struct {
	Done bool
}

// UpdateTaskStepStageDtoV1 reports which stage of a Step Job is running right
// now, so a long Job reads as "Review" rather than an unexplained stretch of
// "Running".
//
// Stage is a string for the same reason the finding fields are: the runner
// cannot import kit's stage_enums. deployment-server maps it back.
// Deliberately its own RPC rather than a heartbeat field — a stage transition
// is an event with a definite order, and folding it into the ~5s heartbeat
// would make "Review started" arrive whenever the next beat happened to fire.
type UpdateTaskStepStageDtoV1 struct {
	TaskID    string
	StepIndex int
	JobID     string
	Stage     string
}

// The stage names the runner sends in UpdateTaskStepStageDtoV1.Stage. Part of
// the same wire mirror as the value functions below: the runner cannot name
// kit's stage_enums.Stage, so it spells the stage and kit parses it. kit owns
// a test pinning both names to the enum values they must resolve to.
const (
	StageImplement = "implement"
	StageReview    = "review"
)

type UpdateTaskStepStageArgsV1 struct {
	types.AuthArgsV1
	Update UpdateTaskStepStageDtoV1
}

type UpdateTaskStepStageReplyV1 struct {
	Done bool
}

// The WIRE MIRROR of kit's review_enums.
//
// The runner computes MustFixOpen — it is the component that holds both the
// findings and the thresholds — but it does NOT import
// github.com/deployment-io/kit and must not start: the runner ships to
// customers and depends only on this module. So the name→value mapping it
// needs lives here, hand-mirrored, and kit owns a test pinning these functions
// to the real enum values so the two numberings cannot drift.
//
// Same normalisation as kit's: case, spaces, hyphens and underscores are
// ignored, because the producer is a model writing prose and the same
// parameter comes back spelled four ways in four runs.
//
// Unknown input returns (0, false). A finding whose parameter or severity does
// not parse is ANNOTATED, never treated as must-fix: a threshold cannot be
// looked up for a parameter nobody can name.

// ReviewParameterValue maps a parameter name to kit's review_enums.Parameter
// value. Zero is never a finding's parameter — it is the all-parameters
// wildcard in a rule — so (0, false) is the only failure shape.
func ReviewParameterValue(s string) (uint, bool) {
	v, ok := reviewParameterValues[normalizeReviewKey(s)]
	return v, ok
}

// ReviewSeverityValue maps a severity name to kit's review_enums.Severity
// value. The ORDER is load-bearing: must-fix is "at or above", so 1..5 must
// stay Info < Low < Medium < High < Critical.
func ReviewSeverityValue(s string) (uint, bool) {
	v, ok := reviewSeverityValues[normalizeReviewKey(s)]
	return v, ok
}

// ReviewCoverageStateValue maps a coverage-state name to kit's
// review_enums.CoverageState value.
func ReviewCoverageStateValue(s string) (uint, bool) {
	v, ok := reviewCoverageStateValues[normalizeReviewKey(s)]
	return v, ok
}

var reviewParameterValues = map[string]uint{
	"security":        1,
	"correctness":     2,
	"specconformance": 3,
	"spec":            3,
	"conformance":     3,
	"testing":         4,
	"tests":           4,
	"test":            4,
	"deployreadiness": 5,
	"deploy":          5,
	"deployment":      5,
	"performance":     6,
	"perf":            6,
	"maintainability": 7,
	"reliability":     8,
	"resilience":      8,
	"robustness":      8,
}

var reviewSeverityValues = map[string]uint{
	"info":          1,
	"informational": 1,
	"note":          1,
	"low":           2,
	"medium":        3,
	"med":           3,
	"moderate":      3,
	"high":          4,
	"critical":      5,
	"crit":          5,
	"blocker":       5,
}

var reviewCoverageStateValues = map[string]uint{
	"checked":    1,
	"notchecked": 2,
	"unchecked":  2,
	"notrun":     2,
	"skipped":    3,
	"skip":       3,
}

// normalizeReviewKey folds untrusted agent text to a comparison key:
// lowercase, with everything that is not a letter or digit removed. Mirrors
// review_enums.normalizeKey exactly.
func normalizeReviewKey(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range strings.ToLower(strings.TrimSpace(s)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}
