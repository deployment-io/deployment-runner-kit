package automations

import (
	"github.com/deployment-io/deployment-runner-kit/enums/automation_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/llm_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/message_enums"
)

type NodeDtoV1 struct {
	ID           string
	Children     []string
	IsStart      bool
	NodeType     automation_enums.NodeType
	Goal         string
	Backstory    string
	ToolType     automation_enums.ToolType
	TriggerType  automation_enums.TriggerType
	LlmModelType llm_enums.ModelType
}

type MessageDataDtoV1 struct {
	Content string
	Type    message_enums.Type
}

type AutomationDataDtoV1 struct {
	ID                 string
	Name               string
	Goal               string
	StartNodeID        string
	Entity             automation_enums.Entity
	NodesMap           map[string]NodeDtoV1
	ChatHistory        []MessageDataDtoV1
	LlmModelType       llm_enums.ModelType
	OpenAIAPIKey       string
	OpenAIBaseUrl      string
	ExtraHelpAvailable bool
}
