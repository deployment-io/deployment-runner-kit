package git_provider_enums

type Type uint

const (
	GitHub Type = iota + 1
	GitLab
	BitBucket
)

var typeToString = map[Type]string{
	GitHub:    "GitHub",
	GitLab:    "GitLab",
	BitBucket: "BitBucket",
}

func (t Type) String() string {
	return typeToString[t]
}
