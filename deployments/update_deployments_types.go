package deployments

import (
	"github.com/deployment-io/deployment-runner-kit/enums/build_enums"
	"github.com/deployment-io/deployment-runner-kit/types"
)

type UpdateDeploymentDtoV1 struct {
	ID                               string
	CloudfrontDistributionID         string
	CloudfrontDistributionArn        string
	CloudfrontDistributionDomainName string
	EcrRepositoryUri                 string
	AlbSecurityGroupId               string
	AlbSecurityIngressRuleId         string
	AlbSecurityEgressRuleId          string
	TargetGroupArn                   string
	ListenerArnPort80                string
	ListenerArnPort443               string
	LoadBalancerArn                  string
	LoadBalancerDns                  string
	EcsServiceArn                    string
	LastDeploymentTs                 int64
	ErrorMessage                     string
	Status                           build_enums.Status
}

type UpdateDeploymentsArgsV1 struct {
	types.AuthArgsV1
	Deployments []UpdateDeploymentDtoV1
}

type UpdateDeploymentsReplyV1 struct {
	Done bool
}
