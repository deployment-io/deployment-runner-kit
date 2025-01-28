package automation_enums

import "fmt"

type TriggerType uint

const (
	Chat TriggerType = iota + 1

	MaxTriggerType //always add new types before MaxTriggerType
)

var triggerTypeToString = map[TriggerType]string{
	Chat: "Chat",
}

func (t TriggerType) String() string {
	return triggerTypeToString[t]
}

func GetTriggerTypeFromNodeType(nodeType NodeType) (TriggerType, error) {
	switch nodeType {
	case ChatTrigger:
		return Chat, nil
	default:
		return 0, fmt.Errorf("node type is not of trigger type")
	}
}
