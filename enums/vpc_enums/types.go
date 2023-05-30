package vpc_enums

type Type uint

const (
	AwsVpc  Type = iota + 1
	MaxType      //always add new types before MaxType
)

var typeToString = map[Type]string{
	AwsVpc: "AWS VPC",
}

func (t Type) String() string {
	return typeToString[t]
}
