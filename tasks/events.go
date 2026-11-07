package tasks

import (
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/deployment-io/deployment-runner-kit/types"
)

// Structured progress events for a Task run.
//
// The runner sends them while a Task Step Job runs (Jobs.AppendTaskEventsV1);
// deployment-server bounds them (Bound), numbers them (Seq) and appends them
// to the run's Job in kit (Job.Events); app-server returns them as-is to the
// Task page. Because kit stores this struct as-is and app-server returns it
// as-is, every field carries json and bson tags, as LiveProgressV1 does.
//
// The type names, stage names and caps below are owned here; kit,
// deployment-server and app-server import them and never define their own.

// Event types. Exactly the payload named next to each is set.
const (
	TaskEventStageStarted          = "stage_started"           // StageDetail
	TaskEventStageFinished         = "stage_finished"          // StageDetail
	TaskEventCheckoutDone          = "checkout_done"           // Checkout
	TaskEventSelfCheck             = "self_check"              // SelfCheck
	TaskEventImplementFinished     = "implement_finished"      // Implement
	TaskEventReviewRoundFinished   = "review_round_finished"   // ReviewRound
	TaskEventFixRoundFinished      = "fix_round_finished"      // FixRound
	TaskEventPreviewDeployFinished = "preview_deploy_finished" // PreviewDeploy
	TaskEventQACheckFinished       = "qa_check_finished"       // QACheck
	TaskEventQACriterionFinished   = "qa_criterion_finished"   // QACriterion
	TaskEventQARoundFinished       = "qa_round_finished"       // QARound
	TaskEventPROpened              = "pr_opened"               // PullRequest
)

// Event stage names (TaskEventV1.Stage).
const (
	TaskEventStageSetup     = "setup"
	TaskEventStageImplement = StageImplement
	TaskEventStageReview    = StageReview
	TaskEventStageQA        = StageQA
	TaskEventStagePR        = "pr"
)

// Preview deploy kinds (TaskEventPreviewDeployV1.Kind).
const (
	TaskEventPreviewKindSite    = "site"
	TaskEventPreviewKindBackend = "backend"
)

// Earlier findings' status (TaskEventFindingStatusV1.Status).
const (
	TaskEventFindingResolved     = "resolved"
	TaskEventFindingStillPresent = "still_present"
)

// MaxTaskEventsPerBatch is the most events one AppendTaskEventsArgsV1 may
// carry; deployment-server rejects a larger batch.
const MaxTaskEventsPerBatch = 100

// Caps applied by Bound. String caps are in runes unless named Bytes.
const (
	TaskEventShortTextRunes      = 200
	TaskEventMediumTextRunes     = 500
	TaskEventLongTextRunes       = 1000
	TaskEventChangesSummaryBytes = 8 * 1024
	TaskEventURLBytes            = 2048
	TaskEventMaxCommands         = 10
	TaskEventMaxListEntries      = 20 // Repos, Roles, Addressed
	TaskEventMaxFindings         = 100
	TaskEventMaxCoverage         = 20
	TaskEventMaxDiffstat         = 20
)

// TaskEventV1 is the envelope of one progress event.
//
// Seq is assigned by deployment-server (1, 2, 3… per Job); the runner sends
// 0. Ts is unix milliseconds on the runner's clock; Seq, not Ts, orders
// events. Round is 0 when the event isn't in a round. Repo is the
// repository's full name (owner/name) for an event that belongs to one
// repository, else "". Now is the one-line text the board and the page's
// "now" line show, or "" for an event that doesn't change it.
type TaskEventV1 struct {
	Seq   int    `json:"seq" bson:"seq"`
	Type  string `json:"type" bson:"type"`
	Ts    int64  `json:"ts" bson:"ts"`
	Stage string `json:"stage" bson:"stage"`
	Round int    `json:"round" bson:"round"`
	Repo  string `json:"repo" bson:"repo"`
	Now   string `json:"now" bson:"now"`

	// StageDetail is the payload of stage_started and stage_finished
	// (named StageDetail, not Stage, because Stage is the envelope's stage).
	StageDetail   *TaskEventStageV1         `json:"stageDetail,omitempty" bson:"stageDetail,omitempty"`
	Checkout      *TaskEventCheckoutV1      `json:"checkout,omitempty" bson:"checkout,omitempty"`
	SelfCheck     *TaskEventSelfCheckV1     `json:"selfCheck,omitempty" bson:"selfCheck,omitempty"`
	Implement     *TaskEventImplementV1     `json:"implement,omitempty" bson:"implement,omitempty"`
	ReviewRound   *TaskEventReviewRoundV1   `json:"reviewRound,omitempty" bson:"reviewRound,omitempty"`
	FixRound      *TaskEventFixRoundV1      `json:"fixRound,omitempty" bson:"fixRound,omitempty"`
	PreviewDeploy *TaskEventPreviewDeployV1 `json:"previewDeploy,omitempty" bson:"previewDeploy,omitempty"`
	QACheck       *TaskEventQACheckV1       `json:"qaCheck,omitempty" bson:"qaCheck,omitempty"`
	QACriterion   *TaskEventQACriterionV1   `json:"qaCriterion,omitempty" bson:"qaCriterion,omitempty"`
	QARound       *TaskEventQARoundV1       `json:"qaRound,omitempty" bson:"qaRound,omitempty"`
	PullRequest   *TaskEventPullRequestV1   `json:"pullRequest,omitempty" bson:"pullRequest,omitempty"`
}

