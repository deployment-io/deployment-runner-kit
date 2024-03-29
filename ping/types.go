package ping

import "github.com/deployment-io/deployment-runner-kit/types"

type ArgsV1 struct {
	types.AuthArgsV1
	Send      string
	FirstPing bool
	GoArch    string
}

type ReplyV1 struct {
	Send               string
	DockerUpgradeImage string
	UpgradeFromTs      int64
	UpgradeToTs        int64
}

type ArgsV2 struct {
	types.AuthArgsV2
	Send      string
	FirstPing bool
}

type ReplyV2 struct {
	Send                         string
	RunnerUpgradeDockerImage     string
	ControllerUpgradeDockerImage string
	UpgradeFromTs                int64
	UpgradeToTs                  int64
}
