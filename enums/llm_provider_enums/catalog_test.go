package llm_provider_enums

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Persisted-value guards. These are the ones that matter most: values 1-4 are
// in live MongoDB documents at Organization.ClaudeAuth.Provider, moved here
// from kit/enums/claude_auth_enums with their numbering intact so that no data
// migration was needed. Reordering for tidiness would silently reinterpret
// every existing org's credential configuration.
// ---------------------------------------------------------------------------

func TestProvider_PersistedNumberingIsFrozen(t *testing.T) {
	cases := map[Provider]uint{
		AnthropicDirect:       1,
		AWSBedrock:            2,
		GoogleVertex:          3,
		AnthropicSubscription: 4,
	}
	for p, want := range cases {
		if uint(p) != want {
			t.Errorf("provider %s = %d, want %d — this value is persisted in ClaudeAuth.Provider; renumbering reinterprets existing org documents", p, uint(p), want)
		}
	}
}

func TestProvider_DisplayStringsUnchangedByTheRename(t *testing.T) {
	// String() feeds ProviderName on GET /organizations/current/claude-auth and
	// is rendered in the dashboard. The Go identifier moved from Subscription to
	// AnthropicSubscription; what a customer reads must not have moved with it.
	if got := AnthropicSubscription.String(); got != "Claude Subscription" {
		t.Errorf("AnthropicSubscription.String() = %q, want %q — user-visible text must survive the identifier rename", got, "Claude Subscription")
	}
	if got := AnthropicDirect.String(); got != "Anthropic Direct" {
		t.Errorf("AnthropicDirect.String() = %q, want %q", got, "Anthropic Direct")
	}
}

func TestProvider_AuthModeRecoversTheFusedAxis(t *testing.T) {
	// Provider fuses vendor-account and auth-mode for persistence reasons;
	// AuthMode() is how callers get the distinction back without a migration.
	cases := map[Provider]AuthMode{
		AnthropicDirect:       AuthAPIKey,
		OpenAIDirect:          AuthAPIKey,
		AWSBedrock:            AuthCloudRole,
		GoogleVertex:          AuthCloudRole,
		AnthropicSubscription: AuthSubscription,
	}
	for p, want := range cases {
		if got := p.AuthMode(); got != want {
			t.Errorf("%s.AuthMode() = %v, want %v", p, got, want)
		}
	}
	// StoresSecret drives whether "configured" can be a key-presence check.
	// Bedrock and subscriptions hold nothing control-plane-side, which is
	// exactly why HasCredentials special-cases them.
	for _, p := range []Provider{AWSBedrock, GoogleVertex, AnthropicSubscription} {
		if p.StoresSecret() {
			t.Errorf("%s must not report storing a control-plane secret", p)
		}
	}
	if !AnthropicDirect.StoresSecret() {
		t.Error("AnthropicDirect stores an API key control-plane-side")
	}
}

func TestModel_VendorIsAModelPropertyNotAProviderOne(t *testing.T) {
	// Regression guard for a real conceptual error in an earlier revision,
	// which hung Vendor off Provider and mapped AWSBedrock -> Anthropic.
	// Bedrock serves Anthropic, Amazon, Google, Meta, Mistral and Qwen models;
	// Vertex serves Claude as well as Gemini. A provider is a ROUTE to a model,
	// so it has no vendor of its own.
	if got := ClaudeOpus48.Vendor(); got != VendorAnthropic {
		t.Errorf("ClaudeOpus48.Vendor() = %v, want VendorAnthropic", got)
	}
	if got := Gpt55.Vendor(); got != VendorOpenAI {
		t.Errorf("Gpt55.Vendor() = %v, want VendorOpenAI", got)
	}
	// The same model keeps its vendor whichever provider serves it — the
	// property the earlier design could not express.
	for _, p := range ClaudeOpus48.Providers() {
		_ = p // vendor is asked of the model, never of p
	}
	for m := ClaudeHaiku45; m < MaxModel; m++ {
		if m.Vendor() == VendorUnknown {
			t.Errorf("model %s has no vendor", m)
		}
	}
}

