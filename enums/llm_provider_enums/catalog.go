package llm_provider_enums

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
	ClaudeCode: {ClaudeHaiku45, ClaudeSonnet46, ClaudeOpus48},
	Codex:      {Gpt55, Gpt53Codex, Gpt54},
	Opencode:   {ClaudeHaiku45, ClaudeSonnet46, ClaudeOpus48, Gpt55, NovaProV1},
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
}

// agentTypeToProviders lists which providers each harness can authenticate to.
//
// The two constraints worth knowing, both enforced in deployment-runner today
// and previously only discoverable by reading it:
//
//   - Codex is OpenAI-only. It has no Bedrock path at all.
//   - Only claude-code may use AnthropicSubscription. maybeApplyClaudeSubscriptionAuth
//     returns early for any other AGENT_TYPE, because a genuine `claude` CLI is
//     what passes Anthropic's client-identity check; routing a subscription
//     token through another agent is prohibited, not merely unsupported.
//
// GoogleVertex is listed for claude-code because the harness supports it, even
// though no model above is served by it yet — the intersection keeps that from
// ever being offered.
var agentTypeToProviders = map[AgentType][]Provider{
	ClaudeCode: {AnthropicDirect, AnthropicSubscription, AWSBedrock},
	Codex:      {OpenAIDirect},
	Opencode:   {AnthropicDirect, AWSBedrock, OpenAIDirect},
}

// Models returns the models this harness can run.
func (h AgentType) Models() []Model { return agentTypeToModels[h] }

// Providers returns the providers this harness can authenticate to.
func (h AgentType) Providers() []Provider { return agentTypeToProviders[h] }

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
	reachable := make(map[Provider]bool, len(agentTypeToProviders[h]))
	for _, p := range agentTypeToProviders[h] {
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
	for h := ClaudeCode; h < MaxAgentType; h++ {
		if h.Supports(m) {
			out = append(out, h)
		}
	}
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
	return name + "/" + modelID
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
	for _, providers := range agentTypeToProviders {
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
