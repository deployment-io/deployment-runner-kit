package jobs

import "github.com/deployment-io/deployment-runner-kit/types"

type UpsertJobHeartbeatArgsV1 struct {
	types.AuthArgsV1
	JobID string
}

type UpsertJobHeartbeatReplyV1 struct {
	Stopping bool
}