// TaskEventStageV1: Outcome is set on stage_finished only (e.g. "No
// findings", "Must-fix open").
type TaskEventStageV1 struct {
	Outcome string `json:"outcome" bson:"outcome"`
}

type TaskEventCheckoutV1 struct {
	Repos []string `json:"repos" bson:"repos"`
}

// TaskEventSelfCheckV1 is the implementer's build/test self-check in one
// repository (the envelope's Repo).
type TaskEventSelfCheckV1 struct {
	Commands []string `json:"commands" bson:"commands"`
	Passed   bool     `json:"passed" bson:"passed"`
	Summary  string   `json:"summary" bson:"summary"`
}

type TaskEventImplementV1 struct {
	Turns          int                   `json:"turns" bson:"turns"`
	InputTokens    int                   `json:"inputTokens" bson:"inputTokens"`
	OutputTokens   int                   `json:"outputTokens" bson:"outputTokens"`
	ChangesSummary string                `json:"changesSummary" bson:"changesSummary"`
	Diffstat       []TaskEventDiffstatV1 `json:"diffstat" bson:"diffstat"`
}

type TaskEventDiffstatV1 struct {
	Repo       string `json:"repo" bson:"repo"`
	Files      int    `json:"files" bson:"files"`
	Insertions int    `json:"insertions" bson:"insertions"`
	Deletions  int    `json:"deletions" bson:"deletions"`
}

// TaskEventReviewRoundV1: Findings are the findings new in this round only;
// Previous is the status of earlier rounds' findings.
type TaskEventReviewRoundV1 struct {
	Findings []TaskEventFindingV1       `json:"findings" bson:"findings"`
	Previous []TaskEventFindingStatusV1 `json:"previous" bson:"previous"`
	Coverage []TaskEventCoverageV1      `json:"coverage" bson:"coverage"`
}

type TaskEventFindingV1 struct {
	Key       string `json:"key" bson:"key"`
	Parameter string `json:"parameter" bson:"parameter"`
	Severity  string `json:"severity" bson:"severity"`
	Location  string `json:"location" bson:"location"`
	What      string `json:"what" bson:"what"`
	Why       string `json:"why" bson:"why"`
	MustFix   bool   `json:"mustFix" bson:"mustFix"`
}

// TaskEventFindingStatusV1: Status is TaskEventFindingResolved or
// TaskEventFindingStillPresent.
type TaskEventFindingStatusV1 struct {
	Key    string `json:"key" bson:"key"`
	Status string `json:"status" bson:"status"`
}

type TaskEventCoverageV1 struct {
	Parameter string `json:"parameter" bson:"parameter"`
	State     string `json:"state" bson:"state"`
	Reason    string `json:"reason" bson:"reason"`
}

// TaskEventFixRoundV1: Addressed holds finding keys.
type TaskEventFixRoundV1 struct {
	Addressed      []string              `json:"addressed" bson:"addressed"`
	ChangesSummary string                `json:"changesSummary" bson:"changesSummary"`
	SelfCheck      *TaskEventSelfCheckV1 `json:"selfCheck,omitempty" bson:"selfCheck,omitempty"`
}

// TaskEventPreviewDeployV1 is one preview deploy. Never the migrate or seed
// command's text or its full output: Summary on each is one line.
type TaskEventPreviewDeployV1 struct {
	Kind           string              `json:"kind" bson:"kind"`
	URL            string              `json:"url" bson:"url"`
	OK             bool                `json:"ok" bson:"ok"`
	Error          string              `json:"error" bson:"error"`
	DatabaseEngine string              `json:"databaseEngine" bson:"databaseEngine"`
	Migrate        *TaskEventCommandV1 `json:"migrate,omitempty" bson:"migrate,omitempty"`
	Seed           *TaskEventCommandV1 `json:"seed,omitempty" bson:"seed,omitempty"`
	Roles          []string            `json:"roles" bson:"roles"`
}

