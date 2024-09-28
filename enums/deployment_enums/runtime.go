package deployment_enums

import (
	"github.com/deployment-io/deployment-runner-kit/enums/commands_enums"
)

type Runtime uint

const (
	Docker Runtime = iota + 1
	Elixir
	Go
	Node
	Python3
	Ruby
	Rust
	MaxRuntime //always add new runtimes before MaxRuntime
)

func GetRuntimes() []Runtime {
	return []Runtime{Docker, Go, Node, Python3, Rust}
}

var runtimeToString = map[Runtime]string{
	Docker:  "Docker",
	Elixir:  "Elixir",
	Go:      "Go",
	Node:    "Node",
	Python3: "Python",
	Ruby:    "Ruby",
	Rust:    "Rust",
}

func (r Runtime) String() string {
	return runtimeToString[r]
}

func (r Runtime) GetBuildCommand() commands_enums.Type {
	switch r {
	case Docker:
		return commands_enums.BuildDockerImage
	default:
		return commands_enums.BuildNixPacksImage
	}
}
