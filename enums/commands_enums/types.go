package commands_enums

/*
CheckoutRepo
*/

type Type uint

const (
	CheckoutRepo        Type = 1
	BuildStaticSite     Type = 2
	DeployStaticSiteAWS Type = 3
	DeployWebServiceAWS Type = 4
	BuildDockerImage    Type = 5
)

var typeToString = map[Type]string{
	CheckoutRepo:        "Checkout Repo",
	BuildStaticSite:     "Build Static Site",
	DeployStaticSiteAWS: "Deploy Static Site AWS",
	DeployWebServiceAWS: "Deploy Web Service AWS",
	BuildDockerImage:    "Build docker image",
}

func (t Type) String() string {
	return typeToString[t]
}
