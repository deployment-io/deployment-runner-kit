package previews

import "github.com/deployment-io/deployment-runner-kit/types"

type GetPreviewsArgsV1 struct {
	types.AuthArgsV1
	IDs []string
}

type GetPreviewDtoV1 struct {
	ID                               string
	CloudfrontDistributionID         string
	CloudfrontDistributionArn        string
	CloudfrontDistributionDomainName string
}

type GetPreviewsReplyV1 struct {
	Previews []GetPreviewDtoV1
}
