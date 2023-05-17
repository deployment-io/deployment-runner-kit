package deployments

import (
	"github.com/deployment-io/deployment-runner-kit/types"
)

type UpdateDeploymentDtoV1 struct {
	ID                               string
	CloudfrontDistributionID         string
	CloudfrontDistributionArn        string
	CloudfrontDistributionDomainName string
}

type UpdateDeploymentsArgsV1 struct {
	types.AuthArgsV1
	Deployments []UpdateDeploymentDtoV1
}

type UpdateDeploymentsReplyV1 struct {
	Done bool
}
