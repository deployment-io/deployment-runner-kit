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
//
// ⚠️ ORDER IS LOAD-BEARING and new entries go at the END. ModelForTier walks
// this list and takes the first non-legacy match, so inserting a model
// mid-list silently changes what "high complexity" resolves to for every
// existing Task. Appending cannot change it on its own — which is the point.
// Moving the frontier answer on is a deliberate edit to agentTypeToDefaultModel
// or legacyModels, never a side effect of adding a model.
var agentTypeToModels = map[AgentType][]Model{
	ClaudeCode: {ClaudeHaiku45, ClaudeSonnet45, ClaudeSonnet46, ClaudeOpus45, ClaudeOpus48,
		ClaudeSonnet5, ClaudeOpus5, ClaudeFable5},
	Codex: {Gpt55, Gpt53Codex, Gpt54, Gpt56Sol, Gpt56Terra, Gpt56Luna},
	Opencode: {ClaudeHaiku45, ClaudeSonnet45, ClaudeSonnet46, ClaudeOpus45, ClaudeOpus48, Gpt55, NovaProV1,
		ClaudeSonnet5, ClaudeOpus5, ClaudeFable5,
		// The Bedrock-only lineup. opencode is the only agent that can reach
		// these: claude-code and codex each speak one vendor's API, while
		// opencode routes by provider id — which is the reason it exists here.
		Qwen3Coder480B, Qwen3CoderNext, DeepSeekV32, Glm47, Glm5, MinimaxM25, Grok43,
		// The 5.6 family, on OpenAIDirect like Gpt55 above — opencode holds an
		// OpenAI key of its own, so it reaches them without going near codex.
		// Appended, so opencode's tier answers are unchanged: Opus 5 still wins
		// frontier and Haiku 4.5 still wins fast, both listed earlier.
		Gpt56Sol, Gpt56Terra, Gpt56Luna},
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
	// The 5.6 family. OpenAIDirect ALONE, deliberately: codex authenticates
	// against api.openai.com, and no other route has been verified for these.
	// Bedrock hosts only OpenAI's open-weight gpt-oss models, and whether the
	// gateways carry 5.6 yet is a question for the day someone checks the
	// registry — not an assumption to make here, where an unserved provider
	// shows up as a picker entry whose Task fails at pickup.
	Gpt56Sol:   {OpenAIDirect},
	Gpt56Terra: {OpenAIDirect},
	Gpt56Luna:  {OpenAIDirect},
	// Bedrock-only: Amazon does not offer Nova through a direct API, so this
	// is the first model whose single provider is a cloud route rather than
	// its vendor. claude-code and codex cannot run it — agentTypeToModels keeps
	// it to opencode.
	NovaProV1: {AWSBedrock},
	// Same providers as their 4.6/4.8 siblings — the direct API and a
	// subscription still serve them, and Bedrock is where they matter most.
	ClaudeSonnet45: {AnthropicDirect, AnthropicSubscription, AWSBedrock},
	ClaudeOpus45:   {AnthropicDirect, AnthropicSubscription, AWSBedrock},
	ClaudeSonnet5:  {AnthropicDirect, AnthropicSubscription, AWSBedrock},
	ClaudeOpus5:    {AnthropicDirect, AnthropicSubscription, AWSBedrock},
	ClaudeFable5:   {AnthropicDirect, AnthropicSubscription, AWSBedrock},
	// The open-weight lineup. No vendor here has a direct provider of its own,
	// so these reach us through a cloud, a platform or a gateway.
	//
	// Bedrock was the only route until Novita and OpenRouter, and that was a
	// problem: AWS grants model access PER ACCOUNT PER REGION, so all seven
	// shipped uninvocable on an account that had not enabled them. Both new
	// providers need one API key and no per-model step.
	//
	// ⚠️ THE THREE ROUTES DO NOT SERVE THE SAME SET, so these lists are not
	// copy-paste. Verified against models.opencode.ai — the registry opencode
	// actually resolves against:
	//
	//	Novita has no Grok         — a GPU platform cannot host a proprietary
	//	                             model; only a gateway can route to xAI
	//	OpenRouter has no 480B     — it carries Qwen3 Coder Next instead
	//
	// Which is why both are worth having rather than either alone.
	Qwen3Coder480B: {AWSBedrock, Novita},
	Qwen3CoderNext: {AWSBedrock, Novita, OpenRouter},
	DeepSeekV32:    {AWSBedrock, Novita, OpenRouter},
	Glm47:          {AWSBedrock, Novita, OpenRouter},
	Glm5:           {AWSBedrock, Novita, OpenRouter},
	MinimaxM25:     {AWSBedrock, Novita, OpenRouter},
	Grok43:         {AWSBedrock, OpenRouter},
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
	// Novita and OpenRouter are opencode-only until proven otherwise:
	// claude-code speaks Anthropic's API and codex authenticates against
	// api.openai.com. OpenRouter may well expose a compatible surface, but that
	// is UNVERIFIED — and this table exists because Bedrock hosts OpenAI's
	// gpt-oss models and codex still cannot reach one.
	Opencode: {AnthropicDirect, AWSBedrock, OpenAIDirect, Novita, OpenRouter},
}

