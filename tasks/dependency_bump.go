package tasks

import (
	"regexp"

	"github.com/deployment-io/deployment-runner-kit/types"
)

// DependencyProvidedV1 is one module or package a repository of a Task
// provides (Go: its go.mod module path).
//
// The dependency types in this file are ECOSYSTEM-NEUTRAL: only the runner's
// per-ecosystem adapters interpret Ecosystem, Name or the versions, and every
// other service (kit, deployment-server, app-server, the dashboard) stores and
// shows them as they are, without a list of ecosystems or a switch on one.
//
// Repository fields here and below are indexes into the Task's Repositories.
type DependencyProvidedV1 struct {
	Repository int
	Ecosystem  string
	Name       string
}

// DependencyEcosystemGo is the Go adapter's ecosystem name. An ecosystem name
// is chosen by its runner adapter; new ones need no change here or anywhere
// but the runner.
const DependencyEcosystemGo = "go"

var dependencyEcosystemPattern = regexp.MustCompile(`^[a-z][a-z0-9-]{0,31}$`)

// ValidDependencyEcosystem reports whether name is a well-formed ecosystem
// name: a lower-case letter, then up to 31 lower-case letters, digits or
// hyphens. It's the only check any service makes on an ecosystem name.
func ValidDependencyEcosystem(name string) bool {
	return dependencyEcosystemPattern.MatchString(name)
}

// Release kinds, in DependencyBumpTargetV1.Release.
const (
	// DependencyReleaseCommit: the merged commit is the release, as for Go.
	DependencyReleaseCommit = "commit"
	// DependencyReleaseRegistry: wait for a published version. Reserved:
	// nothing uses it yet.
	DependencyReleaseRegistry = "registry"
)

// Dependency bump states. kit stores them on the Task repository's bumps; a
// runner reports only DependencyBumpResultStates.
const (
	// DependencyBumpStateBumping is a claimed bump whose Job hasn't reported.
	DependencyBumpStateBumping = "bumping"
	// DependencyBumpStateBumped: the runner changed the pin and pushed.
	DependencyBumpStateBumped = "bumped"
	// DependencyBumpStateAlreadyCurrent: the repository already pins it.
	DependencyBumpStateAlreadyCurrent = "already_current"
	// DependencyBumpStateNotRequired: the repository no longer requires it.
	DependencyBumpStateNotRequired = "not_required"
	// DependencyBumpStateFailed: the bump couldn't be made; Error says why.
	DependencyBumpStateFailed = "failed"
	// DependencyBumpStateWaitingRelease is reserved for registry releases;
	// nothing sends it yet.
	DependencyBumpStateWaitingRelease = "waiting_release"
)

// DependencyBumpResultStates are the states a runner reports in a
// DependencyBumpResultDtoV1.
var DependencyBumpResultStates = []string{
	DependencyBumpStateBumped,
	DependencyBumpStateAlreadyCurrent,
	DependencyBumpStateNotRequired,
	DependencyBumpStateFailed,
}

// Limits, in bytes or entries.
const (
	DependencyNameMaxBytes        = 256
	DependencyDescriptionMaxBytes = 200
	DependencyVersionMaxBytes     = 128
	DependencyErrorMaxBytes       = 4096
	// MaxTaskDependencies caps Provided and Required, each.
	MaxTaskDependencies = 64
	// MaxDependencyBumpsPerJob caps the targets of one dependency bump Job.
	MaxDependencyBumpsPerJob = 16
	// MaxStartUpBumpSources caps the sources of one Step Job's
	// StartUpBumpSources.
	MaxStartUpBumpSources = 16
)

// DependencyRequiredV1 is one dependency between two repositories of a Task:
// Repository requires Name, which ProvidedBy (another repository of the same
// Task) provides, and pins Version now. Description is how messages name it
// ("the Go module github.com/deployment-io/kit"). Version is opaque.
type DependencyRequiredV1 struct {
	Repository  int
	Ecosystem   string
	Name        string
	Description string
	ProvidedBy  int
	Version     string
}

