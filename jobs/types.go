package jobs

import (
	"fmt"
	"github.com/deployment-io/deployment-runner-kit/enums/commands_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/parameters_enums"
)

// ArgsV1 represents args for jobs methods
type ArgsV1 struct {
	A string
	P string
}

type ContextV1 struct {
	OrganizationID string
	EnvironmentID  string
	DeploymentID   string
	BuildID        string
}

// JobDtoV1 represents v1 of a deployment job from the server
type JobDtoV1 struct {
	CommandEnums []commands_enums.Type
	Parameters   map[parameters_enums.Key]interface{}
}

// DtoV1 represents v1 of a deployment jobs from the server
type DtoV1 struct {
	Jobs []JobDtoV1
}

type ParameterType interface {
	int | string | uint
}

func GetParameterValue[T ParameterType](parameters map[parameters_enums.Key]interface{}, k parameters_enums.Key) (T, error) {
	i, ok := parameters[k]
	if !ok {
		var zeroValue T
		return zeroValue, fmt.Errorf("%s is missing", k.String())
	}
	v, ok := i.(T)
	if !ok {
		return v, fmt.Errorf("%s is not of valid type", k.String())
	}
	return v, nil
}
