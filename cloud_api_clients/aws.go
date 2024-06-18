package cloud_api_clients

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/deployment-io/deployment-runner-kit/enums/parameters_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/region_enums"
	"github.com/deployment-io/deployment-runner-kit/jobs"
)

func GetS3Client(parameters map[string]interface{}) (*s3.Client, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	region, err := jobs.GetParameterValue[int64](parameters, parameters_enums.Region)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = region_enums.Type(region).String()
	})

	return s3Client, nil
}

func GetEC2Client(parameters map[string]interface{}) (*ec2.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	region, err := jobs.GetParameterValue[int64](parameters, parameters_enums.Region)
	if err != nil {
		return nil, err
	}

	ec2Client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		o.Region = region_enums.Type(region).String()
	})
	return ec2Client, nil
}

func GetEcsClient(parameters map[string]interface{}) (*ecs.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	region, err := jobs.GetParameterValue[int64](parameters, parameters_enums.Region)
	if err != nil {
		return nil, err
	}

	ecsClient := ecs.NewFromConfig(cfg, func(o *ecs.Options) {
		o.Region = region_enums.Type(region).String()
	})
	return ecsClient, nil
}

func GetEcsClientFromRegion(region string) (*ecs.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	ecsClient := ecs.NewFromConfig(cfg, func(o *ecs.Options) {
		o.Region = region
	})
	return ecsClient, nil
}

func GetElbClient(parameters map[string]interface{}) (*elasticloadbalancingv2.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	region, err := jobs.GetParameterValue[int64](parameters, parameters_enums.Region)
	if err != nil {
		return nil, err
	}

	elbClient := elasticloadbalancingv2.NewFromConfig(cfg, func(options *elasticloadbalancingv2.Options) {
		options.Region = region_enums.Type(region).String()
	})

	return elbClient, nil
}

func GetEcrClient(parameters map[string]interface{}) (*ecr.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	region, err := jobs.GetParameterValue[int64](parameters, parameters_enums.Region)
	if err != nil {
		return nil, err
	}

	ecrClient := ecr.NewFromConfig(cfg, func(options *ecr.Options) {
		options.Region = region_enums.Type(region).String()
	})

	return ecrClient, nil
}

func GetIamClient(parameters map[string]interface{}) (*iam.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	region, err := jobs.GetParameterValue[int64](parameters, parameters_enums.Region)
	if err != nil {
		return nil, err
	}

	iamClient := iam.NewFromConfig(cfg, func(options *iam.Options) {
		options.Region = region_enums.Type(region).String()
	})

	return iamClient, nil
}

func GetCloudfrontClient(parameters map[string]interface{}, region string) (*cloudfront.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	//region, err := jobs.GetParameterValue[int64](parameters, parameters_enums.Region)
	//if err != nil {
	//	return nil, err
	//}

	cloudfrontClient := cloudfront.NewFromConfig(cfg, func(options *cloudfront.Options) {
		options.Region = region
	})

	return cloudfrontClient, nil
}

func GetAcmClient(parameters map[string]interface{}) (*acm.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	region, err := jobs.GetParameterValue[int64](parameters, parameters_enums.CertificateRegion)
	if err != nil {
		return nil, err
	}

	acmClient := acm.NewFromConfig(cfg, func(options *acm.Options) {
		options.Region = region_enums.Type(region).String()
	})

	return acmClient, nil
}

func GetStsClient(region string) (*sts.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	stsClient := sts.NewFromConfig(cfg, func(options *sts.Options) {
		options.Region = region
	})

	return stsClient, nil
}

func GetAsgClient(region string) (*autoscaling.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	asgClient := autoscaling.NewFromConfig(cfg, func(options *autoscaling.Options) {
		options.Region = region
	})

	return asgClient, nil
}

func GetCloudwatchClient(parameters map[string]interface{}) (*cloudwatch.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	region, err := jobs.GetParameterValue[int64](parameters, parameters_enums.Region)
	if err != nil {
		return nil, err
	}

	cloudwatchClient := cloudwatch.NewFromConfig(cfg, func(options *cloudwatch.Options) {
		options.Region = region_enums.Type(region).String()
	})

	return cloudwatchClient, nil
}
