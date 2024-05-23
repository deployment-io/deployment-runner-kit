package registry_enums

type Type uint

const (
	GitHub Type = iota + 1
	GitLab
	Docker
	MaxType //always add new types before MaxType
)

var typeToString = map[Type]string{
	GitHub: "github",
	GitLab: "gitlab",
	Docker: "docker",
}

func (t Type) String() string {
	return typeToString[t]
}
