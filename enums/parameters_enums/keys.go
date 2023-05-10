package parameters_enums

/*
RepoCloneUrl
RepoName
*/

type Key uint

const (
	RepoCloneUrl         Key = 1
	RepoName             Key = 2
	RepoBranch           Key = 3
	RepoProviderToken    Key = 4
	RepoDirectoryPath    Key = 5
	RepoGitProvider      Key = 6
	LoggerType           Key = 7
	OrganizationID       Key = 8
	DeploymentID         Key = 9
	BuildID              Key = 10
	EnvironmentID        Key = 11
	RootDirectory        Key = 12
	BuildCommand         Key = 13
	PublishDirectory     Key = 14
	NodeVersion          Key = 15
	EnvironmentVariables Key = 16
	EnvironmentFiles     Key = 17
)

var keyToString = map[Key]string{
	RepoCloneUrl:         "repo clone url",
	RepoName:             "repo name",
	RepoBranch:           "repo branch",
	RepoProviderToken:    "repo provider token",
	RepoDirectoryPath:    "repo directory path",
	RepoGitProvider:      "repo git provider",
	LoggerType:           "logger type",
	OrganizationID:       "organization id",
	DeploymentID:         "deployment id",
	BuildID:              "build id",
	EnvironmentID:        "environment id",
	RootDirectory:        "root directory",
	BuildCommand:         "build command",
	PublishDirectory:     "publish directory",
	NodeVersion:          "node version",
	EnvironmentVariables: "environment variables",
	EnvironmentFiles:     "environment files",
}

func (k Key) String() string {
	s, ok := keyToString[k]
	if !ok {
		return ""
	}
	return s
}
