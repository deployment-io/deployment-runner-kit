package os_enums

type Type uint

const (
	LINUX Type = iota + 1
	WINDOWS
	MaxType //always add new types before MaxType
)

var typesToString = map[Type]string{
	LINUX:   "linux",
	WINDOWS: "windows",
}

func (t Type) String() string {
	return typesToString[t]
}
