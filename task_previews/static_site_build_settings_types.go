package task_previews

import "github.com/deployment-io/deployment-runner-kit/types"

// StaticSiteBuildSettingsArgsV1 asks deployment-server how the org deploys the
// repository at CloneURL as a static site, so the runner can build the task's
// preview of that repository the way the deploy builds it — with the org's preview
// configuration, whose values the agent never sees.
type StaticSiteBuildSettingsArgsV1 struct {
	types.AuthArgsV1
	TaskID string
	// CloneURL is the repository's clone URL as the Task's repository entry has it.
	// deployment-server compares it case-insensitively, ignoring surrounding
	// whitespace, a trailing "/" and a trailing ".git".
	CloneURL string
}

// StaticSiteBuildSettingsReplyV1 holds the org's static-site deployments of the
// repository and the org's preview configuration, whether or not it deploys the
// repository (a runner builds the preview of a repository it doesn't deploy from
// settings it detects, with the same configuration).
type StaticSiteBuildSettingsReplyV1 struct {
	// Sites are the org's undeleted static-site deployments of the repository,
	// sorted by DeploymentName. Empty when deployment.io doesn't deploy it.
	Sites []StaticSiteBuildSettingsV1
	// ConfigurationSet is true whenever the org has a preview configuration, Sites
	// empty or not; Variables and Files are its decrypted values. False → no values.
	ConfigurationSet bool
	Variables        map[string]string
	Files            []PreviewFileV1
}

// StaticSiteBuildSettingsV1 is how one static-site deployment builds.
type StaticSiteBuildSettingsV1 struct {
	DeploymentName  string
	EnvironmentName string
	// RootDirectory is the directory of the repository the build runs in ("" for
	// the repository root).
	RootDirectory string
	// BuildCommand is "" when the deployment has none.
	BuildCommand string
	// PublishDirectory is the build's output directory, relative to RootDirectory.
	PublishDirectory string
	IsSpa            bool
}

// PreviewFileV1 is one secret file of the preview configuration, written at
// <repository>/<RootDirectory>/<Name> before the build.
type PreviewFileV1 struct {
	Name     string
	Contents string
}
