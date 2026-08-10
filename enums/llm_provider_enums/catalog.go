package llm_provider_enums

import (
	"sort"
	"strings"
)

// The catalogue: which (harness, model, provider) combinations are real.
//
// Held as THREE FLAT TABLES rather than one nested map[AgentType]map[Model][]Provider,
// because each table states an independently true fact and a valid combination
// is their intersection:
//
//	agentTypeToModels    — claude-code cannot run gpt-5.5 whoever serves it
//	modelToProviders   — Opus 4.8 is served by Anthropic and Bedrock, not OpenAI
//	agentTypeToProviders — codex only ever talks to OpenAI
//
// A nested map would restate the same facts once per combination and drift
// between rows. Worked example, and the question this package was built to
// answer: "claude-code supports Opus 4.8, which could come from Anthropic
// Direct or AWS Bedrock" —
//
//	ProvidersFor(ClaudeCode, ClaudeOpus48)
//	  = modelToProviders[ClaudeOpus48]   // Direct, Subscription, Bedrock
//	  ∩ agentTypeToProviders[ClaudeCode]   // Direct, Subscription, Bedrock, Vertex
//	  = Direct, Subscription, Bedrock
//
// This is a STATIC LOOKUP TABLE, not state. Nothing here is persisted; it
// says what is *possible*, never what a given org has configured. Whether an
// org can actually use a provider is a separate question answered by its
// credentials — and for Bedrock, by per-runner readiness.

// agentTypeToModels lists the models each harness can run.
//
// opencode's entries are the same logical models as the others: the
// "anthropic/" and "openai/" prefixes its CLI wants are a rendering of
// (model, provider) at spawn, not part of the model's identity. That is the
// point of keeping Model logical — see PLAN_provider_centric_llm_keys.md §3.3
// and §4.1, where the prefix's second job (disambiguating the harness) is what
// made stripping it dangerous until AgentType became explicit.
var agentTypeToModels = map[AgentType][]Model{
	ClaudeCode: {ClaudeHaiku45, ClaudeSonnet45, ClaudeSonnet46, ClaudeOpus45, ClaudeOpus48},
	Codex:      {Gpt55, Gpt53Codex, Gpt54},
	Opencode:   {ClaudeHaiku45, ClaudeSonnet45, ClaudeSonnet46, ClaudeOpus45, ClaudeOpus48, Gpt55, NovaProV1},
}

// modelToProviders lists which providers can serve each model.
//
// AnthropicSubscription appears on the Claude models because the subscription
// genuinely serves them — the restriction that only claude-code may use it is
// a HARNESS constraint, expressed in agentTypeToProviders, not a model one.
// Keeping those separate is what lets the intersection stay correct.
var modelToProviders = map[Model][]Provider{
	ClaudeHaiku45:  {AnthropicDirect, AnthropicSubscription, AWSBedrock},
	ClaudeSonnet46: {AnthropicDirect, AnthropicSubscription, AWSBedrock},
	ClaudeOpus48:   {AnthropicDirect, AnthropicSubscription, AWSBedrock},
	Gpt55:          {OpenAIDirect},
	Gpt53Codex:     {OpenAIDirect},
	Gpt54:          {OpenAIDirect},
	// Bedrock-only: Amazon does not offer Nova through a direct API, so this
	// is the first model whose single provider is a cloud route rather than
	// its vendor. claude-code and codex cannot run it — agentTypeToModels keeps
	// it to opencode.
	NovaProV1: {AWSBedrock},
	// Same providers as their 4.6/4.8 siblings — the direct API and a
	// subscription still serve them, and Bedrock is where they matter most.
	ClaudeSonnet45: {AnthropicDirect, AnthropicSubscription, AWSBedrock},
	ClaudeOpus45:   {AnthropicDirect, AnthropicSubscription, AWSBedrock},
}

