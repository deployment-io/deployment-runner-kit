package parameters_enums

import "fmt"

/*
RepoCloneUrl
RepoName
*/

type Key uint

const (
	RepoCloneUrl            Key = 1
	RepoName                Key = 2
	RepoBranch              Key = 3
	RepoProviderToken       Key = 4
	RepoDirectoryPath       Key = 5
	RepoGitProvider         Key = 6
	LoggerType              Key = 7
	OrganizationID          Key = 8
	DeploymentID            Key = 9
	BuildID                 Key = 10
	EnvironmentID           Key = 11
	RootDirectory           Key = 12
	BuildCommand            Key = 13
	PublishDirectory        Key = 14
	NodeVersion             Key = 15
	EnvironmentVariables    Key = 16 //map[string][string]
	EnvironmentFiles        Key = 17 //map[string][string]
	Region                  Key = 18
	CloudfrontID            Key = 19
	CommitHash              Key = 20
	Cpu                     Key = 21
	Memory                  Key = 22
	Runtime                 Key = 23
	DockerImageNameAndTag   Key = 24
	DockerBuildArgs         Key = 25
	VpcID                   Key = 26
	InternetGatewayID       Key = 27
	Subnets                 Key = 28 //map[string]string
	ElasticIPAllocations    Key = 29 //map[string]string
	NatGateways             Key = 30 //map[string]string
	RouteTables             Key = 31 //map[string]string
	EcsClusterArn           Key = 32
	VpcCidr                 Key = 33
	Port                    Key = 34
	HealthCheckPath         Key = 35
	PublicSubnets           Key = 36 //[]string/primitive.A
	EcrRepositoryUri        Key = 37
	EcrRepositoryUriWithTag Key = 38
	EcsTaskExecutionRole    Key = 39
	PrivateSubnets          Key = 40 //[]string/primitive.A
	AlbSecurityGroupId      Key = 41
	LoadBalancerArn         Key = 42
	TargetGroupArn          Key = 43
	EcsServiceArn           Key = 44
	TaskDefinitionArn       Key = 45
	InstallationId          Key = 46
)

var keyToString = map[Key]string{
	RepoCloneUrl:            "repo clone url",
	RepoName:                "repo name",
	RepoBranch:              "repo branch",
	RepoProviderToken:       "repo provider token",
	RepoDirectoryPath:       "repo directory path",
	RepoGitProvider:         "repo git provider",
	LoggerType:              "logger type",
	OrganizationID:          "organization id",
	DeploymentID:            "deployment id",
	BuildID:                 "build id",
	EnvironmentID:           "environment id",
	RootDirectory:           "root directory",
	BuildCommand:            "build command",
	PublishDirectory:        "publish directory",
	NodeVersion:             "node version",
	EnvironmentVariables:    "environment variables",
	EnvironmentFiles:        "environment files",
	Region:                  "region",
	CloudfrontID:            "cloudfront",
	CommitHash:              "commit hash",
	Cpu:                     "cpu",
	Memory:                  "memory",
	Runtime:                 "runtime",
	DockerImageNameAndTag:   "docker image name and tag",
	DockerBuildArgs:         "docker build arguments",
	VpcID:                   "vpc id",
	InternetGatewayID:       "internet gateway id",
	Subnets:                 "subnets map",
	ElasticIPAllocations:    "elastic ip allocations map",
	NatGateways:             "nat gateways map",
	RouteTables:             "route tables",
	EcsClusterArn:           "ecs cluster arn",
	VpcCidr:                 "vpc cidr",
	Port:                    "port",
	HealthCheckPath:         "health check path",
	PublicSubnets:           "public subnets",
	EcrRepositoryUri:        "ecr repository uri",
	EcrRepositoryUriWithTag: "ecr repository with tag",
	EcsTaskExecutionRole:    "ecs task execution role",
	PrivateSubnets:          "private subnets",
	AlbSecurityGroupId:      "alb security group id",
	LoadBalancerArn:         "load balancer arn",
	TargetGroupArn:          "target group arn",
	EcsServiceArn:           "ecs service arn",
	TaskDefinitionArn:       "task definition arn",
	InstallationId:          "installation id",
}

func (k Key) String() string {
	s, ok := keyToString[k]
	if !ok {
		return ""
	}
	return s
}

var keyMap = map[Key]string{
	RepoCloneUrl:            "1",
	RepoName:                "2",
	RepoBranch:              "3",
	RepoProviderToken:       "4",
	RepoDirectoryPath:       "5",
	RepoGitProvider:         "6",
	LoggerType:              "7",
	OrganizationID:          "8",
	DeploymentID:            "9",
	BuildID:                 "10",
	EnvironmentID:           "11",
	RootDirectory:           "12",
	BuildCommand:            "13",
	PublishDirectory:        "14",
	NodeVersion:             "15",
	EnvironmentVariables:    "16", //map[string][string]
	EnvironmentFiles:        "17", //map[string][string]
	Region:                  "18",
	CloudfrontID:            "19",
	CommitHash:              "20",
	Cpu:                     "21",
	Memory:                  "22",
	Runtime:                 "23",
	DockerImageNameAndTag:   "24",
	DockerBuildArgs:         "25",
	VpcID:                   "26",
	InternetGatewayID:       "27",
	Subnets:                 "28", //map[string]string
	ElasticIPAllocations:    "29", //map[string]string
	NatGateways:             "30", //map[string]string
	RouteTables:             "31", //map[string]string
	EcsClusterArn:           "32",
	VpcCidr:                 "33",
	Port:                    "34",
	HealthCheckPath:         "35",
	PublicSubnets:           "36", //[]string/primitive.A
	EcrRepositoryUri:        "37",
	EcrRepositoryUriWithTag: "38",
	EcsTaskExecutionRole:    "39",
	PrivateSubnets:          "40", //[]string/primitive.A
	AlbSecurityGroupId:      "41",
	LoadBalancerArn:         "42",
	TargetGroupArn:          "43",
	EcsServiceArn:           "44",
	TaskDefinitionArn:       "45",
	InstallationId:          "46",
}

func (k Key) Key() (string, error) {
	s, ok := keyMap[k]
	if !ok {
		return "", fmt.Errorf("error finding key for %s", k.String())
	}
	return s, nil
}
