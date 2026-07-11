// Package task_previews holds the RPC DTOs for the agent-driven ephemeral preview
// lifecycle (C4). A task preview is one ephemeral Environment (kit#194) that scopes
// a whole stack of services — each a lean Deployment under that env. The runner's
// preview tools lazily ask deployment-server to find-or-create the env + the
// Deployment for a given service and hand back the ids needed to deploy it.
// Distinct from the `previews` package, which is the deployment-dependent PR-preview
// model.
package task_previews

import "github.com/deployment-io/deployment-runner-kit/types"

// Service-type tokens for EnsureTaskPreviewArgsV1.ServiceType. Kept as plain strings
// so this DTO package doesn't import kit's deployments enum (that would be a
// drk→kit import cycle); deployment-server maps these to deployments.Type. A task's
// preview env can hold one of each (a static SPA + a public API + a private DB).
const (
	ServiceTypeStaticSite     = "static-site"
	ServiceTypeWebService     = "web-service"
	ServiceTypePrivateService = "private-service"
	ServiceTypeRDS            = "rds"
)

// EnsureTaskPreviewArgsV1 asks deployment-server to find-or-create the task's
// ephemeral preview Environment (extending its idle TTL on reuse) and, within it,
// the lean Deployment for ONE service — returning the ids the runner needs to
// deploy that service. RunnerID + Region come from the runner's own identity so the
// env is pinned to the runner that serves it.
type EnsureTaskPreviewArgsV1 struct {
	types.AuthArgsV1
	TaskID string

	// The calling runner's identity — the same fields it sends when polling for jobs.
	// deployment-server computes the RunnerName from these + the token's user and pins
	// the ephemeral env to that runner, so GC dispatch + the region constraint line up.
	CloudAccountID string
	Region         string // runner region, region_enums string e.g. "us-east-1"
	GoArch         string // runtime.GOARCH
	GoOS           string // runtime.GOOS

	// ServiceName identifies this service within the task's preview env (a task can
	// preview a full stack — static + web + db — each its own Deployment under the
	// one env). Must be unique per task; used to name/find the Deployment.
	ServiceName string
	// ServiceType is one of the ServiceType* tokens above.
	ServiceType string
}

// EnsureTaskPreviewReplyV1 is the resolved identity of THIS service's Deployment.
// PreviewID is its id (used for bucket + CloudFront resource naming, matching the
// <org>-<deploymentID> production scheme). ExistingDistID is the CloudFront
// distribution id persisted from a prior deploy — empty on the first (create), set
// on reuse (re-upload into the existing one). Web/DB service types will extend this
// with their own resource handles as those tools land.
type EnsureTaskPreviewReplyV1 struct {
	PreviewID      string
	ExistingDistID string
}
