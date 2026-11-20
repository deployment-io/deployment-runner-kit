package task_previews

import (
	"github.com/deployment-io/deployment-runner-kit/tasks"
	"github.com/deployment-io/deployment-runner-kit/types"
)

// CompanionBackendsArgsV1 asks deployment-server which web services outside the
// Task the Task's static-site previews call in production — its "companion
// backends" — so the runner can deploy a preview of each from its production
// branch and build the site previews to call them. deployment-server works it out
// from a person's links (a static site's PreviewBackends) and by matching
// production values, and answers with variable NAMES only, never values.
type CompanionBackendsArgsV1 struct {
	types.AuthArgsV1
	TaskID string
	// Repositories are every repository of the Task, as in StackWiringArgsV1; a
	// CloneURL that isn't the Task's is refused.
	Repositories []StackRepositoryV1
	// Frontends are the static-site deployments the runner builds the Task's site
	// previews from. One that deployment.io doesn't deploy is skipped.
	Frontends []CompanionFrontendV1
	Region    string // the runner's region, region_enums string, e.g. "eu-west-1"
}

// CompanionFrontendV1 is one static-site deployment of the Task, as the runner
// knows it from StaticSiteBuildSettingsV1.
type CompanionFrontendV1 struct {
	// Repository is a Repositories[i].Repository; another is refused.
	Repository string
	// Frontend is the static-site deployment's name
	// (StaticSiteBuildSettingsV1.DeploymentName).
	Frontend string
	// Environment is its environment's name
	// (StaticSiteBuildSettingsV1.EnvironmentName); names aren't unique on their own.
	Environment string
}

// CompanionBackendsReplyV1 holds the companions to deploy and the variables that
// would have linked to a web service outside the Task but won't. A frontend
// variable in neither calls nothing deployment.io deploys outside the Task (or is
// StackWiringV1's: it calls one of the Task's own APIs).
type CompanionBackendsReplyV1 struct {
	// Companions are at most MaxCompanionBackends, sorted by the number of linked
	// variables (most first), then API, then DeploymentID.
	Companions []CompanionBackendV1
	// NotDeployed are sorted by FrontendRepository, Frontend, Variable.
	NotDeployed []CompanionNotDeployedV1
}

// CompanionBackendV1 is one web service outside the Task that the Task's site
// previews call: the runner deploys a preview of it from Repository.BaseBranch.
type CompanionBackendV1 struct {
	// DeploymentID is the web service's Deployment ID, hex.
	DeploymentID string
	// API is the web-service deployment's name.
	API string
	// Repository is the web service's repository: Name, CloneURL, Provider and
	// InstallationID from the deployment's repository, BaseBranch the branch the
	// deployment tracks.
	Repository tasks.RepositoryEntry
	// Links are how each frontend reaches it, sorted by FrontendRepository, then
	// Frontend. A companion has at least one link, each with at least one
	// APIURLVariable.
	Links []CompanionLinkV1
}

// CompanionLinkV1 is how one static-site deployment of the Task (the frontend)
// reaches a companion in production.
type CompanionLinkV1 struct {
	FrontendRepository string
	// Frontend is the static-site deployment's name.
	Frontend string
	// APIURLVariables are the frontend's variables linked to the companion — a
	// person's entries and those whose production value is an http(s) URL on its
	// host — sorted by Name. Path is the production URL's path, as in
	// StackURLVariableV1; "" for a person's entry detection didn't match.
	APIURLVariables []StackURLVariableV1
	// AllowedOriginVariables are the companion's variables whose production value
	// lists the frontend's origin, sorted by Name.
	AllowedOriginVariables []StackOriginVariableV1
	// PersonVariables are the APIURLVariables names that come from a person's
	// entry, sorted.
	PersonVariables []string
}

// CompanionNotDeployedV1 is a frontend variable that links to a web service
// outside the Task whose preview won't be deployed, and why.
type CompanionNotDeployedV1 struct {
	FrontendRepository string
	// Frontend is the static-site deployment's name.
	Frontend string
	// Variable is the frontend's variable name.
	Variable string
	// API is the web service's name; "" for ReasonAmbiguous and ReasonAPIMissing.
	API string
	// Reason is one of the Reason constants.
	Reason string
}

// MaxCompanionBackends is the most companion backends a Task deploys.
const MaxCompanionBackends = 3

// The reasons a CompanionNotDeployedV1 gives.
const (
	ReasonAmbiguous     = "matches more than one web service"
	ReasonAPIMissing    = "the web service a person chose no longer exists or isn't built from a repository"
	ReasonNotDockerfile = "the web service isn't built from a Dockerfile"
	ReasonOtherRegion   = "the web service is in another region"
	ReasonInstallation  = "the web service's repository can't be accessed"
	ReasonLimit         = "a Task deploys at most 3 companion backends"
)
