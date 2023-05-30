package region_enums

type Type uint

const (
	EuWest1 Type = iota + 1
	ApSouth1
	MaxType //always add new types before MaxType
)

var typeToString = map[Type]string{
	EuWest1:  "eu-west-1",
	ApSouth1: "ap-south-1",
}

func (a Type) String() string {
	return typeToString[a]
}
