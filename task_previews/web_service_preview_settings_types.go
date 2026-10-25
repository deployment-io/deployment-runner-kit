package task_previews

import "github.com/deployment-io/deployment-runner-kit/types"

// WebServicePreviewSettingsArgsV1 asks deployment-server how the org deploys the
// repository at CloneURL as a web service, so the runner can build and deploy the
// task's backend preview the way production does — with the org's preview
// configuration (never the production deployment's variables), on deployment.io's
// network in the runner's region and cloud account.
type WebServicePreviewSettingsArgsV1 struct {
	types.AuthArgsV1
	TaskID string
	// CloneURL is the repository's clone URL as the Task's repository entry has it,
	// compared as StaticSiteBuildSettingsArgsV1.CloneURL is.
	CloneURL string

	// The calling runner's identity, as in EnsureTaskPreviewArgsV1: the preview runs
	// on the deployment.io network of this cloud account and region.
	CloudAccountID string
	Region         string // runner region, region_enums string e.g. "us-east-1"
}

// WebServicePreviewSettingsReplyV1 holds the org's previewable web-service
// deployments of the repository, the org's preview configuration and the network a
// backend preview deploys into. Its configuration values are secrets: the runner
// must not log them.
type WebServicePreviewSettingsReplyV1 struct {
	// Services are the org's undeleted web-service deployments of the repository
	// built from it (not deployed from an image) with a Dockerfile, sorted by
	// DeploymentName. Empty when there are none: no backend preview.
	Services []WebServicePreviewSettingsV1
	// RuntimeBuiltServices are the DeploymentNames, sorted, of the matching web
	// services production builds from a language runtime (Nixpacks). They get no
	// backend preview; the runner says why.
	RuntimeBuiltServices []string

	// The preview configuration, exactly as StaticSiteBuildSettingsReplyV1 has it:
	// filled whether or not any service matches.
	ConfigurationSet bool
	Variables        map[string]string
	Files            []PreviewFileV1

	// Network is filled only when NetworkReady.
	Network WebServicePreviewNetworkV1
	// NetworkReady is true when deployment.io's VPC and ECS cluster exist, complete,
	// in the runner's region and cloud account. False → NetworkMissingReason says
	// why, and there is no backend preview.
	NetworkReady         bool
	NetworkMissingReason string

	// NamespaceName is the org's shared Service Connect namespace for previews in
	// the region: "previews-" + the last 4 characters of the org ID, lower-case.
	NamespaceName string
}

// WebServicePreviewSettingsV1 is how one web-service deployment builds and runs.
type WebServicePreviewSettingsV1 struct {
	DeploymentName  string
	EnvironmentName string
	// Region is the deployment's region, a region_enums string e.g. "us-east-1".
	Region string
	// RootDirectory is the directory of the repository the build runs in ("" for
	// the repository root). StartCommand and PreDeployCommand are "" when unset.
	RootDirectory    string
	StartCommand     string
	PreDeployCommand string
	// DockerFilePath is "Dockerfile" when unset, as the production build defaults it.
	DockerFilePath  string
	DockerBuildArgs []string
	Port            int64
	// HealthCheckPath is "/" when unset, as production defaults it.
	HealthCheckPath string
	Cpu             string
	Memory          string

	// VariableNames are the NAMES, sorted, of the variables the production
	// deployment runs with: its environment's variables and those its dotenv
	// secret files define. Never a value. The runner uses them to choose between
	// several database address variables the code reads.
	VariableNames []string
	// PreviewDatabase is the person's choice of a preview database for this
	// service. Nil means automatic (detected from the Task's code).
	PreviewDatabase *PreviewDatabaseV1
}

// WebServicePreviewNetworkV1 is the network a production web-service deploy reads:
// deployment.io's VPC and ECS cluster of the region and cloud account.
type WebServicePreviewNetworkV1 struct {
	VpcID                   string
	VpcCidr                 string
	PublicSubnets           []string
	PrivateSubnets          []string
	EcsClusterArn           string
	EcsTaskExecutionRoleArn string
}
