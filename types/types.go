package types

const (
	True  = "true"
	False = "false"
)

// AuthArgsV1 represents args for auth
type AuthArgsV1 struct {
	OrganizationID string
	Token          string
}
