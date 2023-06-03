package jobs

import (
	"encoding/gob"
	"fmt"
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
	int | int32 | string | uint | map[string]string | primitive.A | region_enums.Type | []string
}

func RegisterGobDataTypes() {
	gob.Register(map[string]string{})
	gob.Register(primitive.A{})
	gob.Register(region_enums.MaxType)
	gob.Register([]string{})
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

func SetParameterValue[T ParameterType](parameters map[parameters_enums.Key]interface{}, k parameters_enums.Key, v T) {
	parameters[k] = v
}
