package parameters_enums

/*
RepoCloneUrl
RepoName
*/

type Key uint

const (
	RepoCloneUrl      Key = 1
	RepoName          Key = 2
	RepoBranch        Key = 3
	RepoProviderToken Key = 4
	RepoDirectoryName Key = 5
	RepoGitProvider   Key = 6
	LoggerType        Key = 7
	OrganizationID    Key = 8
	DeploymentID      Key = 9
	BuildID           Key = 10
)

//var keyToString = map[Key]string{
//	CheckoutRepo: "Checkout Repo",
//}
//
//func (t Key) String() string {
//	return keyToString[t]
//}
