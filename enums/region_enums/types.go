package region_enums

type Type uint

const (
	EuWest1 Type = iota + 1
	ApSouth1
	ApNorthEast1
	ApNorthEast2
	ApNorthEast3
	ApSouthEast1
	ApSouthEast2
	CaCentral1
	EuCentral1
	EuNorth1
	EuWest2
	EuWest3
	SaEast1
	UsEast1
	UsEast2
	UsWest1
	UsWest2
	AfSouth1
	ApEast1
	ApSouth2
	ApSouthEast3
	ApSouthEast4
	EuCentral2
	EuSouth1
	EuSouth2
	MeCentral1
	MeSouth1
	MaxType //always add new types before MaxType
)

var TypeToString = map[Type]string{
	EuWest1:      "eu-west-1",
	ApSouth1:     "ap-south-1",
	ApNorthEast1: "ap-northeast-1",
	ApNorthEast2: "ap-northeast-2",
	ApNorthEast3: "ap-northeast-3",
	ApSouthEast1: "ap-southeast-1",
	ApSouthEast2: "ap-southeast-2",
	CaCentral1:   "ca-central-1",
	EuCentral1:   "eu-central-1",
	EuNorth1:     "eu-north-1",
	EuWest2:      "eu-west-2",
	EuWest3:      "eu-west-3",
	SaEast1:      "sa-east-1",
	UsEast1:      "us-east-1",
	UsEast2:      "us-east-2",
	UsWest1:      "us-west-1",
	UsWest2:      "us-west-2",
	AfSouth1:     "af-south-1",
	ApEast1:      "ap-east-1",
	ApSouth2:     "ap-south-2",
	ApSouthEast3: "ap-southeast-3",
	ApSouthEast4: "ap-southeast-4",
	EuCentral2:   "eu-central-2",
	EuSouth1:     "eu-south-1",
	EuSouth2:     "eu-south-2",
	MeCentral1:   "me-central-1",
	MeSouth1:     "me-south-1",
}

func (a Type) String() string {
	return TypeToString[a]
}
