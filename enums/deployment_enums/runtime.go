package deployment_enums

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
	return []Runtime{Docker, Elixir, Go, Node, Python3, Ruby, Rust}
}

var runtimeToString = map[Runtime]string{
	Docker:  "Docker",
	Elixir:  "Elixir",
	Go:      "Go",
	Node:    "Node",
	Python3: "Python3",
	Ruby:    "Ruby",
	Rust:    "Rust",
}

func (r Runtime) String() string {
	return runtimeToString[r]
}
