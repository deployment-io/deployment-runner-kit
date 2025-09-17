package jobs

import (
	"encoding/gob"
	"fmt"

	"github.com/deployment-io/deployment-runner-kit/agents"
	"github.com/deployment-io/deployment-runner-kit/automations"
	"github.com/deployment-io/deployment-runner-kit/enums/parameters_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/region_enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContextV1 struct {
	OrganizationID string
	EnvironmentID  string
	DeploymentID   string
	BuildID        string
}

type ParameterType interface {
	int | int64 | string | map[string]string | primitive.A | primitive.D | []string | map[string]interface{} |
		[][]string | bool | automations.AutomationDataDtoV1 | agents.AgentDataDtoV1
}

func RegisterGobDataTypes() {
	gob.Register(map[string]string{})
	gob.Register(map[string]interface{}{})
	gob.Register(primitive.A{})
	gob.Register(region_enums.MaxType)
	gob.Register([]string{})
	gob.Register([][]string{{}})
	gob.Register(false)
	gob.Register(automations.AutomationDataDtoV1{})
	gob.Register(agents.AgentDataDtoV1{})
}

func GetParameterValue[T ParameterType](parameters map[string]interface{}, k parameters_enums.Key) (T, error) {
	key, err := k.Key()
	if err != nil {
		var zeroValue T
		return zeroValue, err
	}
	i, ok := parameters[key]
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

func SetParameterValue[T ParameterType](parameters map[string]interface{}, k parameters_enums.Key, v T) error {
	key, err := k.Key()
	if err != nil {
		return err
	}
	parameters[key] = v
	return nil
}