// agentProviderCapabilities lists which provider APIs each CLI can talk to.
//
// AN INDEPENDENT AXIS, not derivable from the models. It is tempting to compute
// it — an agent could just inherit whatever its models reach — and today that
// would even produce the right answer, because no agent's model list overlaps a
// provider it cannot use. That is a COINCIDENCE, and deriving would encode it
// as a rule: adding one Bedrock-served model to codex would silently grant
// codex Bedrock, which it has no path to at all.
//
// The axis is TRANSPORT: which endpoint the CLI knows how to authenticate to
// and speak. Model availability is the other axis, and ProvidersFor intersects
// them. Keeping them apart matters because they genuinely disagree:
//
//   - Codex is OpenAI-only, and NOT because Bedrock lacks OpenAI models —
//     Bedrock does host some (the open-weight gpt-oss family). Our driver runs
//     `codex login --with-api-key` and the CLI authenticates against
//     wss://api.openai.com/v1/responses. Bedrock is not an OpenAI-compatible
//     endpoint, so no model on it is reachable, whoever publishes the model.
//   - opencode reaches Bedrock through its own "amazon-bedrock/…" model id, and
//     claude-code through CLAUDE_CODE_USE_BEDROCK. Same provider, three
//     different answers, none of them a property of the model.
//   - opencode is absent from AnthropicSubscription for a POLICY reason, not a
//     transport one: it could speak the API perfectly well, but a Pro/Max token
//     is validated against Anthropic's client-identity check and routing one
//     through a third-party agent is prohibited. Sonnet IS served by a
//     subscription and opencode DOES run Sonnet, so the model axis would allow
//     it. TestCatalog_SubscriptionIsClaudeCodeOnly is the guard.
//
// That last line is the disproof of "providers depend only on models": opencode
// and codex can support the SAME model on the SAME provider and still differ.
var agentProviderCapabilities = map[AgentType][]Provider{
	ClaudeCode: {AnthropicDirect, AnthropicSubscription, AWSBedrock},
	Codex:      {OpenAIDirect},
	Opencode:   {AnthropicDirect, AWSBedrock, OpenAIDirect},
}

// Models returns the models this harness can run.
func (h AgentType) Models() []Model { return agentTypeToModels[h] }

// Providers returns the provider APIs this harness can talk to. Order carries
// NO preference — see PreferredProvider.
func (h AgentType) Providers() []Provider { return agentProviderCapabilities[h] }

// Providers returns every provider that can serve this model, ignoring which
// harness is asking. Use ProvidersFor when a harness is known.
func (m Model) Providers() []Provider { return modelToProviders[m] }

// Supports reports whether the harness can run the model at all.
func (h AgentType) Supports(m Model) bool {
	for _, candidate := range agentTypeToModels[h] {
		if candidate == m {
			return true
		}
	}
	return false
}

// ProvidersFor returns the providers that can serve this model through this
// harness — the intersection of what serves the model and what the harness can
// reach. Returns nil when the harness cannot run the model at all, so an empty
// result always means "not possible", never "possible but unconfigured".
//
// Order follows modelToProviders so the result is deterministic; callers that
// need a preferred provider should apply their own routing rather than relying
// on position.
func ProvidersFor(h AgentType, m Model) []Provider {
	if !h.Supports(m) {
		return nil
	}
	agentProviders := h.Providers()
	reachable := make(map[Provider]bool, len(agentProviders))
	for _, p := range agentProviders {
		reachable[p] = true
	}
	var out []Provider
	for _, p := range modelToProviders[m] {
		if reachable[p] {
			out = append(out, p)
		}
	}
	return out
}

// AgentTypesFor returns every harness able to run the model. This is the
// inverse of agentTypeToModels, and is what a model-first UI needs in order to
// offer a harness choice once a model is picked — with the first entry as the
// sensible default.
func AgentTypesFor(m Model) []AgentType {
	var out []AgentType
	for _, h := range AllAgentTypes() {
		if h.Supports(m) {
			out = append(out, h)
		}
	}
	// Already in priority order — AllAgentTypes ranks, and this only filters.
	// Callers take [0] as the agent that will run the model, so that ranking
	// decides which harness a shared id resolves to.
	return out
}

// opencodeProviderName maps a provider to the id opencode uses in its
// "provider/model" strings.
//
// Empty means opencode cannot reach the provider that way:
//
//   - AnthropicSubscription — subscription auth is harness-locked to genuine
//     claude-code, so opencode must never be handed one. agentTypeToProviders
//     already excludes it; this is the second line of defence.
//   - GoogleVertex — declared but not wired anywhere.
var opencodeProviderName = map[Provider]string{
	AnthropicDirect: "anthropic",
	OpenAIDirect:    "openai",
	AWSBedrock:      "amazon-bedrock",
}