// ---------------------------------------------------------------------------
// Wire-format guards. These strings are stored in Task documents and accepted
// by the MCP create_task tool; changing one orphans every Task holding it.
// ---------------------------------------------------------------------------

func TestModel_WireStringsAreStable(t *testing.T) {
	cases := map[Model]string{
		ClaudeHaiku45:  "claude-haiku-4-5",
		ClaudeSonnet46: "claude-sonnet-4-6",
		ClaudeOpus48:   "claude-opus-4-8",
		Gpt55:          "gpt-5.5",
		Gpt53Codex:     "gpt-5.3-codex",
		Gpt54:          "gpt-5.4",
		NovaProV1:      "nova-pro-v1",
	}
	for m, want := range cases {
		if got := m.String(); got != want {
			t.Errorf("model %d String() = %q, want %q — this is the stored wire format", uint(m), got, want)
		}
		round, err := GetModel(want)
		if err != nil || round != m {
			t.Errorf("GetModel(%q) = %v, %v; want %d round-tripping", want, round, err, uint(m))
		}
	}
}

func TestHarness_ResolveDefaultsEmptyToClaudeCode(t *testing.T) {
	// Tasks created before the agent type existed carry no AGENT_TYPE, and
	// agentbox defaults an empty value to claude-code. Diverging here would
	// route those Tasks to a different harness than the one that runs them.
	h, err := ResolveAgentType("")
	if err != nil || h != ClaudeCode {
		t.Errorf("ResolveAgentType(\"\") = %v, %v; want ClaudeCode, nil", h, err)
	}
	if _, err := ResolveAgentType("nonexistent"); err == nil {
		t.Error("ResolveAgentType must reject an unknown harness rather than defaulting it")
	}
}

// ---------------------------------------------------------------------------
// Cross-table consistency. This is what makes the package a source of truth
// rather than three lists that happen to sit together: a fact stated in one
// table and contradicted in another is caught here rather than at task time.
// ---------------------------------------------------------------------------

func TestCatalog_EveryHarnessModelPairHasAtLeastOneProvider(t *testing.T) {
	// A harness listing a model it can never actually be served is an offer the
	// product cannot honour — the user picks it and the task fails at spawn.
	for h := ClaudeCode; h < MaxAgentType; h++ {
		for _, m := range h.Models() {
			if got := ProvidersFor(h, m); len(got) == 0 {
				t.Errorf("%s lists model %s but no provider can serve it — agentTypeToProviders and modelToProviders disagree", h, m)
			}
		}
	}
}

func TestCatalog_BedrockModelsAllHaveAProfilePrefix(t *testing.T) {
	// modelToProviders claiming AWSBedrock without a prefix means discovery has
	// nothing to match on, and the logical id reaches Bedrock verbatim — the
	// exact failure the live smoke test produced ("The provided model
	// identifier is invalid").
	for m := ClaudeHaiku45; m < MaxModel; m++ {
		onBedrock := false
		for _, p := range m.Providers() {
			if p == AWSBedrock {
				onBedrock = true
			}
		}
		if onBedrock && m.BedrockProfilePrefix() == "" {
			t.Errorf("model %s lists AWSBedrock but has no profile prefix; discovery cannot resolve it", m)
		}
		if !onBedrock && m.BedrockProfilePrefix() != "" {
			t.Errorf("model %s has a profile prefix but does not list AWSBedrock as a provider", m)
		}
	}
}

// ---------------------------------------------------------------------------
// THE RULE FOR ADDING A MODEL: every entry carries a version, and every entry
// is unique — as a wire id and as a Bedrock profile prefix.
//
// This is not tidiness. Matching against inference profiles is Contains(), so
// an id or prefix that does not pin a version can match a NEIGHBOURING
// version's profile and silently run the wrong model: no error, no log line,
// resolution looking entirely successful.
// ---------------------------------------------------------------------------

