package parameters_enums

import "fmt"

/*
RepoCloneUrl
RepoName
*/

type Key uint

const (
	RepoCloneUrl               Key = 1
	RepoName                   Key = 2
	RepoBranch                 Key = 3
	RepoProviderToken          Key = 4
	RepoDirectoryPath          Key = 5
	RepoGitProvider            Key = 6
	LoggerType                 Key = 7
	OrganizationIDNamespace    Key = 8
	DeploymentID               Key = 9
	BuildID                    Key = 10
	EnvironmentID              Key = 11
	RootDirectory              Key = 12
	BuildCommand               Key = 13
	PublishDirectory           Key = 14
	NodeVersion                Key = 15
	EnvironmentVariables       Key = 16 //map[string][string]
	EnvironmentFiles           Key = 17 //map[string][string]
	Region                     Key = 18
	CloudfrontID               Key = 19
	CommitHash                 Key = 20
	Cpu                        Key = 21
	Memory                     Key = 22
	Runtime                    Key = 23
	DockerImageNameAndTag      Key = 24
	DockerBuildArgs            Key = 25
	VpcID                      Key = 26
	InternetGatewayID          Key = 27
	Subnets                    Key = 28 //map[string]string
	ElasticIPAllocations       Key = 29 //map[string]string
	NatGateways                Key = 30 //map[string]string
	RouteTables                Key = 31 //map[string]string
	EcsClusterArn              Key = 32
	VpcCidr                    Key = 33
	Port                       Key = 34
	HealthCheckPath            Key = 35
	PublicSubnets              Key = 36 //[]string/primitive.A
	EcrRepositoryUri           Key = 37
	DockerRepositoryUriWithTag Key = 38
	EcsTaskExecutionRoleArn    Key = 39
	PrivateSubnets             Key = 40 //[]string/primitive.A
	AlbSecurityGroupId         Key = 41
	LoadBalancerArn            Key = 42
	TargetGroupArn             Key = 43
	EcsServiceArn              Key = 44
	TaskDefinitionArn          Key = 45
	InstallationId             Key = 46

	//CpuArchitecture         Key = 47

	ResponseHeaders              Key = 48
	Domains                      Key = 49
	RedirectDomain               Key = 50
	AcmCertificateArn            Key = 51
	CertificateRegion            Key = 52
	CertificateDomain            Key = 53
	CertificateID                Key = 54
	ParentDomain                 Key = 55
	CreateMultipleNats           Key = 56
	ErrorPages                   Key = 57
	PreviewID                    Key = 58
	IsPreview                    Key = 59
	EcsClusterName               Key = 60
	ShouldSendChatOutput         Key = 61
	JobID                        Key = 62
	SecretName                   Key = 63
	SecretValue                  Key = 64
	IsPrivateRegistry            Key = 65
	DeployedFromImage            Key = 66
	IsPrivateService             Key = 67
	DeploymentName               Key = 68
	EnvironmentName              Key = 69
	RdsEngine                    Key = 70
	CpuMemoryRDSInstance         Key = 71
	AllocatedStorage             Key = 72
	MaxAllocatedStorage          Key = 73
	UseMultiAz                   Key = 74
	UseDeletionProtection        Key = 75
	ApplyRdsChangesImmediately   Key = 76
	ShouldUpdatePassword         Key = 77
	DockerFilePath               Key = 78
	StartCommand                 Key = 79
	NixPacksUsernameKey          Key = 80
	NixPacksPasswordKey          Key = 81
	OrganizationIdFromJob        Key = 82
	AutomationData               Key = 83 //AutomationDtoV1
	AutomationID                 Key = 84
	SmtpUsername                 Key = 85
	SmtpPassword                 Key = 86
	SmtpHost                     Key = 87
	SmtpPort                     Key = 88
	EmailToolFromAddress         Key = 89
	DebugAutomation              Key = 90
	DebugOpenAICallsInAutomation Key = 91
	AgentData                    Key = 92 //AgentDtoV1
	AgentID                      Key = 93
	DebugAgent                   Key = 94
	DebugOpenAICallsInAgent      Key = 95
)

