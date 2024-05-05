package runner_enums

type Mode uint

const (
	LOCAL Mode = iota + 1
	AwsEcs
	MaxMode //always add new modes before MaxMode
)

var modesToString = map[Mode]string{
	LOCAL:  "local",
	AwsEcs: "aws",
}

func (m Mode) String() string {
	return modesToString[m]
}
