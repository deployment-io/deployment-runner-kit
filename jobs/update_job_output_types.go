package jobs

import (
	"github.com/deployment-io/deployment-runner-kit/types"
)

type UpdateJobOutputDtoV1 struct {
	ID     string
	Output string
}

type UpdateJobOutputArgsV1 struct {
	types.AuthArgsV1
	Jobs []UpdateJobOutputDtoV1
}

type UpdateJobOutputReplyV1 struct {
	Done bool
}
