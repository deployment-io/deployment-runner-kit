package build_enums

type Status uint

const (
	Pending Status = 1
	Success Status = 2
	Error   Status = 3
)

var typeToString = map[Status]string{
	Pending: "Pending",
	Success: "Success",
	Error:   "Error",
}

func (t Status) String() string {
	return typeToString[t]
}
