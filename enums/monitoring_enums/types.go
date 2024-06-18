package monitoring_enums

type Type uint

const (
	Cloudwatch Type = iota + 1
)

var typeToString = map[Type]string{
	Cloudwatch: "Cloudwatch",
}

func (t Type) String() string {
	return typeToString[t]
}
