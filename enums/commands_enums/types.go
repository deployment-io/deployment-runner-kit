package commands_enums

/*
CheckoutRepo
*/

type Type uint

const (
	CheckoutRepo        Type = 1
	BuildStaticSite     Type = 2
	DeployAwsStaticSite Type = 3
	DeployAwsWebService Type = 4
	BuildDockerImage    Type = 5
	CreateAwsVpc        Type = 6
	CreateEcsCluster    Type = 7
	UploadImageToEcr    Type = 8
)

var typeToString = map[Type]string{
	CheckoutRepo:        "Checkout Repo",
	BuildStaticSite:     "Build Static Site",
	DeployAwsStaticSite: "Deploy AWS Static Site",
	DeployAwsWebService: "Deploy AWS Web Service",
	BuildDockerImage:    "Build docker image",
	CreateAwsVpc:        "Create AWS VPC",
	CreateEcsCluster:    "Create ECS Cluster",
	UploadImageToEcr:    "Upload image to ECR",
}

func (t Type) String() string {
	return typeToString[t]
}
