package git_provider_enums

type Type uint

const (
	GitHub Type = iota + 1
	GitLab
)

var typeToString = map[Type]string{
	GitHub: "github",
	GitLab: "gitlab",
}

func (t Type) String() string {
	return typeToString[t]
}
