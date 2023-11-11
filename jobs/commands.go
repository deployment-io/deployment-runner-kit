package jobs

import "io"

type Runner interface {
	Run(parameters map[string]interface{}, logsWriter io.Writer) (map[string]interface{}, error)
}

type Command interface {
	Runner
}
