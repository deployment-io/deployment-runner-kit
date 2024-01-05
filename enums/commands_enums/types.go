package commands_enums

/*
CheckoutRepo
*/

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
}

func (t Type) String() string {
	return typeToString[t]
}
