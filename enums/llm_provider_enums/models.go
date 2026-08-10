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

	// The current Claude generation. Sonnet 5 is the everyday coding model and
	// Opus 5 the frontier one; Fable 5 is a SPECIALIST (creative writing and
	// analysis) that costs twice Opus 5, so it is offered but never resolved
	// into by tier — Opus 4.8 and Opus 5 both precede it in agentTypeToModels.
	ClaudeSonnet5
	ClaudeOpus5
	ClaudeFable5

	// Bedrock-only models from vendors with no direct provider here. These are
	// what makes opencode worth having: BYO-model at a fraction of the frontier
	// price, on credentials the org already holds.
	//
	// ⚠️ ALL SEVEN ARE BASE-ID-ONLY ON BEDROCK — no cross-region inference
	// profile exists for any of them, so discovery finds nothing and their
	// concrete ids are DECLARED in modelProviderID rather than resolved. That is
	// the opposite of every model above, where the id is deliberately late-bound.
	Qwen3Coder480B
	Qwen3CoderNext
	DeepSeekV32
	Glm47
	Glm5
	MinimaxM25
	Grok43

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
	ClaudeSonnet5:  "claude-sonnet-5",
	ClaudeOpus5:    "claude-opus-5",
	ClaudeFable5:   "claude-fable-5",
	// Logical ids, in the same style as every other entry — NOT the Bedrock
	// ids. "qwen3-coder-480b", not "qwen.qwen3-coder-480b-a35b-v1:0": the
	// vendor segment and parameter-count suffix are Bedrock's naming, and
	// baking them into the wire format would tie a stored Task to one provider.
	Qwen3Coder480B: "qwen3-coder-480b",
	Qwen3CoderNext: "qwen3-coder-next",
	DeepSeekV32:    "deepseek-v3.2",
	Glm47:          "glm-4.7",
	Glm5:           "glm-5",
	MinimaxM25:     "minimax-m2.5",
	Grok43:         "grok-4.3",
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
// An absent entry means one of TWO different things, which is why
// TestCatalog_BedrockModelsDeclareExactlyOneIDMechanism exists:
//
//	not served by Bedrock at all           — the GPT models
//	served, but with NO inference profile  — declared in modelProviderID
//
// The second case is the whole Bedrock-only lineup (Qwen, DeepSeek, GLM,
// MiniMax, Grok). Nothing to discover means nothing to pin, so their ids are
// stated outright instead.
var modelToBedrockProfilePrefix = map[Model]string{
	ClaudeHaiku45:  "claude-haiku-4-5",
	ClaudeSonnet46: "claude-sonnet-4-6",
	ClaudeOpus48:   "claude-opus-4-8",
	// Matches ids like eu.amazon.nova-pro-v1:0 — the vendor and region
	// segments come from the profile id that discovery returns.
	NovaProV1:      "nova-pro-v1",
	ClaudeSonnet45: "claude-sonnet-4-5",
	ClaudeOpus45:   "claude-opus-4-5",
	// The 5-generation prefixes cannot collide with their 4.x siblings under
	// Contains(): "claude-opus-5" is not a substring of
	// "claude-opus-4-5-20251101-v1:0", and vice versa. That is the property the
	// version pinning exists to hold, and the reason these read oddly short.
	ClaudeSonnet5: "claude-sonnet-5",
	ClaudeOpus5:   "claude-opus-5",
	ClaudeFable5:  "claude-fable-5",
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
	// Vendors reachable only through Bedrock here — no direct provider exists
	// for any of them, which is the point: opencode can serve them on the AWS
	// credential an org already has.
	VendorQwen
	VendorDeepSeek
	VendorZAI
	VendorMiniMax
	VendorXAI

	MaxVendor // always add vendors before MaxVendor
)

