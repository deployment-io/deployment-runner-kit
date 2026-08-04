package llm_provider_enums

import "fmt"

// Model is a LOGICAL model id — the model a user picks, independent of which
// provider serves it or which harness runs it.
//
// ⚠️ LOGICAL, NEVER CONCRETE. A Model names "Claude Opus 4.8"; it does not
// name `eu.anthropic.claude-opus-4-...-v1:0`. Provider-specific ids carry
// versions, revisions and region prefixes that change as models ship and
// differ per account — baking one in here would require a release of this
// module plus a bump in four repos every time AWS publishes a revision, and
// would fail at task time in between. The volatile half is resolved late:
// see BedrockFamily below and the runner's ListInferenceProfiles discovery.
//
// Not persisted today — Task documents store the string id (String() below),
// not the enum. Keep the strings stable for that reason; they are the wire
// format.
type Model uint

const (
	ClaudeHaiku45 Model = iota + 1
	ClaudeSonnet46
	ClaudeOpus48
	Gpt55
	Gpt53Codex
	Gpt54

	MaxModel // always add models before MaxModel
)

// modelToString is the WIRE FORMAT — what Task.Model holds, what the MCP
// create_task tool accepts, and what the dashboard sends. Changing one of
// these strings orphans every Task already storing it.
var modelToString = map[Model]string{
	ClaudeHaiku45:  "claude-haiku-4-5",
	ClaudeSonnet46: "claude-sonnet-4-6",
	ClaudeOpus48:   "claude-opus-4-8",
	Gpt55:          "gpt-5.5",
	Gpt53Codex:     "gpt-5.3-codex",
	Gpt54:          "gpt-5.4",
}

var stringToModel = func() map[string]Model {
	m := make(map[string]Model, len(modelToString))
	for k, v := range modelToString {
		m[v] = k
	}
	return m
}()

func (m Model) String() string {
	return modelToString[m]
}

func (m Model) IsValid() bool {
	return m > 0 && m < MaxModel
}

// GetModel parses a wire-format model id. Unknown ids are an error rather
// than a zero value so callers must decide explicitly what to do with a model
// this build does not know about — historically the lenient path let unknown
// ids through to the agent, which is still valid, but it should be a choice.
func GetModel(s string) (Model, error) {
	if m, ok := stringToModel[s]; ok {
		return m, nil
	}
	return 0, fmt.Errorf("unknown model id %q", s)
}

// modelToBedrockFamily maps a logical model to the Bedrock "family token"
// shared by every concrete inference profile for it.
//
// It deliberately stops at the family. A live run turned up two id shapes in
// one account — `anthropic.claude-sonnet-5` and
// `anthropic.claude-sonnet-4-5-20250929-v1:0` — so the version and revision
// come from discovery, not from here. This map only needs to be right about
// which family a model belongs to, which changes rarely.
//
// Absent entry = not available on Bedrock (the GPT models). Keep in step with
// modelToProviders: a model listing AWSBedrock as a provider needs a family
// here, and catalog_test.go enforces exactly that.
var modelToBedrockFamily = map[Model]string{
	ClaudeHaiku45:  "claude-haiku-4",
	ClaudeSonnet46: "claude-sonnet-4",
	ClaudeOpus48:   "claude-opus-4",
}

// BedrockFamily returns the Bedrock family token for a model, or "" when the
// model is not served by Bedrock. The caller matches this against the ids
// returned by bedrock:ListInferenceProfiles for its own region.
func (m Model) BedrockFamily() string {
	return modelToBedrockFamily[m]
}
