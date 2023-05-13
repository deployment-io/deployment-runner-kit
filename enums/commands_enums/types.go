package commands_enums

/*
CheckoutRepo
*/

type Type uint

const (
	CheckoutRepo Type = iota + 1
	BuildStaticSite
	UploadStaticSite
)

var typeToString = map[Type]string{
	CheckoutRepo:     "Checkout Repo",
	BuildStaticSite:  "Build Static Site",
	UploadStaticSite: "Upload Static Site",
}

func (t Type) String() string {
	return typeToString[t]
}
