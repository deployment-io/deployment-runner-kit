package jobs

import "github.com/deployment-io/deployment-runner-kit/types"

type JobsCountArgsV1 struct {
	types.AuthArgsV1
}

type JobsCountDtoV1 struct {
	Count                  int
	RunnerOnTimeoutSeconds int
}

type JobsCountArgsV2 struct {
	types.AuthArgsV2
}

type SaasJobsCountArgsV1 struct {
	types.AuthArgsV2
}

type SaasJobsCountDtoV1 struct {
	Count                  int
	RunnerOnTimeoutSeconds int
}
