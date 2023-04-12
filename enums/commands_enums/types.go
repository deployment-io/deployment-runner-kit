package commands_enums

/*
CheckoutRepo
*/

type Type uint

const (
	CheckoutRepo Type = iota + 1
)

var typeToString = map[Type]string{
	CheckoutRepo: "Checkout Repo",
}

func (t Type) String() string {
	return typeToString[t]
}
