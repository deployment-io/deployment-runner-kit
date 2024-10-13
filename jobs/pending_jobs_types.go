package jobs

import (
	"github.com/deployment-io/deployment-runner-kit/enums/commands_enums"
	"github.com/deployment-io/deployment-runner-kit/types"
)

// PendingJobsArgsV1 represents v1 args to get pending jobs
type PendingJobsArgsV1 struct {
	types.AuthArgsV1
}

// PendingJobDtoV1 represents v1 of a pending deployment job from the server
type PendingJobDtoV1 struct {
	JobID        string
	CommandEnums []commands_enums.Type
	Parameters   map[string]interface{}
}

// PendingJobsDtoV1 represents v1 of a deployment jobs from the server
type PendingJobsDtoV1 struct {
	Jobs []PendingJobDtoV1
}

type PendingJobsArgsV2 struct {
	types.AuthArgsV2
}

type PendingJobsForSaasArgsV1 struct {
	types.AuthArgsV2
}

type PendingJobForSaasDtoV1 struct {
	JobID          string
	OrganizationID string
	CommandEnums   []commands_enums.Type
	Parameters     map[string]interface{}
}

type PendingJobsForSaasDtoV1 struct {
	Jobs []PendingJobForSaasDtoV1
}