// TaskDependenciesDtoV1 is the dependency graph between a Task's
// repositories, which the runner reads after a Step run and sends to
// deployment-server (Jobs.RecordTaskDependenciesV1), which validates it,
// bounds it and stores it on the Task. JobID is the Step Job (agent run) the
// graph was read from. Empty Provided and Required are a valid graph: this
// run found no dependencies between the Task's repositories.
type TaskDependenciesDtoV1 struct {
	TaskID   string
	JobID    string
	Provided []DependencyProvidedV1
	Required []DependencyRequiredV1
}

type RecordTaskDependenciesArgsV1 struct {
	types.AuthArgsV1
	Dependencies TaskDependenciesDtoV1
}

// RecordTaskDependenciesReplyV1: Done is true once the server handled the
// report, including when it dropped it (an older run, a deleted Task).
type RecordTaskDependenciesReplyV1 struct {
	Done bool
}

// DependencyBumpRequestV1 is the parameter of a dependency bump Job: bump the
// Task repository at Repository to Dependencies. app-server stamps it AT JOB
// CREATION as JSON in the DependencyBumpRequest parameter; the runner only
// reads it.
type DependencyBumpRequestV1 struct {
	Repository   int                      `json:"repository"`
	Dependencies []DependencyBumpTargetV1 `json:"dependencies"`
}

// DependencyBumpTargetV1 is one library to bump to: Name of Ecosystem, to the
// release of the merged Commit of the Task repository at Source. Release is
// DependencyReleaseCommit or DependencyReleaseRegistry. app-server copies
// Ecosystem, Name and Description from the recorded graph.
type DependencyBumpTargetV1 struct {
	Ecosystem   string `json:"ecosystem"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      int    `json:"source"`
	Commit      string `json:"commit"`
	Release     string `json:"release"`
}

// DependencyBumpResultDtoV1 is what became of one target of a dependency bump
// Job. State is one of DependencyBumpResultStates. FromVersion and ToVersion
// are the pins before and after (ToVersion is set when bumped); Error is set
// when failed and holds no secret value.
type DependencyBumpResultDtoV1 struct {
	Ecosystem   string
	Name        string
	State       string
	FromVersion string
	ToVersion   string
	Error       string
}

// UpdateTaskDependencyBumpDtoV1 carries the end of a dependency bump Job to
// deployment-server (Jobs.UpdateTaskDependencyBumpV1), which validates it,
// bounds it and stores it on the Task repository's bumps.
//
// CommitSHA is the bump commit the runner pushed, set when any result is
// bumped. PRURL and PRNumber are set only when the runner opened the
// repository's pull request itself, because the repository had none.
type UpdateTaskDependencyBumpDtoV1 struct {
	TaskID     string
	JobID      string
	Repository int
	CommitSHA  string
	PRURL      string
	PRNumber   int
	Results    []DependencyBumpResultDtoV1
}

type UpdateTaskDependencyBumpArgsV1 struct {
	types.AuthArgsV1
	Update UpdateTaskDependencyBumpDtoV1
}

// UpdateTaskDependencyBumpReplyV1: Done is true once the server handled the
// report, including when it matched no bump (a late or repeated report).
type UpdateTaskDependencyBumpReplyV1 struct {
	Done bool
}

// StartUpBumpSourceV1 is one merged repository of an upstream Task (one
// the Task's Step Job waits for) and what it provides: the runner bumps
// the Task's repositories that require one of Provides to Commit before
// the agent runs. Stamped AT JOB CREATION; the runner only reads it.
type StartUpBumpSourceV1 struct {
	TaskID     string                  `json:"taskId"`
	TaskTitle  string                  `json:"taskTitle"`
	Repository RepositoryEntry         `json:"repository"`
	Commit     string                  `json:"commit"`
	Provides   []StartUpBumpProvidedV1 `json:"provides"`
}

// StartUpBumpProvidedV1 is a module or package the source provides.
type StartUpBumpProvidedV1 struct {
	Ecosystem string `json:"ecosystem"`
	Name      string `json:"name"`
}
