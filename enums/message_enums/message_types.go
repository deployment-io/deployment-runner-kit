package message_enums

import "fmt"

type Type uint

const (
	User Type = iota + 1
	Assistant
	Output
	MaxType //always add new types before MaxType
)

var typeToString = map[Type]string{
	User:      "user",
	Assistant: "assistant",
	Output:    "output",
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
	default:
		return 0, fmt.Errorf("invalid message type")
	}
}
