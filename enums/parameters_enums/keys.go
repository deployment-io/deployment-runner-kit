package parameters_enums

/*
RepoCloneUrl
RepoName
*/

type Key uint

const (
	RepoCloneUrl          Key = 1
	RepoName              Key = 2
	RepoBranch            Key = 3
	RepoProviderToken     Key = 4
	RepoDirectoryPath     Key = 5
	RepoGitProvider       Key = 6
	LoggerType            Key = 7
	OrganizationID        Key = 8
	DeploymentID          Key = 9
	BuildID               Key = 10
	EnvironmentID         Key = 11
	RootDirectory         Key = 12
	BuildCommand          Key = 13
	PublishDirectory      Key = 14
	NodeVersion           Key = 15
	EnvironmentVariables  Key = 16
	EnvironmentFiles      Key = 17
	Region                Key = 18
	CloudfrontID          Key = 19
	CommitHash            Key = 20
	Cpu                   Key = 21
	Memory                Key = 22
	Runtime               Key = 23
	DockerImageNameAndTag Key = 24
	DockerBuildArgs       Key = 25
	VpcID                 Key = 26
	InternetGatewayID     Key = 27
	Subnets               Key = 28
	ElasticIPAllocationID Key = 29
	NatGatewayID          Key = 30
	RouteTables           Key = 31
)

var keyToString = map[Key]string{
	RepoCloneUrl:          "repo clone url",
	RepoName:              "repo name",
	RepoBranch:            "repo branch",
	RepoProviderToken:     "repo provider token",
	RepoDirectoryPath:     "repo directory path",
	RepoGitProvider:       "repo git provider",
	LoggerType:            "logger type",
	OrganizationID:        "organization id",
	DeploymentID:          "deployment id",
	BuildID:               "build id",
	EnvironmentID:         "environment id",
	RootDirectory:         "root directory",
	BuildCommand:          "build command",
	PublishDirectory:      "publish directory",
	NodeVersion:           "node version",
	EnvironmentVariables:  "environment variables",
	EnvironmentFiles:      "environment files",
	Region:                "region",
	CloudfrontID:          "cloudfront",
	CommitHash:            "commit hash",
	Cpu:                   "cpu",
	Memory:                "memory",
	Runtime:               "runtime",
	DockerImageNameAndTag: "docker image name and tag",
	DockerBuildArgs:       "docker build arguments",
	VpcID:                 "vpc id",
	InternetGatewayID:     "internet gateway id",
	Subnets:               "subnets map",
	ElasticIPAllocationID: "elastic ip allocation id",
	NatGatewayID:          "nat gateway id",
	RouteTables:           "route tables",
}

func (k Key) String() string {
	s, ok := keyToString[k]
	if !ok {
		return ""
	}
	return s
}
