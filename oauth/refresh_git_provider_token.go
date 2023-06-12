package oauth

import (
	"fmt"
	"github.com/deployment-io/deployment-runner-kit/types"
)

var ErrRefreshInProcess = fmt.Errorf("a client is already in process of refreshing git token for installation")

type RefreshGitProviderTokenArgsV1 struct {
	types.AuthArgsV1
	InstallationID string
}

type RefreshGitProviderTokenDtoV1 struct {
	Token string
}
