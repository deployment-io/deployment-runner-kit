package commands_enums

type Type uint

const (
	CheckoutRepo                             Type = 1
	BuildStaticSite                          Type = 2
	DeployAwsStaticSite                      Type = 3
	DeployAwsWebService                      Type = 4
	BuildDockerImage                         Type = 5
	CreateAwsVpc                             Type = 6
	CreateEcsCluster                         Type = 7
	UploadImageToEcr                         Type = 8
	AddAwsStaticSiteResponseHeaders          Type = 9
	UpdateAwsStaticSiteDomains               Type = 10
	DeployAwsCloudfrontViewerRequestFunction Type = 11
	CreateAcmCertificate                     Type = 12
	UpdateAwsWebServiceDomain                Type = 13
	VerifyAcmCertificate                     Type = 14
	DeleteAwsStaticSite                      Type = 15
	DeleteAwsWebService                      Type = 16
	ListCloudWatchMetricsAwsEcsWebService    Type = 17
	CreateSecretAwsSecretManager             Type = 18
	DeployAwsPrivateService                  Type = 19
	DeleteAwsPrivateService                  Type = 20
	DeployAwsRdsDatabase                     Type = 21
	DeleteAwsRdsDatabase                     Type = 22
	BuildNixPacksImage                       Type = 23
	// Deprecated: RunNewAutomation and RunNewAgent belonged to the removed
	// automations feature. Nothing creates or runs these commands any more; the
	// constants stay (never renumber or reuse) because old job documents store them.
	RunNewAutomation     Type = 24
	RunNewAgent          Type = 25
	GetDeploymentLogsAws Type = 26
	RunAgentStep         Type = 27
	CommitAndPush        Type = 28
	OpenPullRequest      Type = 29
	RunAssistantSession  Type = 30
	BuildInfraContext    Type = 31
	MaterializeContext   Type = 32
	// RunReviewStage sits between RunAgentStep and CommitAndPush in a Task
	// Step Job: it reviews the diff the implementer produced against the
	// Task's spec and routes must-fix findings back to the implementer
	// before anything is committed. A Job whose review participation is Off
	// returns immediately, so the chain is identical for every Task that
	// does not use the stage.
	RunReviewStage Type = 33
	// RunQAStage is the QA stage. It sits between RunReviewStage and
	// CommitAndPush in a Task Step Job: a separate agent run deploys the
	// Task's static-site previews and checks the spec's acceptance criteria
	// in a real browser before anything is committed.
	RunQAStage Type = 34
	// RunProductionQA is the only command of a production QA Job: once a
	// deployed Task's deliveries have all ended, one report-only agentbox QA
	// run checks the Task's acceptance criteria against the deployed
	// repositories' production URLs, with a read-only browser. It has no
	// checkout, no Step and no fix round.
	RunProductionQA Type = 35
	// DeleteAwsNamespaces is the only command of a Delete environment
	// namespaces Job, sent to an environment's runner when the environment is
	// deleted: it deletes the environment's Cloud Map private DNS namespaces
	// that no service uses any more.
	DeleteAwsNamespaces Type = 36
)

var typeToString = map[Type]string{
	CheckoutRepo:                             "Checkout Repo",
	BuildStaticSite:                          "Build Static Site",
	DeployAwsStaticSite:                      "Deploy AWS Static Site",
	DeployAwsWebService:                      "Deploy AWS Web Service",
	BuildDockerImage:                         "Build docker image",
	CreateAwsVpc:                             "Create AWS VPC",
	CreateEcsCluster:                         "Create ECS Cluster",
	UploadImageToEcr:                         "Upload image to ECR",
	AddAwsStaticSiteResponseHeaders:          "Add AWS Static Site Response Headers",
	UpdateAwsStaticSiteDomains:               "Update AWS Static Site Domains",
	DeployAwsCloudfrontViewerRequestFunction: "CF function on viewer request event",
	CreateAcmCertificate:                     "Create ACM certificate",
	UpdateAwsWebServiceDomain:                "Update AWS Web Service Domain",
	VerifyAcmCertificate:                     "Verify ACM certificate",
	DeleteAwsStaticSite:                      "Delete AWS Static Site",
	DeleteAwsWebService:                      "Delete AWS Web Service",
	ListCloudWatchMetricsAwsEcsWebService:    "List cloudwatch metrics for AWS ECS Web Service",
	CreateSecretAwsSecretManager:             "Create secret in AWS Secret Manager",
	DeployAwsPrivateService:                  "Deploy AWS Private Service",
	DeleteAwsPrivateService:                  "Delete AWS Private Service",
	DeployAwsRdsDatabase:                     "Deploy AWS RDS Database",
	DeleteAwsRdsDatabase:                     "Delete AWS RDS Database",
	BuildNixPacksImage:                       "Build NixPacks Image",
	RunNewAutomation:                         "Run New Automation",
	RunNewAgent:                              "Run New Agent",
	GetDeploymentLogsAws:                     "Get Deployment Logs from AWS CloudWatch",
	RunAgentStep:                             "Run Agent Step",
	CommitAndPush:                            "Commit and Push",
	OpenPullRequest:                          "Open Pull Request",
	RunAssistantSession:                      "Run Assistant Session",
	BuildInfraContext:                        "Build Infra Context",
	MaterializeContext:                       "Materialize Context",
	RunReviewStage:                           "Run Review Stage",
	RunQAStage:                               "Run QA Stage",
	RunProductionQA:                          "Run Production QA",
	DeleteAwsNamespaces:                      "Delete AWS Namespaces",
}

func (t Type) String() string {
	return typeToString[t]
}