var keyToString = map[Key]string{
	RepoCloneUrl:               "repo clone url",
	RepoName:                   "repo name",
	RepoBranch:                 "repo branch",
	RepoProviderToken:          "repo provider token",
	RepoDirectoryPath:          "repo directory path",
	RepoGitProvider:            "repo git provider",
	LoggerType:                 "logger type",
	OrganizationIDNamespace:    "organization id namespace",
	DeploymentID:               "deployment id",
	BuildID:                    "build id",
	EnvironmentID:              "environment id",
	RootDirectory:              "root directory",
	BuildCommand:               "build command",
	PublishDirectory:           "publish directory",
	NodeVersion:                "node version",
	EnvironmentVariables:       "environment variables",
	EnvironmentFiles:           "environment files",
	Region:                     "region",
	CloudfrontID:               "cloudfront",
	CommitHash:                 "commit hash",
	Cpu:                        "cpu",
	Memory:                     "memory",
	Runtime:                    "runtime",
	DockerImageNameAndTag:      "docker image name and tag",
	DockerBuildArgs:            "docker build arguments",
	VpcID:                      "vpc id",
	InternetGatewayID:          "internet gateway id",
	Subnets:                    "subnets map",
	ElasticIPAllocations:       "elastic ip allocations map",
	NatGateways:                "nat gateways map",
	RouteTables:                "route tables",
	EcsClusterArn:              "ecs cluster arn",
	VpcCidr:                    "vpc cidr",
	Port:                       "port",
	HealthCheckPath:            "health check path",
	PublicSubnets:              "public subnets",
	EcrRepositoryUri:           "ecr repository uri",
	DockerRepositoryUriWithTag: "docker repository with tag",
	EcsTaskExecutionRoleArn:    "ecs task execution role arn",
	PrivateSubnets:             "private subnets",
	AlbSecurityGroupId:         "alb security group id",
	LoadBalancerArn:            "load balancer arn",
	TargetGroupArn:             "target group arn",
	EcsServiceArn:              "ecs service arn",
	TaskDefinitionArn:          "task definition arn",
	InstallationId:             "installation id",
	//CpuArchitecture:         "cpu architecture",
	ResponseHeaders:              "response headers",
	Domains:                      "domains",
	RedirectDomain:               "redirect domain",
	AcmCertificateArn:            "ACM certificate Arn",
	CertificateRegion:            "Certificate region",
	CertificateDomain:            "Certificate domain",
	CertificateID:                "Certificate ID",
	ParentDomain:                 "Parent domain",
	CreateMultipleNats:           "Create multiple nats",
	ErrorPages:                   "error pages",
	PreviewID:                    "Preview ID",
	IsPreview:                    "Is the job of type preview",
	EcsClusterName:               "ECS cluster name",
	ShouldSendChatOutput:         "Should output be sent back to the chat",
	JobID:                        "job Id",
	SecretName:                   "secret name",
	SecretValue:                  "secret value",
	IsPrivateRegistry:            "is private registry",
	DeployedFromImage:            "deployed from image",
	IsPrivateService:             "is private service",
	DeploymentName:               "deployment name",
	EnvironmentName:              "environment name",
	RdsEngine:                    "rds engine",
	CpuMemoryRDSInstance:         "cpu memory rds instance",
	AllocatedStorage:             "allocated storage",
	MaxAllocatedStorage:          "max allocated storage",
	UseMultiAz:                   "use multi az",
	UseDeletionProtection:        "use deletion protection",
	ApplyRdsChangesImmediately:   "apply rds changes immediately",
	ShouldUpdatePassword:         "should update password",
	DockerFilePath:               "docker file path",
	StartCommand:                 "start command",
	NixPacksUsernameKey:          "nix packs username",
	NixPacksPasswordKey:          "nix packs password",
	OrganizationIdFromJob:        "org id from job",
	AutomationData:               "automation data",
	AutomationID:                 "automation id",
	SmtpUsername:                 "smtp username",
	SmtpPassword:                 "smtp password",
	SmtpHost:                     "smtp host",
	SmtpPort:                     "smtp port",
	EmailToolFromAddress:         "smtp from address",
	DebugAutomation:              "debug automation",
	DebugOpenAICallsInAutomation: "debug open ai calls in automation",
	AgentData:                    "agent data",
	AgentID:                      "agent id",
	DebugAgent:                   "debug agent",
	DebugOpenAICallsInAgent:      "debug open ai calls in agent",
}