var vendorToString = map[Vendor]string{
	VendorAnthropic: "Anthropic",
	VendorOpenAI:    "OpenAI",
	VendorGoogle:    "Google",
	VendorMeta:      "Meta",
	VendorMistral:   "Mistral",
	VendorAmazon:    "Amazon",
	VendorQwen:      "Qwen",
	VendorDeepSeek:  "DeepSeek",
	VendorZAI:       "Z.ai",
	VendorMiniMax:   "MiniMax",
	VendorXAI:       "xAI",
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
	ClaudeSonnet5:  VendorAnthropic,
	ClaudeOpus5:    VendorAnthropic,
	ClaudeFable5:   VendorAnthropic,
	Qwen3Coder480B: VendorQwen,
	Qwen3CoderNext: VendorQwen,
	DeepSeekV32:    VendorDeepSeek,
	Glm47:          VendorZAI,
	Glm5:           VendorZAI,
	MinimaxM25:     VendorMiniMax,
	Grok43:         VendorXAI,
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
// The Bedrock entries here are NOT the same kind of thing as the Anthropic
// models above them. Those are late-bound on purpose: their concrete ids carry
// regions, dates and revisions that vary per account, so discovery resolves
// them and pinning one here would break between AWS revisions. These seven have
// nothing to discover — no cross-region inference profile exists for any of
// them, so ListInferenceProfiles returns nothing to match and the id has to be
// stated.
//
// Taken verbatim from models.dev's amazon-bedrock registry, which is what
// OPENCODE resolves against — so these must track that registry, not Bedrock's
// own naming, if the two ever disagree.
//
// Nothing else populates this map. A provider whose ids are discovered should
// not appear here — see Provider.ResolvesModelIDAtRuntime.
var modelProviderID = map[Model]map[Provider]string{
	Qwen3Coder480B: {AWSBedrock: "qwen.qwen3-coder-480b-a35b-v1:0"},
	Qwen3CoderNext: {AWSBedrock: "qwen.qwen3-coder-next"},
	DeepSeekV32:    {AWSBedrock: "deepseek.v3.2"},
	Glm47:          {AWSBedrock: "zai.glm-4.7"},
	Glm5:           {AWSBedrock: "zai.glm-5"},
	MinimaxM25:     {AWSBedrock: "minimax.minimax-m2.5"},
	Grok43:         {AWSBedrock: "xai.grok-4.3"},
}

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
	ClaudeSonnet5:  "Sonnet 5",
	ClaudeOpus5:    "Opus 5",
	ClaudeFable5:   "Fable 5",
	// The vendor is rendered separately, so these carry none — a picker shows
	// "Qwen · Qwen3 Coder 480B" and the second word would otherwise repeat.
	// The parameter counts stay because they distinguish real variants.
	Qwen3Coder480B: "Qwen3 Coder 480B",
	Qwen3CoderNext: "Qwen3 Coder Next",
	DeepSeekV32:    "V3.2",
	Glm47:          "GLM-4.7",
	Glm5:           "GLM-5",
	MinimaxM25:     "M2.5",
	Grok43:         "Grok 4.3",
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
	// GLM-5 supersedes it. Kept for the same reason as the 4.5-generation
	// Claude models: Bedrock model access is granted per account, and an org
	// with 4.7 enabled and not 5 should still have something to run.
	Glm47: true,
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

	ClaudeSonnet5: TierBalanced,
	ClaudeOpus5:   TierFrontier,
	// Frontier by price and capability, but a SPECIALIST — creative writing and
	// analysis rather than coding. It never wins a tier resolution because Opus
	// 4.8 and Opus 5 both precede it in agentTypeToModels, so reaching it takes
	// an explicit pick. Same treatment as 5.3-codex: listed honestly, not
	// promoted.
	ClaudeFable5: TierFrontier,

	// The Bedrock-only lineup sits at balanced, deliberately. These are strong
	// coding models at roughly a tenth of frontier cost, but ranking them
	// against Claude would be a claim this catalogue cannot support — none has
	// been run here yet. Balanced states "a reasonable default choice", which is
	// defensible; frontier would state "the most capable available", which is
	// not. Revisit per model once real tasks have run on them.
	Qwen3Coder480B: TierBalanced,
	Qwen3CoderNext: TierBalanced,
	DeepSeekV32:    TierBalanced,
	Glm47:          TierBalanced,
	Glm5:           TierBalanced,
	MinimaxM25:     TierBalanced,
	Grok43:         TierBalanced,
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
