package types

import (
	"github.com/deployment-io/jobs-runner-kit/enums/commands_enums"
	"github.com/deployment-io/jobs-runner-kit/enums/parameters_enums"
)

// Args represents server args
type Args struct {
	A string
	P string
}

// JobDtoV1 represents v1 of a deployment job from the server
type JobDtoV1 struct {
	CommandEnums []commands_enums.Type
	Parameters   map[parameters_enums.Key]interface{}
}

// JobsDtoV1 represents v1 of a deployment jobs from the server
type JobsDtoV1 struct {
	Jobs []JobDtoV1
}
