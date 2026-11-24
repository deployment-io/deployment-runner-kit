package tasks

import "github.com/deployment-io/deployment-runner-kit/enums/region_enums"

// ProductionEvidenceV1 is which production logs the runner reads for a fix
// Task's Step: the services the original Task (SourceTaskID) deployed, each
// with the window of its deploy. kit stamps it AT JOB CREATION of a fix
// Task's Step Job as JSON in the ProductionEvidence parameter; the runner only
// reads it.
//
// Skipped are the original's deployed repositories that have no logs the
// runner can read, each with a one-line reason the agent is shown.
type ProductionEvidenceV1 struct {
	SourceTaskID string                        `json:"sourceTaskId"`
	Services     []ProductionEvidenceServiceV1 `json:"services"`
	Skipped      []ProductionEvidenceSkippedV1 `json:"skipped"`
}

// ProductionEvidenceServiceV1 is one deployed service whose runtime logs the
// runner reads: the streams of build BuildID in the log group of deployment
// DeploymentID, in Region, between From and To (unix seconds). Name is the
// repository's name.
type ProductionEvidenceServiceV1 struct {
	Name         string            `json:"name"`
	DeploymentID string            `json:"deploymentId"`
	BuildID      string            `json:"buildId"`
	Region       region_enums.Type `json:"region"`
	From         int64             `json:"from"`
	To           int64             `json:"to"`
}

// ProductionEvidenceSkippedV1 is one deployed repository whose logs aren't
// read. Name is the repository's name, Reason one line that holds no secret.
type ProductionEvidenceSkippedV1 struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}
