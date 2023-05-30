package cluster_enums

type Type uint

const (
	ECS     Type = iota + 1
	MaxType      //always add new types before MaxType
)

var typeToString = map[Type]string{
	ECS: "ECS",
}

func (t Type) String() string {
	return typeToString[t]
}
