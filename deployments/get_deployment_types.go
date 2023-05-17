package deployments

import "github.com/deployment-io/deployment-runner-kit/types"

type GetDeploymentsArgsV1 struct {
	types.AuthArgsV1
	IDs []string
}

type GetDeploymentDtoV1 struct {
	ID                               string
	CloudfrontDistributionID         string
	CloudfrontDistributionArn        string
	CloudfrontDistributionDomainName string
}

type GetDeploymentsReplyV1 struct {
	Deployments []GetDeploymentDtoV1
}
