package jobs

import "github.com/deployment-io/deployment-runner-kit/types"

type CompletingJobDtoV1 struct {
	Error string
	ID    string
}

type CompletingJobsArgsV1 struct {
	types.AuthArgsV1
	Jobs []CompletingJobDtoV1
}

type CompletingJobsReplyV1 struct {
	Done bool
}
