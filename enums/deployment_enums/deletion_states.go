package deployment_enums

type DeletionState uint

const (
	DeletionNotStarted DeletionState = iota
	DeletionInProcess
	DeletionDone
	MaxDeletionState //always add new states before MaxDeletionState
)
