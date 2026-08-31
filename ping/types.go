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
	// Host specs of the machine the runner is running on. Reported so the
	// control plane can show what a runner actually has and validate any
	// resource setting against it — the Runner model previously carried
	// no notion of the underlying instance at all, so "why did my Task run
	// out of memory" was unanswerable from the dashboard.
	//
	// Sent by the RUNNER, not the controller: the controller is a Fargate
	// task and knows nothing about the EC2 host. The runner scales to zero
	// between jobs, so these are written on its FirstPing and then persist
	// — the values describe the last boot, which is correct until the
	// instance is resized.
	//
	// All three are best-effort and may be zero/empty (a local runner, or
	// IMDS unreachable). Consumers must treat them as optional; the model
	// fields are omitempty so a ping without them never overwrites a
	// previously reported value.
	HostMemoryBytes int64
	HostCPUCores    int64
	// InstanceType is the EC2 instance type from IMDS, e.g. "m6a.large".
	// Empty on any non-EC2 runner.
	InstanceType string
}

type ReplyV2 struct {
	Send                         string
	RunnerUpgradeDockerImage     string
	ControllerUpgradeDockerImage string
	UpgradeFromTs                int64
	UpgradeToTs                  int64
}
