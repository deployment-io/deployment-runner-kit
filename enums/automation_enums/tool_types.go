package automation_enums

import "fmt"

// ToolType is the type of tool that can be run on an entity
type ToolType uint

const (
	GetCPUMemoryUsage ToolType = iota + 1
	SendEmail
	GetApplicationLogs
	QueryCode

	ToolTypeMax //always add new types before ToolTypeMax
)

// Add the following maps and String method below
var toolTypeToStringMap = map[ToolType]string{
	GetCPUMemoryUsage:  "GetCPUMemoryUsage",
	SendEmail:          "Email",
	GetApplicationLogs: "GetApplicationLogs",
	QueryCode:          "QueryCode",
}

// String returns the string representation of the tool type
func (t ToolType) String() string {
	if str, ok := toolTypeToStringMap[t]; ok {
		return str
	}
	return "Unknown ToolType"
}

var toolTypeToEntities = map[ToolType][]Entity{
	GetCPUMemoryUsage: {
		WebService,
		PrivateService,
	},
	GetApplicationLogs: {
		WebService,
		PrivateService,
	},
	QueryCode: {
		StaticSite,
		WebService,
		PrivateService,
		Repository,
	},
}

// GetEntities returns the entities that can run this tool type
func (t ToolType) GetEntities() []Entity {
	return toolTypeToEntities[t]
}

func GetToolTypeFromNodeType(nodeType NodeType) (ToolType, error) {
	switch nodeType {
	case SendEmailTool:
		return SendEmail, nil
	case GetCPUMemoryUsageTool:
		return GetCPUMemoryUsage, nil
	case GetApplicationLogsTool:
		return GetApplicationLogs, nil
	case QueryCodeTool:
		return QueryCode, nil

	default:
		return 0, fmt.Errorf("node type is not of tool type")
	}
}
