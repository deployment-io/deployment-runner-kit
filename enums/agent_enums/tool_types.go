package agent_enums

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
	GetCPUMemoryUsage:  "Get CPU and Memory Usage",
	SendEmail:          "Send an email to a specified email address",
	GetApplicationLogs: "Get Application Logs",
	QueryCode:          "Query Code",
}

// String returns the string representation of the tool type
func (t ToolType) String() string {
	if str, ok := toolTypeToStringMap[t]; ok {
		return str
	}
	return "Unknown ToolType"
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
