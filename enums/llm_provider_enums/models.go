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
	// NovaProV1 is Amazon's own model and the first entry here that is
	// Bedrock-ONLY — there is no direct API for it. It is also the first
	// non-Anthropic, non-OpenAI model, which is what the vendor axis was
	// added for.
	NovaProV1

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
	NovaProV1:      "nova-pro-v1",
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
// modelToBedrockProfilePrefix maps a logical model to the VERSION-PINNED
// prefix shared by every concrete inference profile for that exact model.
//
// ⚠️ IT MUST PIN THE VERSION. An earlier revision used loose family tokens —
// "claude-sonnet-4" for claude-sonnet-4-6 — and matching is Contains(), so it
// also matched `eu.anthropic.claude-sonnet-4-5-20250929-v1:0`. A task asking
// for Sonnet 4.6 would have silently run Sonnet 4.5: no error, no log line,
// resolution looking entirely successful. Only the date and revision may be
// left to discovery, because those genuinely vary per region and account. The
// model VERSION is not something to guess at.
//
// Absent entry = not available on Bedrock (the GPT models). Kept in step with
// modelToProviders and enforced by catalog_test.go.
var modelToBedrockProfilePrefix = map[Model]string{
	ClaudeHaiku45:  "claude-haiku-4-5",
	ClaudeSonnet46: "claude-sonnet-4-6",
	ClaudeOpus48:   "claude-opus-4-8",
	// Matches ids like eu.amazon.nova-pro-v1:0 — the vendor and region
	// segments come from the profile id that discovery returns.
	NovaProV1: "nova-pro-v1",
}

// BedrockProfilePrefix returns the version-pinned prefix for a model, or ""
// when the model is not served by Bedrock. The caller matches this against the
// ids returned by bedrock:ListInferenceProfiles for its own region.
func (m Model) BedrockProfilePrefix() string {
	return modelToBedrockProfilePrefix[m]
}

// Vendor is who MAKES a model — a property of the model, never of the
// provider serving it.
//
// ⚠️ This deliberately does not hang off Provider. AWSBedrock serves models
// from Anthropic, Amazon, Google, Meta, Mistral, Qwen and others; GoogleVertex
// serves Claude alongside Gemini. A provider is a route to a model, not a
// statement about who built it, so asking "which vendor is AWSBedrock" has no
// answer. An earlier revision of this package got that wrong.
type Vendor uint

const (
	VendorUnknown Vendor = iota
	VendorAnthropic
	VendorOpenAI
	VendorGoogle
	VendorMeta
	VendorMistral
	VendorAmazon

	MaxVendor // always add vendors before MaxVendor
)

var vendorToString = map[Vendor]string{
	VendorAnthropic: "Anthropic",
	VendorOpenAI:    "OpenAI",
	VendorGoogle:    "Google",
	VendorMeta:      "Meta",
	VendorMistral:   "Mistral",
	VendorAmazon:    "Amazon",
}

func (v Vendor) String() string {
	return vendorToString[v]
}

// modelToVendor is only as long as the curated model list. Vendors beyond
// Anthropic and OpenAI are declared above because Bedrock and opencode make
// them reachable, not because a model here uses them yet.
var modelToVendor = map[Model]Vendor{
	ClaudeHaiku45:  VendorAnthropic,
	ClaudeSonnet46: VendorAnthropic,
	ClaudeOpus48:   VendorAnthropic,
	Gpt55:          VendorOpenAI,
	Gpt53Codex:     VendorOpenAI,
	Gpt54:          VendorOpenAI,
	NovaProV1:      VendorAmazon,
}

// Vendor returns who makes this model, independent of how it is reached.
// Use this for cost attribution or grouping — never Provider.
func (m Model) Vendor() Vendor {
	return modelToVendor[m]
}
