package automations

import "github.com/deployment-io/deployment-runner-kit/types"

type UpdateResponseDtoV1 struct {
	JobID    string
	Output   string
	NeedHelp bool
}

type UpdateResponseArgsV1 struct {
	types.AuthArgsV1
	Responses []UpdateResponseDtoV1
}

type UpdateResponseReplyV1 struct {
	Done bool
}
