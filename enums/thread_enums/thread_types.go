package thread_enums

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Type uint

const (
	ServiceType Type = iota + 1
	AgentType
	MaxType //always add new types before MaxType
)

func GetType(pageContext string) (primitive.ObjectID, Type, error) {
	idHex, isThreadForService := strings.CutPrefix(pageContext, "service/")
	if isThreadForService {
		deploymentID, err := primitive.ObjectIDFromHex(idHex)
		if err != nil {
			return primitive.NilObjectID, MaxType, err
		}
		return deploymentID, ServiceType, nil
	}
	idHex, isThreadForAgent := strings.CutPrefix(pageContext, "agent/")
	if isThreadForAgent {
		agentID, err := primitive.ObjectIDFromHex(idHex)
		if err != nil {
			return primitive.NilObjectID, MaxType, err
		}
		return agentID, AgentType, nil
	}
	return primitive.NilObjectID, MaxType, nil
}

func (t Type) IsService() bool {
	if t == ServiceType {
		return true
	}
	return false
}

func (t Type) IsAgent() bool {
	if t == AgentType {
		return true
	}
	return false
}
