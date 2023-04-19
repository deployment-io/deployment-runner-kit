package jobs

import "github.com/deployment-io/jobs-runner-kit/enums/parameters_enums"

type Runner interface {
	Run(parameters map[parameters_enums.Key]interface{}, logger Logger, jobContext *ContextV1) (map[parameters_enums.Key]interface{}, error)
}

type Command interface {
	Runner
}