// OpencodeName returns the provider id opencode expects, or "" when opencode
// cannot use this provider.
func (p Provider) OpencodeName() string {
	return opencodeProviderName[p]
}

// OpencodeModelID renders the "provider/model" string opencode takes as its
// --model argument.
//
// modelID is passed as a STRING rather than a Model because the caller may
// have resolved it to something more concrete than the catalogue holds — on
// Bedrock the runner substitutes the discovered inference-profile id
// (eu.amazon.nova-pro-v1:0) for the logical one. The provider prefix is what
// opencode routes on, so it must survive that substitution.
//
// Returns "" when the provider has no opencode representation, so callers get
// an empty result rather than a malformed "/model" string.
func OpencodeModelID(modelID string, p Provider) string {
	name := p.OpencodeName()
	if name == "" || modelID == "" {
		return ""
	}
	// opencode routes Bedrock through its own registry (models.dev), which is
	// ASYMMETRIC about the cross-region geography prefix:
	//
	//   anthropic.*  base AND eu./us./au./jp./global. variants
	//   amazon.*     base ONLY — no geography-prefixed entries at all
	//
	// So the prefix is stripped for Amazon's own models and KEPT for everyone
	// else. Keeping it matters: for Anthropic the prefixed id IS the
	// cross-region inference profile, which newer Claude models on Bedrock
	// generally require — dropping to the base id would silently ask for
	// on-demand throughput they may not offer.
	//
	// A live Nova run exposed this. Discovery resolved eu.amazon.nova-pro-v1:0
	// and opencode rejected it with ProviderModelNotFoundError, suggesting
	// amazon.nova-pro-v1:0.
	//
	// Only the geography goes; the dated revision stays, since that is what
	// discovery exists to find.
	if p == AWSBedrock {
		modelID = stripBedrockGeographyForOpencode(modelID)
	}
	return name + "/" + modelID
}

// stripBedrockGeographyForOpencode removes the cross-region prefix from an
// Amazon-vendor Bedrock id, leaving every other vendor untouched.
//
// Mirrors OPENCODE's registry, not Bedrock's: Bedrock does publish
// eu.amazon.nova-pro-v1:0 — discovery finds it — but opencode has no entry for
// it and cannot route it.
func stripBedrockGeographyForOpencode(modelID string) string {
	for _, geo := range []string{"eu.", "us.", "au.", "jp.", "apac.", "global."} {
		rest, found := strings.CutPrefix(modelID, geo)
		if !found {
			continue
		}
		// Only Amazon's own models lack geography-prefixed registry entries.
		if strings.HasPrefix(rest, "amazon.") {
			return rest
		}
		return modelID
	}
	return modelID
}

// ConfigurableProviders returns every provider an org can actually configure,
// in enum order.
//
// DERIVED, not a second list: a provider is configurable exactly when some
// agent can use it. That makes the catalogue the only place a provider is
// switched on, so a reserved-but-unbuilt one (GoogleVertex) cannot be offered
// by the settings UI or accepted by the API without someone first adding it to
// a real agent — rather than each caller carrying its own reject list that
// drifts.
func ConfigurableProviders() []Provider {
	offered := map[Provider]bool{}
	for agentType := range agentTypeToModels {
		providers := agentType.Providers()
		for _, p := range providers {
			offered[p] = true
		}
	}
	var out []Provider
	for _, p := range AllProviders() {
		if offered[p] {
			out = append(out, p)
		}
	}
	return out
}

// IsConfigurable reports whether an org can configure this provider.
func (p Provider) IsConfigurable() bool {
	for _, c := range ConfigurableProviders() {
		if c == p {
			return true
		}
	}
	return false
}

