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
	default:
		return nil, fmt.Errorf("error finding policy data actions")
	}
}
