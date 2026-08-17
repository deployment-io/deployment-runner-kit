package iam_policies

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/deployment-io/deployment-runner-kit/cloud_api_clients"
	"github.com/deployment-io/deployment-runner-kit/enums/iam_policy_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/parameters_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/region_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/runner_enums"
	"github.com/deployment-io/deployment-runner-kit/iam_policies/aws_policy_schema"
	"net/url"
	"time"
)

// deploymentIoAddedSid marks the one statement in the runner's inline policy that we own.
const deploymentIoAddedSid = "deploymentIoAdded"

func getDeploymentRunnerTaskRoleName(osStr, cpuStr, organizationID, region string) string {
	//dr-task-role-<osCpuStr>-<orgid>-<region>
	osCpuStr := fmt.Sprintf("%s%s", osStr, cpuStr)
	taskRoleName := fmt.Sprintf("dr-task-role-%s", osCpuStr)
	return fmt.Sprintf("%s-%s-%s", taskRoleName, organizationID, region)
}

func getDeploymentRunnerPolicyName(osStr, cpuStr, organizationID, region string) string {
	//dr-policy-<osCpuStr>-<orgid>-<region>
	osCpuStr := fmt.Sprintf("%s%s", osStr, cpuStr)
	policyName := fmt.Sprintf("dr-policy-%s", osCpuStr)
	return fmt.Sprintf("%s-%s-%s", policyName, organizationID, region)
}

func shouldAddPolicy(mode runner_enums.Mode, cloud runner_enums.TargetCloud) bool {
	return (mode == runner_enums.Saas || mode == runner_enums.AwsEcs) && cloud == runner_enums.AwsCloud
}

func AddAwsPolicyForDeploymentRunner(policyType iam_policy_enums.Type, osStr, cpuStr, organizationID, runnerRegion string,
	mode runner_enums.Mode, cloud runner_enums.TargetCloud) error {
	if !shouldAddPolicy(mode, cloud) {
		return nil
	}
	//check and add policy
	regionKey, err := parameters_enums.Region.Key()
	if err != nil {
		return err
	}
	regionType, err := region_enums.GetType(runnerRegion)
	if err != nil {
		return err
	}
	var parameters = map[string]interface{}{
		regionKey: int64(regionType),
	}
	iamClient, err := cloud_api_clients.GetIamClient(parameters)
	if err != nil {
		return err
	}
	runnerTaskRoleName := getDeploymentRunnerTaskRoleName(osStr, cpuStr, organizationID, runnerRegion)
	runnerPolicyName := getDeploymentRunnerPolicyName(osStr, cpuStr, organizationID, runnerRegion)

	//get inline policy
	getRolePolicyOutput, err := iamClient.GetRolePolicy(context.TODO(), &iam.GetRolePolicyInput{
		PolicyName: aws.String(runnerPolicyName),
		RoleName:   aws.String(runnerTaskRoleName),
	})
	if err != nil {
		return err
	}

	//decode policy document
	policyDocumentEncoded := aws.ToString(getRolePolicyOutput.PolicyDocument)
	if len(policyDocumentEncoded) == 0 {
		return fmt.Errorf("got empty policy document from AWS")
	}
	policyDocument, err := url.QueryUnescape(policyDocumentEncoded)
	if err != nil {
		return fmt.Errorf("error decoding policy document: %s", err)
	}
	newActions, err := policyType.GetPolicyDataActions()
	if err != nil {
		return fmt.Errorf("error finding new actions for policy type: %s", err)
	}

	//the whole document is written back, so anything we didn't author has to survive untouched
	newPolicyDocument, changed, err := aws_policy_schema.AddActionsToSid([]byte(policyDocument), deploymentIoAddedSid, newActions)
	if err != nil {
		return fmt.Errorf("error parsing policy data: %s", err)
	}

	if changed {
		//add inline policy
		_, err = iamClient.PutRolePolicy(context.TODO(), &iam.PutRolePolicyInput{
			PolicyDocument: aws.String(string(newPolicyDocument)),
			PolicyName:     aws.String(runnerPolicyName),
			RoleName:       aws.String(runnerTaskRoleName),
		})
		if err != nil {
			return err
		}
		//sleep for 60 seconds since new policies are added. AWS is not fast to update.
		time.Sleep(60 * time.Second)
	}
	return nil
}
