package agents

import (
	"github.com/deployment-io/deployment-runner-kit/enums/agent_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/llm_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/message_enums"
)

type NodeDtoV1 struct {
	ID            string
	Name          string
	Children      []string
	IsStart       bool
	NodeType      agent_enums.NodeType
	Goal          string
	Backstory     string
	ToolType      agent_enums.ToolType
	TriggerType   agent_enums.TriggerType
	LlmModelType  llm_enums.ModelType
	LlmApiVersion llm_enums.ApiVersion
}

type MessageDataDtoV1 struct {
	Content string
	Type    message_enums.Type
}

type AgentDataDtoV1 struct {
	ID                 string
	Name               string
	Goal               string
	StartNodeID        string
	NodesMap           map[string]NodeDtoV1
	ChatHistory        []MessageDataDtoV1
	LlmModelType       llm_enums.ModelType
	LlmApiVersion      llm_enums.ApiVersion
	OpenAIAPIKey       string
	OpenAIBaseUrl      string
	ExtraHelpAvailable bool
}
