package llm_enums

type ApiVersion uint

const (
	FirstJanuary2025Preview ApiVersion = iota + 1

	ApiVersionMax //always add new versions before ApiVersionMax
)

var apiVersionToString = map[ApiVersion]string{
	FirstJanuary2025Preview: "2025-01-01-preview",
}

//

func (a ApiVersion) String() string {
	if a <= 0 || a >= ApiVersionMax {
		return ""
	}
	return apiVersionToString[a]
}
