package deployment_enums

type Runtime uint

const (
	Docker     Runtime = iota + 1
	MaxRuntime         //always add new runtimes before MaxRuntime
)

var runtimeToString = map[Runtime]string{
	Docker: "Docker",
}

func (r Runtime) String() string {
	return runtimeToString[r]
}
