package logs

import (
	"github.com/deployment-io/deployment-runner-kit/types"
)

type AddBuildLogDtoV1 struct {
	ID           string
	Message      string
	ErrorMessage string
	Ts           int64
}

type AddBuildLogsArgsV1 struct {
	types.AuthArgsV1
	BuildLogs []AddBuildLogDtoV1
}

type AddBuildLogsReplyV1 struct {
	Done bool
}
