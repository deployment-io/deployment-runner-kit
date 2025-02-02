package llm_enums

type ModelType uint

const (
	Gpt4O ModelType = iota + 1

	ModelTypeMax //always add new types before ModelTypeMax
)

var modelTypeToString = map[ModelType]string{
	Gpt4O: "gpt-4o",
}

func (m ModelType) String() string {
	return modelTypeToString[m]
}
