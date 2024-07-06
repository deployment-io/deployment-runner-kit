package previews

import (
	"github.com/deployment-io/deployment-runner-kit/enums/build_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/deployment_enums"
	"github.com/deployment-io/deployment-runner-kit/types"
)

type UpdatePreviewDtoV1 struct {
	ID                               string
	BuildTs                          int64
	CommitHash                       string
	CommitMessage                    string
	Status                           build_enums.Status
	ErrorMessage                     string
	TaskDefinitionArn                string
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
	DnsName                          string
	DeletionState                    deployment_enums.DeletionState
	NamespaceName                    string
}

type UpdatePreviewsArgsV1 struct {
	types.AuthArgsV1
	Previews []UpdatePreviewDtoV1
}

type UpdatePreviewsReplyV1 struct {
	Done bool
}
