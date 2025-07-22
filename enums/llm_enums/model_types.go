package llm_enums

type ModelType uint

const (
	Gpt4O ModelType = iota + 1
	Gpt41

	ModelTypeMax //always add new types before ModelTypeMax
)

var modelTypeToString = map[ModelType]string{
	Gpt4O: "gpt-4o",
	Gpt41: "gpt-4.1",
}

func (m ModelType) String() string {
	if m <= 0 || m >= ModelTypeMax {
		return ""
	}
	return modelTypeToString[m]
}
