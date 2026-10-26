package oauth

import (
	"fmt"
	"github.com/deployment-io/deployment-runner-kit/types"
)

var ErrRefreshInProcess = fmt.Errorf("a client is already in process of refreshing git token for installation")

type RefreshGitProviderTokenArgsV1 struct {
	types.AuthArgsV1
	InstallationID string
	// RejectedTokenSHA256 is the lower-case hex SHA-256 of the token the git
	// provider just rejected, or "" for an ordinary request (asking for a
	// token to start with). When it matches the stored token the server mints
	// a new one whatever the stored expiry; when it doesn't, another caller
	// already replaced the rejected token and the stored one is returned. The
	// token itself is never sent back. Older servers ignore the field.
	RejectedTokenSHA256 string
}

type RefreshGitProviderTokenDtoV1 struct {
	Token string
}