// Models returns the models this harness OFFERS — retired ones excluded.
//
// Filtered here rather than at each call site because this is what every
// offering path already reads, kit's SupportedModels and IsModelSupported
// included. A disabled model therefore stops being creatable without kit
// needing to learn the concept.
//
// Use AllModels when the question is capability or history rather than what to
// offer.
func (h AgentType) Models() []Model {
	out := make([]Model, 0, len(agentTypeToModels[h]))
	for _, m := range agentTypeToModels[h] {
		if !m.IsDisabled() {
			out = append(out, m)
		}
	}
	return out
}

// AllModels returns every model this harness can run, INCLUDING retired ones.
//
// For callers asking what a harness is capable of, or rendering a model a Task
// already holds. Offering paths want Models.
func (h AgentType) AllModels() []Model { return agentTypeToModels[h] }

// Providers returns the provider APIs this harness can talk to. Order carries
// NO preference — see PreferredProvider.
func (h AgentType) Providers() []Provider { return agentProviderCapabilities[h] }

// Providers returns every provider that can serve this model, ignoring which
// harness is asking. Use ProvidersFor when a harness is known.
func (m Model) Providers() []Provider { return modelToProviders[m] }

// Supports reports whether the harness can run the model at all.
//
// Deliberately reads the FULL list, retired models included: a Task already
// holding a disabled model is still runnable, and this answers capability
// rather than what to offer. ProvidersFor depends on it, which is what keeps a
// disabled model resolving at spawn instead of failing as "not possible".
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
	// The registry ids, not our display names. "novita-ai" reads oddly and is
	// correct — it is the key in models.opencode.ai, and opencode resolves the
	// "provider/model" string against exactly that.
	Novita:     "novita-ai",
	OpenRouter: "openrouter",
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
// vendor decides whether a Bedrock id keeps its cross-region geography prefix,
// and is why this takes a Vendor rather than sniffing the id's own vendor
// segment: the rule is a fact about who MAKES the model, so the type that
// already carries that fact should state it. Pass VendorUnknown when the model
// is not in the catalogue and the id will be left alone.
//
// Returns "" when the provider has no opencode representation, so callers get
// an empty result rather than a malformed "/model" string.
func OpencodeModelID(modelID string, p Provider, vendor Vendor) string {
	name := p.OpencodeName()
	if name == "" || modelID == "" {
		return ""
	}
	// opencode routes Bedrock through its own registry (models.dev), which
	// carries geography-prefixed ids for some vendors and only base ids for
	// others — see opencodeBedrockGeographyVendors for the split and the counts.
	//
	// Keeping the prefix matters where it belongs: for Anthropic the prefixed
	// id IS the cross-region inference profile, which newer Claude models on
	// Bedrock generally require, so dropping to the base id would silently ask
	// for on-demand throughput they may not offer.
	//
	// A live Nova run exposed the other direction. Discovery resolved
	// eu.amazon.nova-pro-v1:0 and opencode rejected it with
	// ProviderModelNotFoundError, suggesting amazon.nova-pro-v1:0.
	//
	// Only the geography goes; the dated revision stays, since that is what
	// discovery exists to find.
	if p == AWSBedrock && !vendor.OpencodeKeepsBedrockGeography() {
		modelID = stripBedrockGeography(modelID)
	}
	return name + "/" + modelID
}

// stripBedrockGeography removes the cross-region prefix from a Bedrock id.
//
// WHETHER to call this is the vendor's business, not this function's — it just
// performs the removal.
func stripBedrockGeography(modelID string) string {
	for _, geo := range []string{"eu.", "us.", "au.", "jp.", "apac.", "global."} {
		if rest, found := strings.CutPrefix(modelID, geo); found {
			return rest
		}
	}
	return modelID
}

// ConfigurableProviders returns every provider an org can actually configure,
// in DISPLAY order.
//
// DERIVED, not a second list: a provider is configurable exactly when some
// agent can use it. That makes the catalogue the only place a provider is
// switched on, so a reserved-but-unbuilt one (GoogleVertex) cannot be offered
// by the settings UI or accepted by the API without someone first adding it to
// a real agent — rather than each caller carrying its own reject list that
// drifts.
//
// Ordering used to be AllProviders(), i.e. by the PERSISTED ENUM NUMBER — so
// the settings page was laid out by the order providers happened to be added,
// and could not be changed without renumbering values that live in customer
// documents. displayRank replaces that with a stated rule, the same move
// already made for provider preference, capability tier and agent priority.
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
	sort.SliceStable(out, func(i, j int) bool {
		return displayRank(out[i]) < displayRank(out[j])
	})
	return out
}

