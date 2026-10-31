package tasks

import "github.com/deployment-io/deployment-runner-kit/types"

// ProductionWatchWindowSeconds is how long a production watch Job watches a
// Task's production after its deploy.
const ProductionWatchWindowSeconds = 1800

// ProductionWatchTargetV1 is one deployed repository a production watch Job
// watches, stamped at Job creation as a JSON array in the
// ProductionWatchTargets parameter. Empty AWS fields are the ones the Kind
// does not have (a static site has no load balancer, a private service no
// URL).
type ProductionWatchTargetV1 struct {
	// Repository is "<index>-<owner/name>".
	Repository   string `json:"repository"`
	DeploymentID string `json:"deploymentId"`
	// Kind is one of the ProductionWatchKind* values below.
	Kind string `json:"kind"`
	// Region is a region_enums.Type.
	Region                   int64  `json:"region"`
	LoadBalancerArn          string `json:"loadBalancerArn,omitempty"`
	TargetGroupArn           string `json:"targetGroupArn,omitempty"`
	EcsServiceArn            string `json:"ecsServiceArn,omitempty"`
	ClusterName              string `json:"clusterName,omitempty"`
	CloudfrontDistributionID string `json:"cloudfrontDistributionId,omitempty"`
	// URLs are the production URL and, for a web service, its health URL.
	URLs []string `json:"urls,omitempty"`
	// DeployStartedAt and DeployEndedAt are the delivery's StartedAt and
	// EndedAt, unix seconds.
	DeployStartedAt int64 `json:"deployStartedAt"`
	DeployEndedAt   int64 `json:"deployEndedAt"`
	// BuildID is the delivered build.
	BuildID string `json:"buildId"`
}

// The target kinds in ProductionWatchTargetV1.Kind.
const (
	ProductionWatchKindWebService     = "web_service"
	ProductionWatchKindPrivateService = "private_service"
	ProductionWatchKindStaticSite     = "static_site"
)

// The WIRE values of a production watch report. kit owns a test pinning each
// to its own value; deployment-server rejects a report holding any other.
const (
	// Report states, ProductionWatchReportDtoV1.State.
	ProductionWatchStateRunning = "running"
	ProductionWatchStateDone    = "done"

	// Verdicts, ProductionWatchReportDtoV1.Verdict ("" while running).
	ProductionWatchVerdictNoChange         = "no_change"
	ProductionWatchVerdictErrorsUp         = "errors_up"
	ProductionWatchVerdictNotEnoughTraffic = "not_enough_traffic"

	// Signal kinds, ProductionWatchSignalDtoV1.Kind.
	ProductionWatchSignalUnhealthyTargets = "unhealthy_targets"
	ProductionWatchSignalTaskCrashed      = "task_crashed"
	ProductionWatchSignalURLFailing       = "url_failing"
	ProductionWatchSignalErrorsUp         = "errors_up"
)

// ProductionWatchTargetReportDtoV1 is one target's numbers at a poll. Before
// is the window of the same length before the deploy, After the time since
// its end. P95Ms is 0 when unknown.
type ProductionWatchTargetReportDtoV1 struct {
	Repository      string
	RequestsBefore  int
	RequestsAfter   int
	Errors5xxBefore int
	Errors5xxAfter  int
	P95MsBefore     int
	P95MsAfter      int
	// UnhealthyTargets is the count at the latest poll.
	UnhealthyTargets int
	// CrashedTasks is the count since the deploy.
	CrashedTasks int
	// FailingURLs are at most 5.
	FailingURLs []string
}

// ProductionWatchSignalDtoV1 is one thing the watch saw go wrong. Detail is
// one line and holds no secret value. SeenAt is unix seconds.
type ProductionWatchSignalDtoV1 struct {
	Repository string
	Kind       string
	Detail     string
	SeenAt     int64
}

// ProductionWatchReportDtoV1 is one poll of a production watch. PolledAt and
// WindowEndsAt are unix seconds.
type ProductionWatchReportDtoV1 struct {
	State        string
	Verdict      string
	Targets      []ProductionWatchTargetReportDtoV1
	Signals      []ProductionWatchSignalDtoV1
	PolledAt     int64
	WindowEndsAt int64
}

// UpdateTaskProductionWatchDtoV1 carries one poll of a production watch Job
// to deployment-server (Jobs.UpdateTaskProductionWatchV1), which bounds it
// and stores it on the Task.
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
// now — the Task's watch already ended (timed out, or the deployment was
// deployed again or rolled back) or the Task is no longer Succeeded.
type UpdateTaskProductionWatchReplyV1 struct {
	Done bool
	Stop bool
}
