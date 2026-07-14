package deployment_enums

type DeletionState uint

const (
	DeletionNotStarted DeletionState = iota
	DeletionInProcess
	DeletionDone
	DeleteUndone //this is when delete is aborted ex. when PR is reopened
	// PreviewActive is a non-zero "alive" state for ephemeral preview environments.
	// The zero value DeletionNotStarted is erased by `omitempty` on the DeletionState
	// field, so a fresh active env stores no deletionState at all — which makes "active"
	// impossible to express in a partial index. Ephemeral previews are created with
	// PreviewActive instead so liveness is a stored, queryable ($eq) fact, enabling a
	// partial-unique index (one active preview env per task) without touching omitempty
	// or migrating existing data. Preview-only; permanent envs keep the zero convention.
	PreviewActive
	MaxDeletionState //always add new states before MaxDeletionState
)
