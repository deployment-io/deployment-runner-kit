package builds

import (
	"github.com/deployment-io/deployment-runner-kit/enums/build_enums"
	"github.com/deployment-io/deployment-runner-kit/types"
)

type UpdateBuildDtoV1 struct {
	ID            string
	BuildTs       int64
	CommitHash    string
	CommitMessage string
	Status        build_enums.Status
}

type UpdateBuildsArgsV1 struct {
	types.AuthArgsV1
	Builds []UpdateBuildDtoV1
}

type UpdateBuildsReplyV1 struct {
	Done bool
}
