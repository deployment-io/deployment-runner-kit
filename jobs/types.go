package jobs

import (
	"encoding/gob"
	"fmt"
	"github.com/deployment-io/deployment-runner-kit/enums/parameters_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/region_enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type ContextV1 struct {
	OrganizationID string
	EnvironmentID  string
	DeploymentID   string
	BuildID        string
}

type ParameterType interface {
	int | string | uint | map[string]string | primitive.A | region_enums.Type
}

func RegisterGobDataTypes() {
	gob.Register(map[string]string{})
	gob.Register(primitive.A{})
	gob.Register(region_enums.MaxType)
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

func EncodeEnvironmentVariablesToString(envVariables map[string]string) string {
	envString := ""
	for k, v := range envVariables {
		entry := fmt.Sprintf("%s=%s", k, v)
		envString += entry + "\n"
	}
	return envString
}

func DecodeEnvironmentVariablesToSlice(envVariables string) []string {
	envString := strings.TrimSpace(envVariables)
	envVariablesInSlice := strings.Split(envString, "\n")
	var envVariablesSlice []string
	for _, evIn := range envVariablesInSlice {
		ev := strings.TrimSpace(evIn)
		envVariablesSlice = append(envVariablesSlice, ev)
	}
	return envVariablesSlice
}