// PreferredProvider returns which of these candidates should serve a Job, given
// what the org has configured.
//
// EXISTS SO PREFERENCE IS NOT LIST ORDER. agentTypeToProviders answers
// MEMBERSHIP — which providers can serve an agent — and a list literal's order
// is the kind of thing someone tidies alphabetically. Deciding billing as a
// side effect of that is indefensible, so the rule lives here, named, with its
// reason attached.
//
// The rule is: prefer a credential the org has ALREADY PAID FOR. A subscription
// is a flat fee whether or not we use it, so reaching for a metered API key
// instead charges the org twice — silently, since nothing fails. Everything
// else ties and falls back to catalogue order, which is arbitrary but
// deterministic; no defensible reason ranks an API key against a cloud role.
//
// INTERIM. When per-model routing lands, an org's own configured order replaces
// this and the guessing stops.
func PreferredProvider(candidates []Provider, isConfigured func(Provider) bool) (Provider, bool) {
	best, found := Provider(0), false
	for _, p := range candidates {
		if !isConfigured(p) {
			continue
		}
		if !found || preferenceRank(p) < preferenceRank(best) {
			best, found = p, true
		}
	}
	return best, found
}

// preferenceRank orders providers by whether the org is already paying for them
// regardless of use. Lower wins; ties keep the caller's order.
func preferenceRank(p Provider) int {
	if p.AuthMode() == AuthSubscription {
		return 0
	}
	return 1
}

// agentTypeToDefaultModel is the model a picker should preselect.
//
// A product choice, so it lives beside the catalogue rather than in whichever
// client renders the picker first — two clients would otherwise preselect
// differently for the same agent.
var agentTypeToDefaultModel = map[AgentType]Model{
	// Anthropic's documented pick for agentic coding.
	ClaudeCode: ClaudeOpus48,
	// OpenAI's recommended default for Codex.
	Codex: Gpt55,
	// Balanced, and reuses an org's existing Anthropic credential — so opencode
	// is usable without configuring a new provider.
	Opencode: ClaudeSonnet46,
}

// DefaultModel returns the model to preselect for this agent.
func (h AgentType) DefaultModel() Model { return agentTypeToDefaultModel[h] }

// ModelsFor returns the models this agent can run that the org can actually
// serve, in catalogue order.
//
// The intersection callers keep needing: a model is offerable only when the
// agent runs it AND some provider serving it is configured. Offering more than
// that means a picker shows a model whose Task fails at pickup — which is what
// a client-side model list cannot avoid, because it cannot see the org.
//
// isConfigured is passed in rather than the org itself: this package must not
// learn what an organization is.
func ModelsFor(h AgentType, isConfigured func(Provider) bool) []Model {
	var out []Model
	for _, m := range agentTypeToModels[h] {
		for _, p := range ProvidersFor(h, m) {
			if isConfigured(p) {
				out = append(out, m)
				break
			}
		}
	}
	// Legacy generations last, current order preserved within each group. Done
	// HERE rather than in agentTypeToModels because that list is
	// capability-ascending and kit's recommendedByComplexity indexes into it —
	// reordering it would make "high complexity" resolve to a superseded model.
	//
	// Sorted server-side so every client agrees. A picker that ordered these
	// itself would be one more copy of a catalogue fact.
	sort.SliceStable(out, func(i, j int) bool {
		return !out[i].IsLegacy() && out[j].IsLegacy()
	})
	return out
}

// AllAgentTypes returns every agent, in priority order — the same order
// AgentTypesFor uses to break a tie when several can run one model.
func AllAgentTypes() []AgentType {
	// Swept from the declaration map rather than counted through a range.
	// Counting assumes the values are contiguous, which stops being true the
	// moment one is reserved or retired — the same reason AllProviders is
	// derived.
	out := make([]AgentType, 0, len(agentTypeToString))
	for h := range agentTypeToString {
		out = append(out, h)
	}
	// PRIORITY order, not numeric. This is what app-server iterates to build
	// the agent picker, so sorting by value would put the enum's declaration
	// order on screen — the same dependency AgentTypesFor just stopped
	// carrying. Value is the tie-break only, so the result stays deterministic
	// if two agents ever share a rank.
	sort.Slice(out, func(i, j int) bool {
		if out[i].Priority() != out[j].Priority() {
			return out[i].Priority() < out[j].Priority()
		}
		return out[i] < out[j]
	})
	return out
}