func TestCatalog_EveryModelIdCarriesAVersion(t *testing.T) {
	// A digit is the enforceable proxy for "has a version": claude-opus-4-8,
	// gpt-5.5, nova-pro-v1 all qualify; a bare "nova-pro" does not. The first
	// draft of the Nova entry was exactly that, and would have matched any
	// future Nova Pro v2 profile the moment AWS published one.
	for m := ClaudeHaiku45; m < MaxModel; m++ {
		if !strings.ContainsAny(m.String(), "0123456789") {
			t.Errorf("model id %q carries no version; a versionless id can match a future version's profile", m)
		}
	}
}

func TestCatalog_ModelIdsAndPrefixesAreUnique(t *testing.T) {
	ids := map[string]Model{}
	prefixes := map[string]Model{}
	for m := ClaudeHaiku45; m < MaxModel; m++ {
		if prev, dup := ids[m.String()]; dup {
			t.Errorf("wire id %q is used by both %d and %d", m.String(), uint(prev), uint(m))
		}
		ids[m.String()] = m

		p := m.BedrockProfilePrefix()
		if p == "" {
			continue
		}
		if prev, dup := prefixes[p]; dup {
			t.Errorf("Bedrock prefix %q is shared by %s and %s; both would resolve to the same profile", p, prev, m)
		}
		prefixes[p] = m
	}
}

func TestCatalog_BedrockPrefixPinsTheModelVersion(t *testing.T) {
	// The prefix must pin the same model AND version the wire id names, leaving
	// only the date and revision to discovery — those genuinely vary per region
	// and account.
	//
	// Regression guard for a real bug: claude-sonnet-4-6 carried the prefix
	// "claude-sonnet-4", which Contains-matches
	// eu.anthropic.claude-sonnet-4-5-20250929-v1:0. A Sonnet 4.6 task would
	// have silently run Sonnet 4.5.
	for m := ClaudeHaiku45; m < MaxModel; m++ {
		p := m.BedrockProfilePrefix()
		if p == "" {
			continue
		}
		if p != m.String() {
			t.Errorf("model %s has Bedrock prefix %q; it must pin the same version as the wire id. If a model genuinely needs a different prefix, that is a deliberate decision — document it and relax this assertion for that entry only", m, p)
		}
		if !strings.ContainsAny(p, "0123456789") {
			t.Errorf("Bedrock prefix %q for %s carries no version; it can match a neighbouring version's profile", p, m)
		}
	}
}

func TestCatalog_SubscriptionIsClaudeCodeOnly(t *testing.T) {
	// The runner refuses subscription auth for any harness but claude-code — a
	// genuine `claude` CLI is what passes Anthropic's client-identity check, so
	// this is prohibited rather than merely unsupported. If the table ever
	// offered it elsewhere, the UI would present an option the runner drops.
	for h := ClaudeCode; h < MaxAgentType; h++ {
		for _, p := range h.Providers() {
			if p == AnthropicSubscription && h != ClaudeCode {
				t.Errorf("%s lists AnthropicSubscription; only claude-code may use it", h)
			}
		}
	}
}

func TestCatalog_TheWorkedExample(t *testing.T) {
	// "claude-code supports Opus 4.8, which could be provided by Anthropic
	// Direct or AWS Bedrock" — the question this package exists to answer.
	got := ProvidersFor(ClaudeCode, ClaudeOpus48)
	want := map[Provider]bool{AnthropicDirect: true, AnthropicSubscription: true, AWSBedrock: true}
	if len(got) != len(want) {
		t.Fatalf("ProvidersFor(ClaudeCode, ClaudeOpus48) = %v, want %d providers", got, len(want))
	}
	for _, p := range got {
		if !want[p] {
			t.Errorf("unexpected provider %s", p)
		}
	}
	// Codex cannot run a Claude model regardless of provider.
	if got := ProvidersFor(Codex, ClaudeOpus48); got != nil {
		t.Errorf("ProvidersFor(Codex, ClaudeOpus48) = %v, want nil — codex cannot run Claude models", got)
	}
	// opencode can run Opus, but never via the subscription.
	for _, p := range ProvidersFor(Opencode, ClaudeOpus48) {
		if p == AnthropicSubscription {
			t.Error("opencode must not be offered AnthropicSubscription")
		}
	}
}

