package types

const (
	True  = "true"
	False = "false"
)

// AuthArgsV1 represents args for auth
type AuthArgsV1 struct {
	OrganizationID string
	Token          string
	DockerImage    string
}

type AuthArgsV2 struct {
	OrganizationID string
	Token          string
	DockerImage    string
	CloudAccountID string
	RunnerRegion   string
	GoArch         string
	GoOS           string
}
