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
	RunNewAutomation                         Type = 24
	RunNewAgent                              Type = 25
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
}

func (t Type) String() string {
	return typeToString[t]
}