func TestCatalog_NovaIsBedrockOnlyAndOpencodeOnly(t *testing.T) {
	// Nova is the first model that exercises the axes rather than riding on
	// the old claude-code/codex symmetry, so its shape is worth pinning.

	// Bedrock-only: Amazon offers no direct API for Nova, so its single
	// provider is a cloud ROUTE rather than its vendor. That distinction is
	// exactly why Vendor hangs off Model and not Provider.
	got := NovaProV1.Providers()
	if len(got) != 1 || got[0] != AWSBedrock {
		t.Errorf("NovaProV1.Providers() = %v, want exactly [AWSBedrock]", got)
	}
	if NovaProV1.Vendor() != VendorAmazon {
		t.Errorf("NovaProV1.Vendor() = %v, want VendorAmazon", NovaProV1.Vendor())
	}

	// Only opencode can run it — claude-code speaks the Anthropic API and
	// codex is OpenAI-only, so neither can drive a Nova model however it is
	// reached.
	harnesses := AgentTypesFor(NovaProV1)
	if len(harnesses) != 1 || harnesses[0] != Opencode {
		t.Errorf("AgentTypesFor(NovaProV1) = %v, want exactly [Opencode]", harnesses)
	}
	if ClaudeCode.Supports(NovaProV1) || Codex.Supports(NovaProV1) {
		t.Error("neither claude-code nor codex can run a Nova model")
	}

	// And it must be resolvable through opencode, or offering it is a lie.
	if len(ProvidersFor(Opencode, NovaProV1)) == 0 {
		t.Error("no provider can serve NovaPro through opencode")
	}
}

