package cpu_architecture_enums

type Type uint

const (
	ARM Type = iota + 1
	AMD
	MaxType //always add new types before MaxType
)

var typesToString = map[Type]string{
	ARM: "arm",
	AMD: "amd",
}

func (t Type) String() string {
	return typesToString[t]
}