type TaskEventCommandV1 struct {
	Ran      bool   `json:"ran" bson:"ran"`
	ExitCode int    `json:"exitCode" bson:"exitCode"`
	Summary  string `json:"summary" bson:"summary"`
}

// TaskEventQACheckV1 is one check the QA agent ran. Criterion is a 0-based
// index into the criteria. Path never has a query string.
type TaskEventQACheckV1 struct {
	Criterion   int    `json:"criterion" bson:"criterion"`
	Tool        string `json:"tool" bson:"tool"`
	Method      string `json:"method" bson:"method"`
	Path        string `json:"path" bson:"path"`
	Role        string `json:"role" bson:"role"`
	Status      int    `json:"status" bson:"status"`
	DurationMs  int    `json:"durationMs" bson:"durationMs"`
	OK          bool   `json:"ok" bson:"ok"`
	Result      string `json:"result" bson:"result"`
	EvidenceURL string `json:"evidenceUrl" bson:"evidenceUrl"`
}

// TaskEventQACriterionV1: State is one of the QACriterion* names (qa.go).
type TaskEventQACriterionV1 struct {
	Criterion int    `json:"criterion" bson:"criterion"`
	Text      string `json:"text" bson:"text"`
	State     string `json:"state" bson:"state"`
	Reason    string `json:"reason" bson:"reason"`
}

type TaskEventQARoundV1 struct {
	Final      bool   `json:"final" bson:"final"`
	Passed     int    `json:"passed" bson:"passed"`
	Failed     int    `json:"failed" bson:"failed"`
	NotChecked int    `json:"notChecked" bson:"notChecked"`
	Summary    string `json:"summary" bson:"summary"`
}

// TaskEventPullRequestV1: Updated is true when an existing PR was updated
// rather than opened.
type TaskEventPullRequestV1 struct {
	Number         int    `json:"number" bson:"number"`
	URL            string `json:"url" bson:"url"`
	Title          string `json:"title" bson:"title"`
	Draft          bool   `json:"draft" bson:"draft"`
	Updated        bool   `json:"updated" bson:"updated"`
	InlineComments int    `json:"inlineComments" bson:"inlineComments"`
}

// AppendTaskEventsArgsV1 is one batch of events for one run (Job) of a Task.
// BatchID is a random id the runner makes per batch and reuses when it
// retries that batch, so a retry after a lost reply stores nothing twice.
// At most MaxTaskEventsPerBatch events.
type AppendTaskEventsArgsV1 struct {
	types.AuthArgsV1
	TaskID  string
	JobID   string
	BatchID string
	Events  []TaskEventV1
}

type AppendTaskEventsReplyV1 struct {
	Done bool
}

var taskEventStages = map[string]bool{
	TaskEventStageSetup:     true,
	TaskEventStageImplement: true,
	TaskEventStageReview:    true,
	TaskEventStageQA:        true,
	TaskEventStagePR:        true,
}

var taskEventVerdicts = map[string]bool{
	TaskEventStageFinished:         true,
	TaskEventImplementFinished:     true,
	TaskEventReviewRoundFinished:   true,
	TaskEventFixRoundFinished:      true,
	TaskEventPreviewDeployFinished: true,
	TaskEventQACriterionFinished:   true,
	TaskEventQARoundFinished:       true,
	TaskEventPROpened:              true,
}

// IsVerdict reports whether the event is a verdict — the outcome of a
// stage, round, criterion, deploy or PR. Storage never drops a verdict for
// the soft caps (only the hard ones); the others are intermediate.
func (e *TaskEventV1) IsVerdict() bool {
	return taskEventVerdicts[e.Type]
}