func TestCatalog_HarnessesForIsTheInverseOfModels(t *testing.T) {
	for m := ClaudeHaiku45; m < MaxModel; m++ {
		for _, h := range AgentTypesFor(m) {
			if !h.Supports(m) {
				t.Errorf("AgentTypesFor(%s) returned %s, which does not support it", m, h)
			}
		}
	}
	if len(AgentTypesFor(Gpt55)) < 2 {
		t.Error("gpt-5.5 should be runnable by both codex and opencode; the harness axis is the point")
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func TestOpencodeModelID_RendersProviderPrefix(t *testing.T) {
	// The prefix is what opencode routes on, so these strings are a contract
	// with the CLI, not cosmetic.
	cases := []struct {
		model    string
		provider Provider
		want     string
	}{
		{"claude-sonnet-4-6", AnthropicDirect, "anthropic/claude-sonnet-4-6"},
		{"gpt-5.5", OpenAIDirect, "openai/gpt-5.5"},
		{"nova-pro-v1", AWSBedrock, "amazon-bedrock/nova-pro-v1"},
		// The runner substitutes the discovered profile id before rendering, so
		// the prefix must survive a model string the catalogue never held.
		{"eu.amazon.nova-pro-v1:0", AWSBedrock, "amazon-bedrock/eu.amazon.nova-pro-v1:0"},
	}
	for _, c := range cases {
		if got := OpencodeModelID(c.model, c.provider); got != c.want {
			t.Errorf("OpencodeModelID(%q, %s) = %q, want %q", c.model, c.provider, got, c.want)
		}
	}
}

func TestOpencodeModelID_EmptyForProvidersOpencodeCannotUse(t *testing.T) {
	// Subscription auth is harness-locked to genuine claude-code — routing a
	// subscription token through a third-party harness is prohibited, not just
	// unsupported. agentTypeToProviders already excludes it; returning "" here is
	// the second line of defence, so a caller that skips that check still
	// cannot build a usable id.
	if got := OpencodeModelID("claude-opus-4-8", AnthropicSubscription); got != "" {
		t.Errorf("OpencodeModelID with AnthropicSubscription = %q, want \"\"", got)
	}
	if got := OpencodeModelID("claude-opus-4-8", GoogleVertex); got != "" {
		t.Errorf("OpencodeModelID with GoogleVertex = %q, want \"\" (not wired)", got)
	}
	// A malformed "/model" is worse than nothing — it looks valid.
	if got := OpencodeModelID("", AWSBedrock); got != "" {
		t.Errorf("OpencodeModelID with empty model = %q, want \"\"", got)
	}
}

func TestOpencodeProviderNamesCoverEveryReachableProvider(t *testing.T) {
	// If agentTypeToProviders says opencode can reach a provider, that provider
	// must have an opencode name — otherwise the catalogue offers a combination
	// no id can be rendered for, and the task fails at spawn with a malformed
	// model.
	for _, p := range Opencode.Providers() {
		if p.OpencodeName() == "" {
			t.Errorf("opencode can reach %s but it has no opencode provider name", p)
		}
	}
}

func TestIDFor_DefaultsToTheLogicalID(t *testing.T) {
	// The override map is empty on purpose, so every model currently reports
	// its logical id at every provider. This pins the FALLBACK, which is the
	// behaviour that matters when no override exists.
	for m := ClaudeHaiku45; m < MaxModel; m++ {
		for _, p := range m.Providers() {
			if got := m.IDFor(p); got != m.String() {
				t.Errorf("IDFor(%s, %s) = %q, want the logical id %q", m, p, got, m)
			}
		}
	}
}

func TestResolvesModelIDAtRuntime_IsAPropertyNotAProviderCheck(t *testing.T) {
	// Callers branch on this rather than on a provider's name, so that adding
	// Vertex does not mean hunting down every `if p == AWSBedrock`.
	runtime := map[Provider]bool{
		AWSBedrock:   true,
		GoogleVertex: true,
	}
	for _, p := range AllProviders() {
		if got := p.ResolvesModelIDAtRuntime(); got != runtime[p] {
			t.Errorf("%s.ResolvesModelIDAtRuntime() = %v, want %v", p, got, runtime[p])
		}
	}
	// A key-based provider declares its ids; nothing to discover.
	for _, p := range []Provider{AnthropicDirect, OpenAIDirect, AnthropicSubscription} {
		if p.ResolvesModelIDAtRuntime() {
			t.Errorf("%s declares its model ids; it must not require discovery", p)
		}
	}
}

func TestRuntimeResolvedProvidersHaveDiscoveryInput(t *testing.T) {
	// A provider that resolves ids at runtime needs something to match on. For
	// Bedrock that is the version-pinned profile prefix — without it discovery
	// has nothing to search for and the logical id reaches the API verbatim,
	// which is exactly how the first live run failed.
	for m := ClaudeHaiku45; m < MaxModel; m++ {
		for _, p := range m.Providers() {
			if !p.ResolvesModelIDAtRuntime() {
				continue
			}
			if p == AWSBedrock && m.BedrockProfilePrefix() == "" {
				t.Errorf("%s is served by %s, which resolves ids at runtime, but it has no profile prefix to discover with", m, p)
			}
		}
	}
}

func TestAgentTypesFor_PriorityOrder(t *testing.T) {
	// The first entry is the default agent for a model, and it comes from the
	// ENUM DECLARATION ORDER. That makes the declaration load-bearing in a way
	// nothing else would catch: inserting a new agent type in the middle would
	// silently re-route every model it can run, with no compile error.
	//
	// The rule being encoded is "prefer the specialised agent over the
	// general-purpose one".
	cases := map[Model]AgentType{
		ClaudeOpus48:   ClaudeCode, // opencode can run it too
		ClaudeSonnet46: ClaudeCode,
		ClaudeHaiku45:  ClaudeCode,
		Gpt55:          Codex,    // opencode can run it too
		NovaProV1:      Opencode, // the only option
	}
	for m, want := range cases {
		got := AgentTypesFor(m)
		if len(got) == 0 {
			t.Errorf("AgentTypesFor(%s) is empty", m)
			continue
		}
		if got[0] != want {
			t.Errorf("AgentTypesFor(%s)[0] = %s, want %s — the specialised agent must outrank the general-purpose one", m, got[0], want)
		}
	}
}

func TestProviderKeys_AreStable(t *testing.T) {
	// These are PERSISTED as the keys of LLMConfig.Providers and appear in API
	// paths. Changing one silently orphans an org's credentials for that
	// provider — the entry survives under the old key and is simply never read
	// again, so the org looks unconfigured with nothing to explain why.
	//
	// Same contract as the model wire ids: append, never rename.
	want := map[Provider]string{
		AnthropicDirect:       "anthropic-direct",
		AWSBedrock:            "aws-bedrock",
		GoogleVertex:          "google-vertex",
		AnthropicSubscription: "anthropic-subscription",
		OpenAIDirect:          "openai-direct",
	}
	for p, k := range want {
		if got := p.Key(); got != k {
			t.Errorf("%s.Key() = %q, want %q — changing a key orphans stored credentials", p, got, k)
		}
		back, err := ProviderFromKey(k)
		if err != nil || back != p {
			t.Errorf("ProviderFromKey(%q) = %v, %v; want %s round-tripping", k, back, err, p)
		}
	}
}

func TestProviderKeys_AreDistinctFromDisplayStrings(t *testing.T) {
	// The whole point of a separate vocabulary: String() is user-facing and must
	// stay free to change, so no key may equal a display string. If they ever
	// coincide, someone will "simplify" by using one for both and couple a copy
	// edit to stored data.
	for _, p := range AllProviders() {
		if p.Key() == "" {
			t.Errorf("%s has no storage key", p)
			continue
		}
		if p.Key() == p.String() {
			t.Errorf("%s: key and display string are both %q; they must stay separate vocabularies", p, p.Key())
		}
		if strings.ContainsAny(p.Key(), " .$") {
			t.Errorf("%s key %q contains a character that is awkward in a Mongo path or URL", p, p.Key())
		}
	}
}

func TestProviderKeys_AreUnique(t *testing.T) {
	seen := map[string]Provider{}
	for _, p := range AllProviders() {
		k := p.Key()
		if prev, dup := seen[k]; dup {
			t.Errorf("key %q is shared by %s and %s; one org entry would serve both", k, prev, p)
		}
		seen[k] = p
	}
}

func TestProviderFromKey_RejectsUnknown(t *testing.T) {
	// A corrupted document or a provider from a newer build must surface, not
	// read as "unconfigured".
	if _, err := ProviderFromKey("not-a-provider"); err == nil {
		t.Error("an unknown key must be an error, not a zero Provider")
	}
	if _, err := ProviderFromKey(""); err == nil {
		t.Error("an empty key must be an error")
	}
}

func TestClaudeCodeUseBedrock_IsTheCLIsOwnContract(t *testing.T) {
	// This is Anthropic's variable, not ours — a deployed claude-code is
	// looking for this exact name and value. Ours belong in typed job
	// parameters (parameters_enums.AgentProvider), not here.
	if EnvClaudeCodeUseBedrock != "CLAUDE_CODE_USE_BEDROCK" || ClaudeCodeUseBedrockValue != "1" {
		t.Errorf("EnvClaudeCodeUseBedrock=%q value=%q; claude-code reads CLAUDE_CODE_USE_BEDROCK=1", EnvClaudeCodeUseBedrock, ClaudeCodeUseBedrockValue)
	}
}

func TestApplyClaudeCodeUseBedrock_OnlyForClaudeCodeOnBedrock(t *testing.T) {
	env := map[string]string{}
	ApplyClaudeCodeUseBedrock(env, AWSBedrock, ClaudeCode)
	if env[EnvClaudeCodeUseBedrock] != ClaudeCodeUseBedrockValue {
		t.Error("claude-code on Bedrock needs its own switch set")
	}

	// No other agent understands this variable. codex ignores it; opencode
	// selects Bedrock through its model id instead. Write-if-needed means these
	// simply never get it — there is no strip step that could miss one.
	for _, c := range []struct {
		p Provider
		a AgentType
	}{
		{AWSBedrock, Codex},
		{AWSBedrock, Opencode},
		{AnthropicDirect, ClaudeCode},
		{AnthropicSubscription, ClaudeCode},
	} {
		env := map[string]string{}
		ApplyClaudeCodeUseBedrock(env, c.p, c.a)
		if _, ok := env[EnvClaudeCodeUseBedrock]; ok {
			t.Errorf("%v/%v must not get claude-code's Bedrock switch", c.p, c.a)
		}
	}
}

// The numbering is the storage format. AnthropicSubscription is 4 in live
// documents, so removing an earlier constant to "clean up" would silently
// reinterpret every subscription org's stored provider as something else.
// GoogleVertex is withheld by the catalogue, not by deletion — this pins that.
func TestProviderNumbering_IsStableAcrossReservedSlots(t *testing.T) {
	for p, want := range map[Provider]uint{
		AnthropicDirect:       1,
		AWSBedrock:            2,
		GoogleVertex:          3,
		AnthropicSubscription: 4,
		OpenAIDirect:          5,
	} {
		if uint(p) != want {
			t.Errorf("%v = %d, want %d — renumbering reinterprets stored documents", p, uint(p), want)
		}
	}
}

func TestConfigurableProviders_ExcludesReservedOnes(t *testing.T) {
	if GoogleVertex.IsConfigurable() {
		t.Error("GoogleVertex is reserved and unbuilt; offering it would accept a config nothing can serve")
	}
	// Everything an agent lists must be configurable, or an org could pick a
	// model it is then unable to run.
	for agentType := range agentTypeToModels {
		providers := agentType.Providers()
		for _, p := range providers {
			if !p.IsConfigurable() {
				t.Errorf("%v lists %v, which is not configurable", agentType, p)
			}
		}
	}
	// Enum order, so callers need not sort.
	got := ConfigurableProviders()
	for i := 1; i < len(got); i++ {
		if got[i-1] >= got[i] {
			t.Errorf("ConfigurableProviders is not in enum order: %v", got)
		}
	}
}

// Explicit values only help if nothing reintroduces a positional assumption.
// AllProviders must cover every declared provider and stay numerically ordered,
// because ConfigurableProviders and the settings UI both read it as an order.
func TestAllProviders_CoversEveryDeclaredProviderInOrder(t *testing.T) {
	all := AllProviders()
	if len(all) != len(providerToString) {
		t.Errorf("AllProviders returned %d, want %d — a declared provider is being skipped", len(all), len(providerToString))
	}
	for i := 1; i < len(all); i++ {
		if all[i-1] >= all[i] {
			t.Errorf("AllProviders is not in numeric order: %v", all)
		}
	}
	for _, p := range all {
		if !p.IsValid() {
			t.Errorf("%v is declared but not valid", p)
		}
	}
	// A value with no declaration must be invalid — the gap left by a reserved
	// or future provider cannot read as usable.
	if Provider(99).IsValid() {
		t.Error("an undeclared value must not be valid; a range check would have accepted it")
	}
}

// The invariant this pair exists to hold: an agent is put into Bedrock mode by
// exactly the same rule that decides where it reads its model. Split across two
// repos it could not be tested at all, and it had already drifted.
func TestApplyModelEnv_AgreesWithTheBedrockSwitch(t *testing.T) {
	const id = "eu.anthropic.claude-sonnet-4-6-20260101-v1:0"
	for _, p := range AllProviders() {
		for _, agentType := range []AgentType{ClaudeCode, Codex, Opencode} {
			switchEnv := map[string]string{}
			ApplyClaudeCodeUseBedrock(switchEnv, p, agentType)
			inBedrockMode := switchEnv[EnvClaudeCodeUseBedrock] == ClaudeCodeUseBedrockValue

			modelEnv := map[string]string{}
			ApplyModelEnv(modelEnv, id, p, agentType)
			_, readsAnthropicModel := modelEnv[EnvClaudeCodeBedrockModel]

			if inBedrockMode != readsAnthropicModel {
				t.Errorf("%v/%v: Bedrock mode = %v but ANTHROPIC_MODEL set = %v — an agent in the mode must be given it, and one given it must be in the mode",
					p, agentType, inBedrockMode, readsAnthropicModel)
			}

			// EnvModel is written for EVERY spawn. agentbox turns it into
			// --model, so leaving it holding the logical id means a flag that
			// overrides whatever ANTHROPIC_MODEL was resolved to.
			if modelEnv[EnvModel] != id {
				t.Errorf("%v/%v: MODEL = %q, want the resolved id — agentbox passes it as --model", p, agentType, modelEnv[EnvModel])
			}
			if readsAnthropicModel && modelEnv[EnvClaudeCodeBedrockModel] != modelEnv[EnvModel] {
				t.Errorf("%v/%v: MODEL and ANTHROPIC_MODEL disagree (%q vs %q)", p, agentType, modelEnv[EnvModel], modelEnv[EnvClaudeCodeBedrockModel])
			}
		}
	}
}

// Preference must NOT come from agentTypeToProviders' order — that list answers
// membership, and someone reordering it must not change who gets billed. So
// this passes the candidates in the WORST order and still expects the
// subscription to win.
func TestPreferredProvider_PrefersWhatTheOrgAlreadyPaysFor(t *testing.T) {
	configured := func(want ...Provider) func(Provider) bool {
		set := map[Provider]bool{}
		for _, p := range want {
			set[p] = true
		}
		return func(p Provider) bool { return set[p] }
	}
	worstOrder := []Provider{AnthropicDirect, AWSBedrock, AnthropicSubscription}

	got, ok := PreferredProvider(worstOrder, configured(AnthropicDirect, AnthropicSubscription))
	if !ok || got != AnthropicSubscription {
		t.Errorf("PreferredProvider = %v (%v); a subscribed org must not be billed metered alongside its subscription", got, ok)
	}

	// With no subscription configured, the first configured candidate wins —
	// arbitrary but deterministic, since nothing defensibly ranks an API key
	// against a cloud role.
	got, ok = PreferredProvider(worstOrder, configured(AWSBedrock, AnthropicDirect))
	if !ok || got != AnthropicDirect {
		t.Errorf("PreferredProvider = %v (%v), want the first configured candidate", got, ok)
	}

	// A candidate the org has NOT configured is never chosen.
	if got, ok := PreferredProvider(worstOrder, configured(AWSBedrock)); !ok || got != AWSBedrock {
		t.Errorf("PreferredProvider = %v (%v), want the only configured candidate", got, ok)
	}
	if _, ok := PreferredProvider(worstOrder, configured()); ok {
		t.Error("nothing configured must report no provider, not a guess")
	}
}

// The capability table and the model table are independent, so they can drift
// apart silently — a provider quietly missing from an agent looks exactly like
// one deliberately withheld. This pairs them: every provider that serves one of
// an agent's models, but which the agent cannot use, must be NAMED as an
// exclusion.
//
// Deriving capabilities instead of checking them was tried and is wrong. Today
// the derivation happens to produce the right answer, because no agent's model
// list overlaps a provider it cannot use — but that is a coincidence, and
// encoding it as a rule means adding one Bedrock-served model to codex silently
// grants codex Bedrock, which its CLI has no path to.
func TestAgentProviders_EveryGapIsNamed(t *testing.T) {
	for agentType, models := range agentTypeToModels {
		capable := map[Provider]bool{}
		for _, p := range agentType.Providers() {
			capable[p] = true
		}
		excluded := map[Provider]bool{}
		for _, p := range agentProviderExclusions[agentType] {
			excluded[p] = true
			if capable[p] {
				t.Errorf("%v excludes %v but also lists it as a capability", agentType, p)
			}
		}
		for _, m := range models {
			for _, p := range modelToProviders[m] {
				if capable[p] || excluded[p] {
					continue
				}
				t.Errorf("%v runs %v, which %v serves, but %v can use neither — if that is deliberate, name it in agentProviderExclusions",
					agentType, m, p, agentType)
			}
		}
	}
}

// A capability an agent cannot exercise is worth knowing about: it means the
// agent lists a provider none of its models are served by, so the intersection
// is empty and the capability is inert. Not an error — claude-code could
// legitimately support a provider before any model is offered through it — but
// it should be visible rather than silent.
func TestAgentProviders_ReportsInertCapabilities(t *testing.T) {
	for agentType, models := range agentTypeToModels {
		served := map[Provider]bool{}
		for _, m := range models {
			for _, p := range modelToProviders[m] {
				served[p] = true
			}
		}
		for _, p := range agentType.Providers() {
			if !served[p] {
				t.Logf("note: %v can talk to %v but runs no model served there — the capability is inert", agentType, p)
			}
		}
	}
}
