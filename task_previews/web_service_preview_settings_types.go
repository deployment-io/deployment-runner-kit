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
	// Branch is the branch the deployment deploys (its repository branch), ""
	// when unknown.
	Branch string
	// ProductionEnvironment is true when the deployment's environment is a
	// permanent, not deleted (or being deleted) environment classified
	// Production. With Branch, it lets the runner choose between several
	// services of the same repository and root directory (Staging and
	// Production): the one deploying the Task's base branch, in Production. A
	// runner built before these fields ignores them (gob skips fields its
	// struct lacks); an older deployment-server leaves them "" and false.
	ProductionEnvironment bool
	// Region is the deployment's region, a region_enums string e.g. "us-east-1".
	Region string
	// RootDirectory is the directory of the repository the build runs in ("" for
	// the repository root). StartCommand and PreDeployCommand are "" when unset.
	RootDirectory    string
	StartCommand     string
	PreDeployCommand string
	// DockerFilePath is "Dockerfile" when unset, as the production build defaults it.
	DockerFilePath string
	// DockerBuildArgs are the NAMES of the production build args ("NAME", the
	// part before the first "="), never their values: the runner reads only the
	// names, so production build-arg values don't leave deployment-server.
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
	// PreviewMigrateCommand is the command a person set for this service's
	// previews to run against the preview's own database before the API starts
	// (its migrations). "" means a preview runs no migrations: the production
	// PreDeployCommand never runs in a preview, since it is production
	// configuration and can name a production database or carry credentials.
	PreviewMigrateCommand string
	// PreviewSeedCommand is the command a person set for this service's
	// previews to run against the preview's own database after
	// PreviewMigrateCommand, on every preview deploy: it creates a user for each
	// of PreviewRoles. "" means the preview runs no seed command. A runner built
	// before this field ignores it (gob skips fields its struct lacks).
	PreviewSeedCommand string
	// PreviewRoles are the test roles a person set for this service's previews:
	// the seed command creates a user for each, with an email and password
	// generated for the preview. Nil when there are none. A runner built before
	// this field ignores it (gob skips fields its struct lacks).
	PreviewRoles []PreviewRoleV1
	// PreviewSignIn is the sign-in request a person set for this service's
	// previews: how QA signs in as one of PreviewRoles with the email and
	// password generated for it. Nil when the service has none. A runner built
	// before this field ignores it (gob skips fields its struct lacks).
	PreviewSignIn *PreviewSignInV1
	// PreviewBrowserSignInPath is the path of the page where a person signs in
	// with an email and password, e.g. /signin, on the site QA's browser checks.
	// "" when unset: QA's browser doesn't sign in as this service's roles. A
	// runner built before this field ignores it (gob skips fields its struct
	// lacks).
	PreviewBrowserSignInPath string
	// BuildSecrets are this service's preview build secrets, decrypted: name →
	// value, nil when it has none. Never production's build secrets. They are for
	// this service's backend preview build only, and stay runner-side: never in an
	// agent's tool result, a log or the Step log. A runner built before this field
	// ignores it (gob skips fields its struct lacks).
	BuildSecrets map[string]string
}

// PreviewRoleV1 is a test role of a web service's previews. Name is lower case
// (letters, digits, - and _); Description says what the user can do, "" when
// unset.
type PreviewRoleV1 struct {
	Name        string
	Description string
}

// PreviewSignInV1 is how a web service's preview signs in a test user: a POST
// to Path with the user's email and password, answered with a session.
type PreviewSignInV1 struct {
	// Path is the path the app signs in at: starts with /, no ? or #.
	Path string
	// BodyFormat is "json" or "form". With "json", a "." in a field name nests
	// objects (user.email → {"user": {"email": …}}); with "form" the names are
	// used as they are (user[email] stays user[email]).
	BodyFormat string
	// EmailField is the name of the body field holding the email.
	EmailField string
	// PasswordField is the name of the body field holding the password.
	PasswordField string
	// Session is how the app answers with the session: "cookie" (its cookies)
	// or "bearer" (a token in a field of its JSON body).
	Session string
	// TokenField is the dot path into the JSON body that holds the token, e.g.
	// token or data.access_token. Set only with "bearer"; "" with "cookie".
	TokenField string
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