// payloads returns how many payloads are set and whether the one Type names
// is among them; ok is false for an unknown Type.
func (e *TaskEventV1) payloads() (set int, matching bool, ok bool) {
	var want bool
	ok = true
	count := func(isSet bool, forTypes ...string) {
		if !isSet {
			return
		}
		set++
		for _, t := range forTypes {
			if t == e.Type {
				matching = true
			}
		}
	}
	switch e.Type {
	case TaskEventStageStarted, TaskEventStageFinished, TaskEventCheckoutDone,
		TaskEventSelfCheck, TaskEventImplementFinished, TaskEventReviewRoundFinished,
		TaskEventFixRoundFinished, TaskEventPreviewDeployFinished,
		TaskEventQACheckFinished, TaskEventQACriterionFinished,
		TaskEventQARoundFinished, TaskEventPROpened:
		want = true
	}
	if !want {
		return 0, false, false
	}
	count(e.StageDetail != nil, TaskEventStageStarted, TaskEventStageFinished)
	count(e.Checkout != nil, TaskEventCheckoutDone)
	count(e.SelfCheck != nil, TaskEventSelfCheck)
	count(e.Implement != nil, TaskEventImplementFinished)
	count(e.ReviewRound != nil, TaskEventReviewRoundFinished)
	count(e.FixRound != nil, TaskEventFixRoundFinished)
	count(e.PreviewDeploy != nil, TaskEventPreviewDeployFinished)
	count(e.QACheck != nil, TaskEventQACheckFinished)
	count(e.QACriterion != nil, TaskEventQACriterionFinished)
	count(e.QARound != nil, TaskEventQARoundFinished)
	count(e.PullRequest != nil, TaskEventPROpened)
	return set, matching, ok
}

// Bound validates the event and caps it in place. It returns false — the
// event must be dropped — when Type is unknown, Stage is not one of the five
// stage names, the payload Type names is nil, or any other payload is set.
// Otherwise it caps strings (cut on a rune boundary, with "…" appended),
// lists and URLs (kept only as https URLs without userinfo, else ""),
// strips a query string from QACheck.Path, and turns negative ints into 0.
func (e *TaskEventV1) Bound() bool {
	set, matching, ok := e.payloads()
	if !ok || !matching || set != 1 {
		return false
	}
	if !taskEventStages[e.Stage] {
		return false
	}

	e.Seq = nonNegative(e.Seq)
	if e.Ts < 0 {
		e.Ts = 0
	}
	e.Round = nonNegative(e.Round)
	e.Repo = capRunes(e.Repo, TaskEventShortTextRunes)
	e.Now = capRunes(e.Now, TaskEventShortTextRunes)

	switch {
	case e.StageDetail != nil:
		p := e.StageDetail
		p.Outcome = capRunes(p.Outcome, TaskEventShortTextRunes)
	case e.Checkout != nil:
		p := e.Checkout
		p.Repos = capStrings(p.Repos, TaskEventMaxListEntries, TaskEventShortTextRunes)
	case e.SelfCheck != nil:
		boundSelfCheck(e.SelfCheck)
	case e.Implement != nil:
		p := e.Implement
		p.Turns = nonNegative(p.Turns)
		p.InputTokens = nonNegative(p.InputTokens)
		p.OutputTokens = nonNegative(p.OutputTokens)
		p.ChangesSummary = capBytes(p.ChangesSummary, TaskEventChangesSummaryBytes)
		if len(p.Diffstat) > TaskEventMaxDiffstat {
			p.Diffstat = p.Diffstat[:TaskEventMaxDiffstat]
		}
		for i := range p.Diffstat {
			d := &p.Diffstat[i]
			d.Repo = capRunes(d.Repo, TaskEventShortTextRunes)
			d.Files = nonNegative(d.Files)
			d.Insertions = nonNegative(d.Insertions)
			d.Deletions = nonNegative(d.Deletions)
		}
	case e.ReviewRound != nil:
		p := e.ReviewRound
		if len(p.Findings) > TaskEventMaxFindings {
			p.Findings = p.Findings[:TaskEventMaxFindings]
		}
		for i := range p.Findings {
			f := &p.Findings[i]
			f.Key = capRunes(f.Key, TaskEventShortTextRunes)
			f.Parameter = capRunes(f.Parameter, TaskEventShortTextRunes)
			f.Severity = capRunes(f.Severity, TaskEventShortTextRunes)
			f.Location = capRunes(f.Location, TaskEventMediumTextRunes)
			f.What = capRunes(f.What, TaskEventLongTextRunes)
			f.Why = capRunes(f.Why, TaskEventLongTextRunes)
		}
		if len(p.Previous) > TaskEventMaxFindings {
			p.Previous = p.Previous[:TaskEventMaxFindings]
		}
		for i := range p.Previous {
			s := &p.Previous[i]
			s.Key = capRunes(s.Key, TaskEventShortTextRunes)
			s.Status = capRunes(s.Status, TaskEventShortTextRunes)
		}
		if len(p.Coverage) > TaskEventMaxCoverage {
			p.Coverage = p.Coverage[:TaskEventMaxCoverage]
		}
		for i := range p.Coverage {
			c := &p.Coverage[i]
			c.Parameter = capRunes(c.Parameter, TaskEventShortTextRunes)
			c.State = capRunes(c.State, TaskEventShortTextRunes)
			c.Reason = capRunes(c.Reason, TaskEventShortTextRunes)
		}
	case e.FixRound != nil:
		p := e.FixRound
		p.Addressed = capStrings(p.Addressed, TaskEventMaxListEntries, TaskEventShortTextRunes)
		p.ChangesSummary = capBytes(p.ChangesSummary, TaskEventChangesSummaryBytes)
		if p.SelfCheck != nil {
			boundSelfCheck(p.SelfCheck)
		}
	case e.PreviewDeploy != nil:
		p := e.PreviewDeploy
		p.Kind = capRunes(p.Kind, TaskEventShortTextRunes)
		p.URL = boundURL(p.URL)
		p.Error = capRunes(p.Error, TaskEventShortTextRunes)
		p.DatabaseEngine = capRunes(p.DatabaseEngine, TaskEventShortTextRunes)
		boundCommand(p.Migrate)
		boundCommand(p.Seed)
		p.Roles = capStrings(p.Roles, TaskEventMaxListEntries, TaskEventShortTextRunes)
	case e.QACheck != nil:
		p := e.QACheck
		p.Criterion = nonNegative(p.Criterion)
		p.Tool = capRunes(p.Tool, TaskEventShortTextRunes)
		p.Method = capRunes(p.Method, TaskEventShortTextRunes)
		if i := strings.IndexAny(p.Path, "?#"); i >= 0 {
			p.Path = p.Path[:i]
		}
		p.Path = capRunes(p.Path, TaskEventMediumTextRunes)
		p.Role = capRunes(p.Role, TaskEventShortTextRunes)
		p.Status = nonNegative(p.Status)
		p.DurationMs = nonNegative(p.DurationMs)
		p.Result = capRunes(p.Result, TaskEventShortTextRunes)
		p.EvidenceURL = boundURL(p.EvidenceURL)
	case e.QACriterion != nil:
		p := e.QACriterion
		p.Criterion = nonNegative(p.Criterion)
		p.Text = capRunes(p.Text, TaskEventMediumTextRunes)
		p.State = capRunes(p.State, TaskEventShortTextRunes)
		p.Reason = capRunes(p.Reason, TaskEventLongTextRunes)
	case e.QARound != nil:
		p := e.QARound
		p.Passed = nonNegative(p.Passed)
		p.Failed = nonNegative(p.Failed)
		p.NotChecked = nonNegative(p.NotChecked)
		p.Summary = capRunes(p.Summary, TaskEventShortTextRunes)
	case e.PullRequest != nil:
		p := e.PullRequest
		p.Number = nonNegative(p.Number)
		p.URL = boundURL(p.URL)
		p.Title = capRunes(p.Title, TaskEventShortTextRunes)
		p.InlineComments = nonNegative(p.InlineComments)
	}
	return true
}

