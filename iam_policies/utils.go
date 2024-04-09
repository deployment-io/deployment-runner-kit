package iam_policies

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/deployment-io/deployment-runner-kit/cloud_api_clients"
	"github.com/deployment-io/deployment-runner-kit/enums/iam_policy_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/parameters_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/region_enums"
	"github.com/deployment-io/deployment-runner-kit/iam_policies/aws_policy_schema"
	"net/url"
)

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

func AddAwsPolicyForDeploymentRunner(policyType iam_policy_enums.Type, osStr, cpuStr, organizationID, runnerRegion string) error {
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

	//parse policy document data
	policyDocumentEncoded := aws.ToString(getRolePolicyOutput.PolicyDocument)
	if len(policyDocumentEncoded) == 0 {
		return fmt.Errorf("got empty policy document from AWS")
	}
	policyDocument, err := url.QueryUnescape(policyDocumentEncoded)
	if err != nil {
		return fmt.Errorf("error decoding policy document: %s", err)
	}
	var policyDocumentData aws_policy_schema.PolicyDocumentData
	err = json.Unmarshal([]byte(policyDocument), &policyDocumentData)
	if err != nil {
		return fmt.Errorf("error parsing policy data: %s", err)
	}

	newActions, err := policyType.GetPolicyDataActions()
	if err != nil {
		return fmt.Errorf("error finding new actions for policy type: %s", err)
	}

	statementIndex := -1
	var oldActionsSet = make(map[string]bool)
	for index, statement := range policyDocumentData.Statement {
		if len(statement.Sid) > 0 && statement.Sid == "deploymentIoAdded" {
			for _, action := range statement.Action {
				//create a set for old actions
				oldActionsSet[action] = true
			}
			statementIndex = index
			break
		}
	}

	if statementIndex == -1 {
		//add and append new statement
		policyDocumentData.Statement = append(policyDocumentData.Statement, aws_policy_schema.Statement{
			Sid:      "deploymentIoAdded",
			Effect:   "Allow",
			Action:   []string{},
			Resource: []string{"*"},
		})
		statementIndex = len(policyDocumentData.Statement) - 1
	}

	var newActionsUpdateList []string
	for _, newAction := range newActions {
		_, exists := oldActionsSet[newAction]
		if !exists {
			newActionsUpdateList = append(newActionsUpdateList, newAction)
		}
	}

	if len(newActionsUpdateList) > 0 {
		policyDocumentData.Statement[statementIndex].Action = append(policyDocumentData.Statement[statementIndex].Action, newActionsUpdateList...)
		//add inline policy
		newPolicyDocument, err := json.Marshal(policyDocumentData)
		if err != nil {
			return err
		}
		_, err = iamClient.PutRolePolicy(context.TODO(), &iam.PutRolePolicyInput{
			PolicyDocument: aws.String(string(newPolicyDocument)),
			PolicyName:     aws.String(runnerPolicyName),
			RoleName:       aws.String(runnerTaskRoleName),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
