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
	// The 4.5-generation Claude models. Kept in the catalogue because BEDROCK
	// AVAILABILITY LAGS the direct API and varies by account and region: an org
	// whose Bedrock has Sonnet 4.5 but not 4.6 could otherwise run no Claude
	// model there at all, since the profile prefixes are version-pinned and 4.6
	// deliberately will not match a 4.5 profile.
	//
	// Ordinary catalogue models, not a Bedrock special case — the direct API
	// still serves them, and pinning an older version is a legitimate choice.
	ClaudeSonnet45
	ClaudeOpus45

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
	ClaudeSonnet45: "claude-sonnet-4-5",
	ClaudeOpus45:   "claude-opus-4-5",
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
	NovaProV1:      "nova-pro-v1",
	ClaudeSonnet45: "claude-sonnet-4-5",
	ClaudeOpus45:   "claude-opus-4-5",
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
	ClaudeSonnet45: VendorAnthropic,
	ClaudeOpus45:   VendorAnthropic,
}

// Vendor returns who makes this model, independent of how it is reached.
// Use this for cost attribution or grouping — never Provider.
func (m Model) Vendor() Vendor {
	return modelToVendor[m]
}

// opencodeBedrockGeographyVendors names the vendors whose Bedrock models opencode
// lists UNDER their cross-region geography prefix.
//
// This mirrors OPENCODE's registry (models.dev), not Bedrock's. Bedrock
// publishes a geography-prefixed profile for nearly everything, but models.dev
// only carries one where that profile is the invocable id, and it is starkly
// split — counted from the registry:
//
//	anthropic  11 base + 45 geography-prefixed   overwhelmingly prefixed
//	amazon      4 base +  0                      never prefixed
//	meta        5 base +  2                      MIXED — see below
//	openai, mistral, qwen, nvidia, google, zai, minimax, writer, xai — base only
//
// Amazon is NOT the special case it first appears to be: ten of the fifteen
// vendors are base-only, so KEEPING the prefix is the exception and this map
// is the exception list.
//
// Meta is deliberately ABSENT despite having two prefixed entries. A
// vendor-wide "keep" would be wrong for the other five, so Meta's split is
// per-MODEL and this map cannot express it. Nothing forces the question yet —
// no Meta model is in the catalogue — and inventing an answer here would bake
// in a guess that no test could exercise. When a Meta model is added,
// TestBedrockModelsHaveAConsideredGeographyRule fails and the decision gets
// made against the registry as it stands then, which may need a per-model
// override rather than an entry here.
//
// An absent vendor means "strip", which is the right default for a new vendor
// and, when wrong, fails loudly — opencode answers ProviderModelNotFoundError
// rather than quietly running something else.
var opencodeBedrockGeographyVendors = map[Vendor]bool{
	VendorAnthropic: true,
}

// OpencodeKeepsBedrockGeography reports whether OPENCODE expects this vendor's
// Bedrock ids to keep their cross-region prefix.
//
// ⚠️ OPENCODE-SPECIFIC, and named so on purpose. This is not a statement about
// what Bedrock accepts. Bedrock publishes eu.amazon.nova-pro-v1:0 and will
// serve it — discovery found that id there. Only opencode's registry lacks an
// entry for it, and only opencode routes through that registry.
//
// So this must NEVER be applied to claude-code, which reaches Bedrock directly
// and correctly uses the prefixed profile id for EVERY vendor, Amazon included.
// A neutral-sounding name here would invite exactly that mistake.
//
// VendorUnknown answers true, so a model outside the catalogue passes through
// untouched rather than being rewritten on a guess.
func (v Vendor) OpencodeKeepsBedrockGeography() bool {
	return v == VendorUnknown || opencodeBedrockGeographyVendors[v]
}

// modelProviderID overrides the id used at a specific provider, for models
// whose provider-side name differs from our logical one.
//
// DELIBERATELY EMPTY. No divergence is currently KNOWN — but that is
// unverified rather than established, which is the reason this seam exists at
// all. Our ids are ours: opencode resolves models through its own registry
// (models.dev), and that registry need not agree with either our names or the
// vendor API's. Populate an entry the moment a real run shows a mismatch;
// guessing now would bake in a second unverified assumption.
//
// Not for providers that resolve ids at runtime — see
// Provider.ResolvesModelIDAtRuntime.
var modelProviderID = map[Model]map[Provider]string{}

// IDFor returns the model id to use at a given provider: an override when the
// provider names the model differently, otherwise our logical id.
//
// Only meaningful for providers that DECLARE their ids. Where
// Provider.ResolvesModelIDAtRuntime is true the concrete id is discovered and
// this value is at best a starting point.
func (m Model) IDFor(p Provider) string {
	if byProvider, ok := modelProviderID[m]; ok {
		if id, ok := byProvider[p]; ok {
			return id
		}
	}
	return m.String()
}