func boundSelfCheck(p *TaskEventSelfCheckV1) {
	p.Commands = capStrings(p.Commands, TaskEventMaxCommands, TaskEventShortTextRunes)
	p.Summary = capRunes(p.Summary, TaskEventShortTextRunes)
}

func boundCommand(p *TaskEventCommandV1) {
	if p == nil {
		return
	}
	p.ExitCode = nonNegative(p.ExitCode)
	p.Summary = capRunes(p.Summary, TaskEventShortTextRunes)
}

const ellipsis = "…"

// capRunes returns s when it has at most n runes, else its first n-1 runes
// followed by "…" (n runes in all).
func capRunes(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	i, count := 0, 0
	for i < len(s) && count < n-1 {
		_, size := utf8.DecodeRuneInString(s[i:])
		i += size
		count++
	}
	return s[:i] + ellipsis
}

// capBytes returns s when it is at most n bytes, else a prefix cut on a rune
// boundary followed by "…", at most n bytes in all.
func capBytes(s string, n int) string {
	if len(s) <= n {
		return s
	}
	i := n - len(ellipsis)
	for i > 0 && !utf8.RuneStart(s[i]) {
		i--
	}
	return s[:i] + ellipsis
}

func capStrings(list []string, entries, runes int) []string {
	if len(list) > entries {
		list = list[:entries]
	}
	for i := range list {
		list[i] = capRunes(list[i], runes)
	}
	return list
}

// boundURL caps u at TaskEventURLBytes and keeps it only when it then parses
// as an absolute https URL with a host and no userinfo.
func boundURL(u string) string {
	if u == "" {
		return ""
	}
	u = capBytes(u, TaskEventURLBytes)
	parsed, err := url.Parse(u)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil {
		return ""
	}
	return u
}

func nonNegative(n int) int {
	if n < 0 {
		return 0
	}
	return n
}
