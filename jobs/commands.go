package jobs

type Runner interface {
	Run(parameters map[string]interface{}, logger Logger) (map[string]interface{}, error)
}

type Command interface {
	Runner
}