func (k Key) String() string {
	s, ok := keyToString[k]
	if !ok {
		return ""
	}
	return s
}

var keyMap = map[Key]string{
	RepoCloneUrl:               "1",
	RepoName:                   "2",
	RepoBranch:                 "3",
	RepoProviderToken:          "4",
	RepoDirectoryPath:          "5",
	RepoGitProvider:            "6",
	LoggerType:                 "7",
	OrganizationIDNamespace:    "8",
	DeploymentID:               "9",
	BuildID:                    "10",
	EnvironmentID:              "11",
	RootDirectory:              "12",
	BuildCommand:               "13",
	PublishDirectory:           "14",
	NodeVersion:                "15",
	EnvironmentVariables:       "16", //map[string][string]
	EnvironmentFiles:           "17", //map[string][string]
	Region:                     "18",
	CloudfrontID:               "19",
	CommitHash:                 "20",
	Cpu:                        "21",
	Memory:                     "22",
	Runtime:                    "23",
	DockerImageNameAndTag:      "24",
	DockerBuildArgs:            "25",
	VpcID:                      "26",
	InternetGatewayID:          "27",
	Subnets:                    "28", //map[string]string
	ElasticIPAllocations:       "29", //map[string]string
	NatGateways:                "30", //map[string]string
	RouteTables:                "31", //map[string]string
	EcsClusterArn:              "32",
	VpcCidr:                    "33",
	Port:                       "34",
	HealthCheckPath:            "35",
	PublicSubnets:              "36", //[]string/primitive.A
	EcrRepositoryUri:           "37",
	DockerRepositoryUriWithTag: "38",
	EcsTaskExecutionRoleArn:    "39",
	PrivateSubnets:             "40", //[]string/primitive.A
	AlbSecurityGroupId:         "41",
	LoadBalancerArn:            "42",
	TargetGroupArn:             "43",
	EcsServiceArn:              "44",
	TaskDefinitionArn:          "45",
	InstallationId:             "46",
	//CpuArchitecture:         "47",
	ResponseHeaders:              "48", //[][]string/primitive.A
	Domains:                      "49", //[]string/primitive.A
	RedirectDomain:               "50", //[]string/primitive.A - 0th element is from and 1st element is to
	AcmCertificateArn:            "51",
	CertificateRegion:            "52",
	CertificateDomain:            "53",
	CertificateID:                "54",
	ParentDomain:                 "55",
	CreateMultipleNats:           "56",
	ErrorPages:                   "57", //[][]string/primitive.A
	PreviewID:                    "58",
	IsPreview:                    "59",
	EcsClusterName:               "60",
	ShouldSendChatOutput:         "61",
	JobID:                        "62",
	SecretName:                   "63",
	SecretValue:                  "64",
	IsPrivateRegistry:            "65",
	DeployedFromImage:            "66",
	IsPrivateService:             "67",
	DeploymentName:               "68",
	EnvironmentName:              "69",
	RdsEngine:                    "70",
	CpuMemoryRDSInstance:         "71",
	AllocatedStorage:             "72",
	MaxAllocatedStorage:          "73",
	UseMultiAz:                   "74",
	UseDeletionProtection:        "75",
	ApplyRdsChangesImmediately:   "76",
	ShouldUpdatePassword:         "77",
	DockerFilePath:               "78",
	StartCommand:                 "79",
	NixPacksUsernameKey:          "80",
	NixPacksPasswordKey:          "81",
	OrganizationIdFromJob:        "82",
	AutomationData:               "83",
	AutomationID:                 "84",
	SmtpUsername:                 "85",
	SmtpPassword:                 "86",
	SmtpHost:                     "87",
	SmtpPort:                     "88",
	EmailToolFromAddress:         "89",
	DebugAutomation:              "90",
	DebugOpenAICallsInAutomation: "91",
	AgentData:                    "92",
	AgentID:                      "93",
	DebugAgent:                   "94",
	DebugOpenAICallsInAgent:      "95",
}

func (k Key) Key() (string, error) {
	s, ok := keyMap[k]
	if !ok {
		return "", fmt.Errorf("error finding key for %s", k.String())
	}
	return s, nil
}