// modelToDisplayName is what a human calls the model, without its vendor —
// "Sonnet 4.6", not "claude-sonnet-4-6" and not "Anthropic Sonnet 4.6".
//
// Separate from String(), which is the WIRE id and must never move: a copy edit
// to a label would otherwise change what gets sent to a provider. Vendor is
// separate too, so a caller composes "Anthropic · Sonnet 4.6" or groups by
// vendor without this map having to pick one presentation.
//
// Lived in the dashboard until now, which meant a model added here was invisible
// there until someone remembered to mirror it — and mirrored ids drifted from
// the catalogue three times.
var modelToDisplayName = map[Model]string{
	ClaudeHaiku45:  "Haiku 4.5",
	ClaudeSonnet46: "Sonnet 4.6",
	ClaudeOpus48:   "Opus 4.8",
	Gpt55:          "GPT-5.5",
	Gpt53Codex:     "GPT-5.3 Codex",
	Gpt54:          "GPT-5.4",
	NovaProV1:      "Nova Pro",
	ClaudeSonnet45: "Sonnet 4.5",
	ClaudeOpus45:   "Opus 4.5",
}

// DisplayName returns the human-facing name, falling back to the wire id so an
// unnamed model still renders as something rather than blank.
func (m Model) DisplayName() string {
	if name, ok := modelToDisplayName[m]; ok {
		return name
	}
	return m.String()
}

// legacyModels are superseded by a newer generation but kept offerable.
//
// They exist because BEDROCK AVAILABILITY LAGS the direct API: an account with
// Sonnet 4.5 and not 4.6 needs the older entry to run anything at all, since
// the profile prefixes are version-pinned and refuse to cross generations.
//
// A DISPLAY concern only. It must not touch agentTypeToModels' order, which is
// capability-ascending and indexed by kit's recommendedByComplexity — sorting
// legacy to the end of that list would make "high complexity" resolve to Opus
// 4.5. ModelsFor applies the ordering instead, where it affects pickers alone.
var legacyModels = map[Model]bool{
	ClaudeSonnet45: true,
	ClaudeOpus45:   true,
}

// IsLegacy reports whether a newer generation of this model exists.
func (m Model) IsLegacy() bool { return legacyModels[m] }

// Tier is how capable a model is, independent of when it shipped.
//
// A PROPERTY OF THE MODEL, not its position in a list. Recommendation used to
// index into agentTypeToModels — first entry for a simple task, last for a
// complex one — which quietly required that list to stay capability-ascending
// forever. Two things break that: adding an older generation (Sonnet 4.5 sits
// between Haiku and Sonnet 4.6 by capability, not by release date), and the
// passage of time. Today's frontier model is next year's balanced one, and
// re-tagging it should be an edit HERE rather than a silent reshuffle of a
// list that several other things read.
type Tier uint

const (
	TierUnknown Tier = iota
	// TierFast — cheapest and quickest; enough for mechanical work.
	TierFast
	// TierBalanced — the default choice for most tasks.
	TierBalanced
	// TierFrontier — the most capable available, for genuinely hard work.
	TierFrontier

	MaxTier
)

var tierToString = map[Tier]string{
	TierFast:     "fast",
	TierBalanced: "balanced",
	TierFrontier: "frontier",
}

func (t Tier) String() string { return tierToString[t] }

// modelToTier states each model's capability band.
//
// EXPECTED TO CHANGE as newer generations arrive — that is the point. When a
// frontier model is superseded it moves down a band here, and every caller
// follows without any list being reordered.
var modelToTier = map[Model]Tier{
	ClaudeHaiku45:  TierFast,
	ClaudeSonnet45: TierBalanced,
	ClaudeSonnet46: TierBalanced,
	ClaudeOpus45:   TierFrontier,
	ClaudeOpus48:   TierFrontier,
	// The codex lineup is not cleanly tiered — 5.3-codex is task-specialised
	// rather than weaker — so all three sit at balanced and the recommendation
	// falls back to the agent's default.
	Gpt55:      TierBalanced,
	Gpt53Codex: TierBalanced,
	Gpt54:      TierBalanced,
	NovaProV1:  TierBalanced,
}

// Tier returns the model's capability band.
func (m Model) Tier() Tier { return modelToTier[m] }

// ModelForTier returns the model an agent should use at this tier, preferring
// a current generation over a superseded one.
//
// Falls back to the agent's default when nothing matches, so a caller always
// gets a usable model rather than an empty string — a tier with no model is a
// catalogue gap, not something the caller should have to handle.
func ModelForTier(h AgentType, t Tier) Model {
	var current, legacy Model
	for _, m := range agentTypeToModels[h] {
		if m.Tier() != t {
			continue
		}
		// The agent's own default wins its tier outright — an explicit choice
		// beats any incidental one. This is what keeps a lineup that is not
		// cleanly tiered (codex, where all three sit at balanced) from
		// resolving to whichever model happened to be listed first.
		if m == h.DefaultModel() {
			return m
		}
		if m.IsLegacy() {
			if legacy == 0 {
				legacy = m
			}
			continue
		}
		if current == 0 {
			current = m
		}
	}
	switch {
	case current != 0:
		return current
	case legacy != 0:
		return legacy
	}
	return h.DefaultModel()
}
