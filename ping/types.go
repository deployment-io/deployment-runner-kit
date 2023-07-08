package ping

import "github.com/deployment-io/deployment-runner-kit/types"

type ArgsV1 struct {
	types.AuthArgsV1
	Send      string
	FirstPing bool
	GoArch    string
}

type ReplyV1 struct {
	Send string
}
