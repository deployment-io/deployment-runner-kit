package region_enums

type Type uint

const (
	EuWest1 Type = iota + 1
	ApSouth1
	MaxType //always add new types before MaxType
)

//ap-northeast-1:
//AMI: ami-0a8be33c358a98829
//ap-northeast-2:
//AMI: ami-04ea56ea021b86658
//ap-northeast-3:
//AMI: ami-098436f864fd2d3b0
//ap-south-1:
//AMI: ami-0fbb5026d37f5f347
//ap-southeast-1:
//AMI: ami-0108ebda02138676e
//ap-southeast-2:
//AMI: ami-09a2afe68765e6248
//ca-central-1:
//AMI: ami-04ecbc5a516929022
//eu-central-1:
//AMI: ami-01ad80fd556799f44
//eu-north-1:
//AMI: ami-003e63396fb3cf41f
//eu-west-1:
//AMI: ami-0cbbbe6104042e885
//eu-west-2:
//AMI: ami-035d9fe64fff63b27
//eu-west-3:
//AMI: ami-050cc42ff95d4ed9d
//sa-east-1:
//AMI: ami-0f68db8769ddddc2e
//us-east-1:
//AMI: ami-03a8890aa01a8eaff
//us-east-2:
//AMI: ami-0280e3450a5a28e52
//us-west-1:
//AMI: ami-0057dd7b23e1c372e
//us-west-2:
//AMI: ami-036aba10f939dbcbb

var TypeToString = map[Type]string{
	EuWest1:  "eu-west-1",
	ApSouth1: "ap-south-1",
}

func (a Type) String() string {
	return TypeToString[a]
}
