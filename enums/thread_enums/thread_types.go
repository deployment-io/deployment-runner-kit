package thread_enums

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type Type uint

const (
	ServiceType Type = iota + 1
	MaxType          //always add new types before MaxType
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
	return primitive.NilObjectID, MaxType, nil
}

func (t Type) IsService() bool {
	if t == ServiceType {
		return true
	}
	return false
}
