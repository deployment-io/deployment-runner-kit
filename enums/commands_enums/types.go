package commands_enums

/*
CheckoutRepo
*/

type Type uint

const (
	CheckoutRepo Type = iota + 1
	BuildStaticSite
	DeployStaticSiteAWS
)

var typeToString = map[Type]string{
	CheckoutRepo:        "Checkout Repo",
	BuildStaticSite:     "Build Static Site",
	DeployStaticSiteAWS: "Deploy Static Site AWS",
}

func (t Type) String() string {
	return typeToString[t]
}
