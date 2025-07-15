package automation_enums

import "fmt"

type NodeType uint

const (
	CustomAgent NodeType = iota + 1
	SendEmailTool
	ChatTrigger
	GetCPUMemoryUsageTool
	GetApplicationLogsTool
	QueryCodeTool

	NodeTypeMax //always add new types before NodeTypeMax
)

func GetNodeTypeFromKey(key string) (NodeType, error) {
	switch key {
	case "customAgent":
		return CustomAgent, nil
	case "chatTrigger":
		return ChatTrigger, nil
	case "sendEmailTool":
		return SendEmailTool, nil
	case "getCpuMemoryUsageTool":
		return GetCPUMemoryUsageTool, nil
	case "getApplicationLogsTool":
		return GetApplicationLogsTool, nil
	case "queryCodeTool":
		return QueryCodeTool, nil
	}
	return 0, fmt.Errorf("unknown node type %s", key)
}

var typeToString = map[NodeType]string{
	CustomAgent:            "Custom Agent",
	SendEmailTool:          "Email Tool",
	ChatTrigger:            "Chat Trigger",
	GetCPUMemoryUsageTool:  "CPU Memory Usage Tool",
	GetApplicationLogsTool: "Application Logs Tool",
	QueryCodeTool:          "Query Code Tool",
}

func (n NodeType) String() string {
	return typeToString[n]
}

func (n NodeType) IsToolType() bool {
	switch n {
	case SendEmailTool:
		return true
	case GetCPUMemoryUsageTool:
		return true
	case GetApplicationLogsTool:
		return true
	case QueryCodeTool:
		return true
	default:
		return false
	}
}

func (n NodeType) IsTriggerType() bool {
	switch n {
	case ChatTrigger:
		return true
	default:
		return false
	}
}

func (n NodeType) IsAgentType() bool {
	switch n {
	case CustomAgent:
		return true
	default:
		return false
	}
}