// displayRank orders the settings cards. Lower sorts first.
//
// THE RULE IS "how much does this card ask of the user", not familiarity or
// alphabet — so it generalises to providers not yet built instead of needing a
// new opinion each time:
//
//  1. paste an API key          Anthropic, OpenAI, Novita, OpenRouter
//  2. an optional fallback key  Claude Subscription
//  3. nothing to paste at all   AWS Bedrock — the runner assumes a role, so
//     there is no key and no console to visit
//
// Ties keep AllProviders() order via a stable sort, which is why the API-key
// group stays in the order those providers were added rather than shuffling.
// Vertex lands in group 3 with Bedrock when it ships, without a decision.
func displayRank(p Provider) int {
	switch p.AuthMode() {
	case AuthAPIKey:
		return 1
	case AuthSubscription:
		return 2
	default:
		// AuthCloudRole and anything unrecognised. Last is the safe place for
		// a provider whose card shape we do not know.
		return 3
	}
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
// That rule is now the DEFAULT, not the answer: an org that has stated its own
// order gets that instead, and the guessing stops. See ProviderPreference.
func PreferredProvider(candidates []Provider, isConfigured func(Provider) bool, pref ProviderPreference) (Provider, bool) {
	// An override is the most specific thing anyone said, so it wins — but only
	// if it is still usable. A model that loses a provider, or an org that
	// removes a credential, leaves a stale override behind, and honouring one
	// blindly would break routing on a setting the user has forgotten setting.
	if pref.Override.IsValid() && isConfigured(pref.Override) && servesModel(candidates, pref.Override) {
		return pref.Override, true
	}
	best, found := Provider(0), false
	for _, p := range candidates {
		if !isConfigured(p) {
			continue
		}
		if !found || pref.rank(p) < pref.rank(best) {
			best, found = p, true
		}
	}
	return best, found
}

// ProviderPreference is an org's own routing choice, and the zero value means
// "we have not been told" — which is why every field is optional and the
// default rule below still applies.
type ProviderPreference struct {
	// Order ranks providers highest-first. NOT an allowlist: a provider absent
	// from it is still eligible, just last. That distinction matters because a
	// model may have exactly one route (Nova on Bedrock, Grok on OpenRouter),
	// and an order that silently excluded it would make the model unrunnable
	// through a setting that never mentioned it.
	Order []Provider
	// Override pins ONE model to ONE provider, for the case an org-wide order
	// cannot express: Qwen3 Coder Next serves 262k of output through one
	// provider and 65k through another, so the general preference is right for
	// every other model and wrong for that one.
	Override Provider
}

// rank scores a provider for this preference. Lower wins.
func (pref ProviderPreference) rank(p Provider) int {
	if len(pref.Order) == 0 {
		// Nobody has said anything, so fall back to the paid-for-already rule.
		return preferenceRank(p)
	}
	// ⚠️ AN EXPLICIT ORDER OUTRANKS THE SUBSCRIPTION RULE. That rule exists to
	// stop us silently billing an org twice, which is worth defending against a
	// coin toss — but not against someone who has said what they want. Guessing
	// is the thing being replaced here; overriding a stated choice to protect
	// the org from itself would keep the guessing and add condescension.
	for i, ordered := range pref.Order {
		if ordered == p {
			return i
		}
	}
	// Ranked after everything named, in catalogue order among themselves.
	return len(pref.Order) + int(p)
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
	//
	// ⚠️ Moving this moves TWO things: what a picker preselects, and what a
	// high-complexity session resolves to — the default wins its tier outright
	// in ModelForTier. Change it deliberately, never as a side effect of adding
	// a model.
	ClaudeCode: ClaudeOpus5,
	// OpenAI's best coding model, and their recommended default for Codex —
	// same rationale as the line above.
	//
	// ⚠️ Same two-for-one as ClaudeCode: this also decides what a
	// high-complexity codex session resolves to, because the default wins its
	// tier outright in ModelForTier. Both movements are intended here. Sol is
	// the only frontier-tier model codex runs, so the tier answer would have
	// landed on it either way; what the default adds is the picker's
	// preselection, moving off Gpt55.
	Codex: Gpt56Sol,
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
	for _, m := range h.Models() {
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

// servesModel reports whether a provider is among a model's candidates. Named
// for what the caller is asking rather than the mechanics, since "contains"
// already means something else in this package's tests.
func servesModel(candidates []Provider, want Provider) bool {
	for _, p := range candidates {
		if p == want {
			return true
		}
	}
	return false
}
