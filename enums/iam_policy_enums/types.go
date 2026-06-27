package iam_policy_enums

import "fmt"

type Type uint

const (
	AwsDeploymentRunnerControllerStart = iota + 1
	AwsStaticSiteDeployment
	AwsWebServiceDeployment
	AwsLogs
	AwsEcrUpload
	AwsCertificateManager
	AwsSecretsManager
	AwsRdsDeployment
	AwsInfraContext
	MaxType //always add new types before MaxType
)

var typeToString = map[Type]string{
	//AwsCloudfront: "AwsCloudfront",
}

func (t Type) String() string {
	return typeToString[t]
}

func (t Type) GetPolicyDataActions() ([]string, error) {
	switch t {
	case AwsDeploymentRunnerControllerStart:
		return []string{
			"application-autoscaling:*",
			"autoscaling:*",
			"ec2:*",
			"ecs:*",
			"secretsmanager:*",
		}, nil
	case AwsStaticSiteDeployment:
		return []string{
			"cloudfront:*",
			"s3:*",
			"secretsmanager:*",
			"route53:*",
		}, nil
	case AwsWebServiceDeployment:
		return []string{
			"application-autoscaling:*",
			"autoscaling:*",
			"ec2:*",
			"ecs:*",
			"elasticloadbalancing:*",
			"secretsmanager:*",
			"servicediscovery:*",
			"route53:*",
		}, nil
	case AwsEcrUpload:
		return []string{
			"ecr:*",
		}, nil
	case AwsLogs:
		return []string{
			"cloudwatch:*",
			"logs:*",
			"events:*",
		}, nil
	case AwsCertificateManager:
		return []string{
			"acm:*",
			"elasticloadbalancing:*",
		}, nil
	case AwsSecretsManager:
		return []string{
			"secretsmanager:*",
		}, nil
	case AwsRdsDeployment:
		return []string{
			"rds:*",
		}, nil
	case AwsInfraContext:
		// Read access for the infra-context connectors (Explore/Context observed layer). ecs:* for
		// the AWS ECS source (ListClusters/ListServices/DescribeServices/DescribeTaskDefinition).
		// Matches the codebase's service:* bundle convention and the deploy bundle's ecs:*, so it's a
		// no-op on a runner that has already deployed (no redundant policy-propagation wait) and only
		// self-provisions on a context-only runner that has never deployed. Add ecr:* here when image
		// label / digest recovery (OCI org.opencontainers.image.source) lands.
		return []string{
			"ecs:*",
		}, nil
	default:
		return nil, fmt.Errorf("error finding policy data actions")
	}
}
