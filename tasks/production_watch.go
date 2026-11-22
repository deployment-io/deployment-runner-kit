package tasks

import (
	"github.com/deployment-io/deployment-runner-kit/enums/region_enums"
	"github.com/deployment-io/deployment-runner-kit/types"
)

// ProductionWatchWindowSeconds is how long a production watch Job watches a
// deployed Task's production after its deploy: 30 minutes.
const ProductionWatchWindowSeconds = 1800

// ProductionWatchTargetV1 is one deployed repository a production watch Job
// watches. app-server stamps the Job's targets AT JOB CREATION as a JSON
// array in the ProductionWatchTargets parameter; the runner only reads them.
//
// Kind is one of the ProductionWatchKind* names below, set by app-server and
// only read by the runner. The AWS identifiers are empty when the kind has
// none (a static site has no load balancer or ECS service; a private service
// no load balancer). URLs are the production URL and, for a web service, its
// health URL. DeployStartedAt and DeployEndedAt are the delivery's StartedAt
// and EndedAt, in unix seconds; BuildID is the delivered build.
type ProductionWatchTargetV1 struct {
	// Repository is "<index>-<owner/name>".
	Repository               string            `json:"repository"`
	DeploymentID             string            `json:"deploymentId"`
	Kind                     string            `json:"kind"`
	Region                   region_enums.Type `json:"region"`
	LoadBalancerArn          string            `json:"loadBalancerArn,omitempty"`
	TargetGroupArn           string            `json:"targetGroupArn,omitempty"`
	EcsServiceArn            string            `json:"ecsServiceArn,omitempty"`
	ClusterName              string            `json:"clusterName,omitempty"`
	CloudfrontDistributionID string            `json:"cloudfrontDistributionId,omitempty"`
	URLs                     []string          `json:"urls,omitempty"`
	DeployStartedAt          int64             `json:"deployStartedAt"`
	DeployEndedAt            int64             `json:"deployEndedAt"`
	BuildID                  string            `json:"buildId"`
}

// Target kinds, set by app-server in ProductionWatchTargetV1.Kind.
const (
	ProductionWatchKindWebService     = "web_service"
	ProductionWatchKindPrivateService = "private_service"
	ProductionWatchKindStaticSite     = "static_site"
)

// The WIRE names the runner sends in a ProductionWatchReportDtoV1. The runner
// sends only these State, Verdict and signal Kind values; deployment-server
// validates them and kit stores the same strings. kit's other watch states
// (couldnt_finish, ended_early) are written by app-server, never sent by a
// runner.
const (
	ProductionWatchStateRunning = "running"
	ProductionWatchStateDone    = "done"
)

// Verdicts, sent only with ProductionWatchStateDone.
const (
	ProductionWatchVerdictNoChange         = "no_change"
	ProductionWatchVerdictErrorsUp         = "errors_up"
	ProductionWatchVerdictNotEnoughTraffic = "not_enough_traffic"
)

// Signal kinds, in ProductionWatchSignalDtoV1.Kind.
const (
	ProductionWatchSignalUnhealthyTargets = "unhealthy_targets"
	ProductionWatchSignalTaskCrashed      = "task_crashed"
	ProductionWatchSignalURLFailing       = "url_failing"
	ProductionWatchSignalErrorsUp         = "errors_up"
)

// ProductionWatchTargetReportDtoV1 is one target's numbers at a poll,
// CUMULATIVE since the deploy. Before is the window of the same length before
// the deploy, After the time since it ended. P95Ms is 0 when unknown.
// UnhealthyTargets is the count at this poll, CrashedTasks the count since the
// deploy. FailingURLs are the target's URLs failing at this poll.
type ProductionWatchTargetReportDtoV1 struct {
	Repository       string
	RequestsBefore   int
	RequestsAfter    int
	Errors5xxBefore  int
	Errors5xxAfter   int
	P95MsBefore      int
	P95MsAfter       int
	UnhealthyTargets int
	CrashedTasks     int
	FailingURLs      []string
}

// ProductionWatchSignalDtoV1 is one thing that got worse since the deploy.
// Kind is one of the ProductionWatchSignal* names. Detail is one line and
// holds no secret value. SeenAt is when it was first seen, in unix seconds.
type ProductionWatchSignalDtoV1 struct {
	Repository string
	Kind       string
	Detail     string
	SeenAt     int64
}

// ProductionWatchReportDtoV1 is one poll of a production watch Job.
//
// Every report is CUMULATIVE: Targets are the numbers since the deploy, and
// Signals are every signal the Job has seen so far (the first of each
// Repository and Kind), so a report the server drops loses nothing the next
// one doesn't carry.
//
// State is ProductionWatchStateRunning or ProductionWatchStateDone; Verdict is
// set only with done. PolledAt and WindowEndsAt are unix seconds.
type ProductionWatchReportDtoV1 struct {
	State        string
	Verdict      string
	Targets      []ProductionWatchTargetReportDtoV1
	Signals      []ProductionWatchSignalDtoV1
	PolledAt     int64
	WindowEndsAt int64
}

// UpdateTaskProductionWatchDtoV1 carries one poll of a production watch Job
// to deployment-server (Jobs.UpdateTaskProductionWatchV1), which validates
// it, bounds it and stores it on the Task.
type UpdateTaskProductionWatchDtoV1 struct {
	TaskID string
	JobID  string
	Report ProductionWatchReportDtoV1
}

type UpdateTaskProductionWatchArgsV1 struct {
	types.AuthArgsV1
	Update UpdateTaskProductionWatchDtoV1
}

// UpdateTaskProductionWatchReplyV1: Stop tells the runner to end its watch
// now (the Task or its watch is gone, over, or no longer this Job's).
type UpdateTaskProductionWatchReplyV1 struct {
	Done bool
	Stop bool
}
