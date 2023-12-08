package deployment_enums

type DeletionState uint

const (
	DeletionNotStarted DeletionState = iota
	DeletionInProcess
	DeletionDone
	DeleteUndone     //this is when delete is aborted ex. when PR is reopened
	MaxDeletionState //always add new states before MaxDeletionState
)
