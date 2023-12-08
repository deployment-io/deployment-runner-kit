package logs

import (
	"github.com/deployment-io/deployment-runner-kit/types"
)

type AddJobLogDtoV1 struct {
	ID           string
	Message      string
	ErrorMessage string
	Ts           int64
}

type AddJobLogsArgsV1 struct {
	types.AuthArgsV1
	JobLogs []AddJobLogDtoV1
}

type AddJobLogsReplyV1 struct {
	Done bool
}
