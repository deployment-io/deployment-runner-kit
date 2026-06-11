package message_enums

import "fmt"

type Type uint

const (
	User Type = iota + 1
	Assistant
	Output
	// TurnEnd is a control record, not chat content: the session agent
	// finished its turn and is waiting for the next user message. Carries no
	// content; the UI uses it to gate the composer and group a turn's
	// messages.
	TurnEnd
	MaxType //always add new types before MaxType
)

var typeToString = map[Type]string{
	User:      "user",
	Assistant: "assistant",
	Output:    "output",
	TurnEnd:   "turnEnd",
}

func (t Type) String() string {
	return typeToString[t]
}

func StringToType(typeStr string) (Type, error) {
	switch typeStr {
	case "user":
		return User, nil
	case "assistant":
		return Assistant, nil
	case "output":
		return Output, nil
	case "turnEnd":
		return TurnEnd, nil
	default:
		return 0, fmt.Errorf("invalid message type")
	}
}
