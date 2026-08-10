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
		ClaudeSonnet45: "claude-sonnet-4-5",
		ClaudeOpus45:   "claude-opus-4-5",
		ClaudeSonnet5:  "claude-sonnet-5",
		ClaudeOpus5:    "claude-opus-5",
		ClaudeFable5:   "claude-fable-5",
		Qwen3Coder480B: "qwen3-coder-480b",
		Qwen3CoderNext: "qwen3-coder-next",
		DeepSeekV32:    "deepseek-v3.2",
		Glm47:          "glm-4.7",
		Glm5:           "glm-5",
		MinimaxM25:     "minimax-m2.5",
		Grok43:         "grok-4.3",
	}
	// EXHAUSTIVE. Without this, adding a model leaves its wire id unpinned and
	// this test still passes — which is exactly what happened when Sonnet 4.5
	// and Opus 4.5 were added. Pinning must be a deliberate step, not one that
	// depends on remembering.
	if len(cases) != len(modelToString) {
		t.Errorf("%d models declared, %d pinned — add the new one's wire id here on purpose", len(modelToString), len(cases))
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
	for _, h := range AllAgentTypes() {
		for _, m := range h.Models() {
			if got := ProvidersFor(h, m); len(got) == 0 {
				t.Errorf("%s lists model %s but no provider can serve it — agentTypeToProviders and modelToProviders disagree", h, m)
			}
		}
	}
}

// A Bedrock model reaches its concrete id one of TWO ways, and must declare
// exactly one of them.
//
//	profile prefix   — a cross-region inference profile exists; discovery
//	                   resolves the dated revision at run time
//	declared id      — no profile exists (the Qwen/DeepSeek/GLM/MiniMax/Grok
//	                   lineup), so the id is stated in modelProviderID
//
// NEITHER means the logical id reaches Bedrock verbatim — "The provided model
// identifier is invalid", which is how the first live run failed. BOTH is a
// contradiction: discovery would overwrite the declared id, so one of the two
// statements is dead and nobody could tell which was intended.
func TestCatalog_BedrockModelsDeclareExactlyOneIDMechanism(t *testing.T) {
	for m := ClaudeHaiku45; m < MaxModel; m++ {
		onBedrock := false
		for _, p := range m.Providers() {
			if p == AWSBedrock {
				onBedrock = true
			}
		}
		prefix := m.BedrockProfilePrefix() != ""
		declared := m.IDFor(AWSBedrock) != m.String()

		switch {
		case onBedrock && !prefix && !declared:
			t.Errorf("model %s lists AWSBedrock but declares neither a profile prefix nor a Bedrock id; the logical id would reach the API verbatim", m)
		case onBedrock && prefix && declared:
			t.Errorf("model %s declares BOTH a profile prefix and a Bedrock id; discovery would overwrite the declared id, so one of them is dead", m)
		case !onBedrock && prefix:
			t.Errorf("model %s has a profile prefix but does not list AWSBedrock as a provider", m)
		case !onBedrock && declared:
			t.Errorf("model %s declares a Bedrock id but does not list AWSBedrock as a provider", m)
		}
	}
}

// The declared ids are what OPENCODE resolves against, so their shape is a
// contract with models.dev rather than with Bedrock. Each must carry the
// vendor segment models.dev uses — a bare "glm-5" is not findable there, and
// that is precisely the class of failure the Nova run exposed.
func TestCatalog_DeclaredBedrockIDsCarryTheirVendorSegment(t *testing.T) {
	segment := map[Vendor]string{
		VendorQwen:     "qwen.",
		VendorDeepSeek: "deepseek.",
		VendorZAI:      "zai.",
		VendorMiniMax:  "minimax.",
		VendorXAI:      "xai.",
		VendorAmazon:   "amazon.",
	}
	for m := ClaudeHaiku45; m < MaxModel; m++ {
		id := m.IDFor(AWSBedrock)
		if id == m.String() {
			continue // resolved by discovery, not declared
		}
		want, ok := segment[m.Vendor()]
		if !ok {
			t.Errorf("model %s declares Bedrock id %q but its vendor %v has no known registry segment", m, id, m.Vendor())
			continue
		}
		if !strings.HasPrefix(id, want) {
			t.Errorf("declared Bedrock id for %s is %q, want it to start with %q — opencode looks it up in models.dev by that name", m, id, want)
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
	for _, h := range AllAgentTypes() {
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
		vendor   Vendor
		want     string
	}{
		{"claude-sonnet-4-6", AnthropicDirect, VendorAnthropic, "anthropic/claude-sonnet-4-6"},
		{"gpt-5.5", OpenAIDirect, VendorOpenAI, "openai/gpt-5.5"},
		{"nova-pro-v1", AWSBedrock, VendorAmazon, "amazon-bedrock/nova-pro-v1"},

		// Amazon's own models have NO geography-prefixed entries in opencode's
		// registry, so the prefix discovery found must come off. This is the
		// case a live Nova run failed on.
		{"eu.amazon.nova-pro-v1:0", AWSBedrock, VendorAmazon, "amazon-bedrock/amazon.nova-pro-v1:0"},
		{"us.amazon.nova-pro-v1:0", AWSBedrock, VendorAmazon, "amazon-bedrock/amazon.nova-pro-v1:0"},

		// Anthropic's ARE in the registry prefixed, and the prefixed id IS the
		// cross-region inference profile — which newer Claude models on Bedrock
		// generally require. Stripping here would quietly request on-demand
		// throughput they may not offer, so it must pass through untouched.
		{"eu.anthropic.claude-sonnet-4-5-20250929-v1:0", AWSBedrock, VendorAnthropic,
			"amazon-bedrock/eu.anthropic.claude-sonnet-4-5-20250929-v1:0"},
		{"us.anthropic.claude-opus-4-5-20251101-v1:0", AWSBedrock, VendorAnthropic,
			"amazon-bedrock/us.anthropic.claude-opus-4-5-20251101-v1:0"},

		// An id that already carries no geography is left alone either way.
		{"amazon.nova-pro-v1:0", AWSBedrock, VendorAmazon, "amazon-bedrock/amazon.nova-pro-v1:0"},

		// A model outside the catalogue has no vendor, and an unknown vendor
		// must not be rewritten on a guess.
		{"eu.acme.something-v1:0", AWSBedrock, VendorUnknown,
			"amazon-bedrock/eu.acme.something-v1:0"},

		// Only Bedrock is touched — an id starting with a geography-looking
		// token at another provider must survive.
		{"eu.something", AnthropicDirect, VendorAnthropic, "anthropic/eu.something"},
		{"eu.something", OpenAIDirect, VendorAmazon, "openai/eu.something"},
	}
	for _, c := range cases {
		if got := OpencodeModelID(c.model, c.provider, c.vendor); got != c.want {
			t.Errorf("OpencodeModelID(%q, %s, %v) = %q, want %q", c.model, c.provider, c.vendor, got, c.want)
		}
	}
}

// Every geography AWS uses must be recognised, for a vendor that strips. A
// prefix missing from the list is a silent pass-through: the id keeps its
// geography, opencode has no entry for it, and the run fails at spawn.
func TestOpencodeModelID_StripsEveryGeography(t *testing.T) {
	for _, geo := range []string{"eu.", "us.", "au.", "jp.", "apac.", "global."} {
		got := OpencodeModelID(geo+"amazon.nova-pro-v1:0", AWSBedrock, VendorAmazon)
		if want := "amazon-bedrock/amazon.nova-pro-v1:0"; got != want {
			t.Errorf("geography %q: got %q, want %q", geo, got, want)
		}
	}
}

// The strip decision belongs to the vendor, so it must follow the VENDOR
// argument rather than the id's own vendor segment. These two disagree on
// purpose: a caller that passed the wrong vendor should produce the wrong
// answer here, which is what makes the argument load-bearing rather than
// decorative.
func TestOpencodeModelID_FollowsVendorArgNotTheIDText(t *testing.T) {
	if got := OpencodeModelID("eu.amazon.nova-pro-v1:0", AWSBedrock, VendorAnthropic); got !=
		"amazon-bedrock/eu.amazon.nova-pro-v1:0" {
		t.Errorf("VendorAnthropic should keep the prefix, got %q", got)
	}
	if got := OpencodeModelID("eu.anthropic.claude-sonnet-4-5-20250929-v1:0", AWSBedrock, VendorAmazon); got !=
		"amazon-bedrock/anthropic.claude-sonnet-4-5-20250929-v1:0" {
		t.Errorf("VendorAmazon should strip the prefix, got %q", got)
	}
}

// Every model the catalogue offers on Bedrock must have a vendor whose strip
// rule was actually considered. Adding a Bedrock-hosted model from a new vendor
// fails HERE, at CI, rather than at spawn on a customer's runner — which is the
// whole reason the rule hangs off the Vendor enum.
func TestBedrockModelsHaveAConsideredGeographyRule(t *testing.T) {
	// ⚠️ KEYED BY MODEL, NOT VENDOR, even though the RULE is vendor-level. Two
	// vendors split per model in the registry — meta (llama3 bare, llama4 us.)
	// and deepseek (v3.2 bare, r1 us.) — so a vendor-keyed guard would wave
	// DeepSeek R1 through the moment V3.2 is present, and R1 would break
	// silently. Model-keyed costs one line per model and cannot be outgrown.
	const (
		keep  = true  // models.dev carries geography-prefixed entries
		strip = false // models.dev carries only the base id
	)
	considered := map[Model]bool{
		ClaudeHaiku45:  keep,
		ClaudeSonnet45: keep,
		ClaudeSonnet46: keep,
		ClaudeOpus45:   keep,
		ClaudeOpus48:   keep,
		ClaudeSonnet5:  keep,
		ClaudeOpus5:    keep,
		ClaudeFable5:   keep,
		NovaProV1:      strip,
		Qwen3Coder480B: strip,
		Qwen3CoderNext: strip,
		DeepSeekV32:    strip,
		Glm47:          strip,
		Glm5:           strip,
		MinimaxM25:     strip,
		Grok43:         strip,
	}
	for m := ClaudeHaiku45; m < MaxModel; m++ {
		onBedrock := false
		for _, p := range m.Providers() {
			if p == AWSBedrock {
				onBedrock = true
			}
		}
		if !onBedrock {
			continue
		}
		want, ok := considered[m]
		if !ok {
			t.Errorf("model %s is on Bedrock but its geography rule was never "+
				"checked — look it up in models.dev, then record it here and, if it "+
				"needs keeping, in opencodeBedrockGeographyVendors", m)
			continue
		}
		// The vendor-level rule must agree with what was checked per model. A
		// mismatch means that vendor has split and the rule can no longer be
		// stated vendor-wide.
		if got := m.Vendor().OpencodeKeepsBedrockGeography(); got != want {
			t.Errorf("model %s: vendor %v says keep=%v but the registry says keep=%v — "+
				"this vendor now splits per model, so the rule needs a per-model override",
				m, m.Vendor(), got, want)
		}
	}
}

func TestOpencodeModelID_EmptyForProvidersOpencodeCannotUse(t *testing.T) {
	// Subscription auth is harness-locked to genuine claude-code — routing a
	// subscription token through a third-party harness is prohibited, not just
	// unsupported. agentTypeToProviders already excludes it; returning "" here is
	// the second line of defence, so a caller that skips that check still
	// cannot build a usable id.
	if got := OpencodeModelID("claude-opus-4-8", AnthropicSubscription, VendorAnthropic); got != "" {
		t.Errorf("OpencodeModelID with AnthropicSubscription = %q, want \"\"", got)
	}
	if got := OpencodeModelID("claude-opus-4-8", GoogleVertex, VendorAnthropic); got != "" {
		t.Errorf("OpencodeModelID with GoogleVertex = %q, want \"\" (not wired)", got)
	}
	// A malformed "/model" is worse than nothing — it looks valid.
	if got := OpencodeModelID("", AWSBedrock, VendorAmazon); got != "" {
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
	// Overrides exist only for Bedrock, and only for the models with no
	// inference profile. Everywhere else the logical id must come back
	// unchanged — this pins the FALLBACK, which is what almost every lookup
	// hits.
	//
	// Deliberately asserts on EVERY provider except that one pair, rather than
	// listing the models expected to be clean: an override added to a second
	// provider by accident fails here rather than at spawn.
	for m := ClaudeHaiku45; m < MaxModel; m++ {
		for _, p := range m.Providers() {
			if p == AWSBedrock && m.BedrockProfilePrefix() == "" {
				continue // declared id, covered by the two tests above
			}
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
	// A provider that resolves ids at runtime needs something to match on, OR a
	// declared id to fall back to. Neither means the logical id reaches the API
	// verbatim, which is exactly how the first live run failed.
	//
	// ResolvesModelIDAtRuntime is a property of the PROVIDER, so it stays true
	// for Bedrock even for models that have nothing to discover. That is not a
	// contradiction: discovery runs, finds no profile, logs it, and leaves the
	// declared id in place — see applyAgentModelEnv, which only overwrites on a
	// non-empty result.
	for m := ClaudeHaiku45; m < MaxModel; m++ {
		for _, p := range m.Providers() {
			if !p.ResolvesModelIDAtRuntime() || p != AWSBedrock {
				continue
			}
			if m.BedrockProfilePrefix() == "" && m.IDFor(p) == m.String() {
				t.Errorf("%s is served by %s, which resolves ids at runtime, but it has neither a profile prefix to discover with nor a declared id to fall back to", m, p)
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
	if len(want) != len(providerKey) {
		t.Errorf("%d providers have keys, %d pinned — a new slug must be pinned deliberately", len(providerKey), len(want))
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

// A model with no display name renders as its wire id — legible but wrong for a
// picker, and the kind of gap nobody notices until a customer sees it.
func TestModel_EveryCatalogueModelHasADisplayName(t *testing.T) {
	for _, models := range agentTypeToModels {
		for _, m := range models {
			if _, named := modelToDisplayName[m]; !named {
				t.Errorf("%v has no display name; a picker would show its wire id", m)
			}
			// The label must NOT be the wire id: they have different jobs, and
			// letting them coincide invites someone to use one for the other.
			if m.DisplayName() == m.String() {
				t.Errorf("%v: display name equals the wire id %q", m, m.String())
			}
		}
	}
}

// A default no picker can offer is worse than none: the UI preselects something
// the org cannot run, and the first Task fails.
func TestAgentType_DefaultModelIsOneItRuns(t *testing.T) {
	for _, h := range AllAgentTypes() {
		def := h.DefaultModel()
		if def == 0 {
			t.Errorf("%v has no default model", h)
			continue
		}
		if !h.Supports(def) {
			t.Errorf("%v defaults to %v, which it cannot run", h, def)
		}
	}
}

// ---------------------------------------------------------------------------
// Retirement. Disabling is the soft delete for a persisted enum: the wire id is
// in Task documents, so the entry must keep resolving even once nothing offers
// it. These cases pin BOTH halves — a disabled model that stopped resolving
// would be a delete with extra steps, and one that kept being offered would be
// no retirement at all.
// ---------------------------------------------------------------------------

func TestDisabled_StillResolvesForTasksThatAlreadyHoldIt(t *testing.T) {
	if !NovaProV1.IsDisabled() {
		t.Fatal("this test needs a disabled model to be meaningful")
	}
	// The reason we disable instead of deleting: a stored Task holds this
	// string, and every one of these is on the path from that string to a
	// running agent, or to a rendered row in the UI.
	m, err := GetModel("nova-pro-v1")
	if err != nil || m != NovaProV1 {
		t.Fatalf("GetModel(\"nova-pro-v1\") = %v, %v — a retired model must still parse, or history and in-flight Jobs break", m, err)
	}
	if m.DisplayName() == "" || m.DisplayName() == m.String() {
		t.Errorf("DisplayName() = %q; a retired model must still render as itself", m.DisplayName())
	}
	if m.Vendor() != VendorAmazon {
		t.Errorf("Vendor() = %v, want VendorAmazon", m.Vendor())
	}
	// Capability is not offering: the harness can still RUN it, which is what
	// keeps ProvidersFor non-empty and lets an existing Job spawn.
	if !Opencode.Supports(m) {
		t.Error("Opencode.Supports(nova-pro-v1) = false; a retired model must stay runnable for Tasks that already chose it")
	}
	if len(ProvidersFor(Opencode, m)) == 0 {
		t.Error("ProvidersFor returned nothing for a retired model; it would fail at spawn as 'not possible'")
	}
	if m.IDFor(AWSBedrock) == "" {
		t.Error("IDFor returned empty for a retired model")
	}
}

func TestDisabled_IsOfferedNowhere(t *testing.T) {
	for _, h := range AllAgentTypes() {
		for _, m := range h.Models() {
			if m.IsDisabled() {
				t.Errorf("%s.Models() offers retired model %s", h, m)
			}
		}
		// Every tier, since a recommendation is the one path that could hand a
		// user a retired model they never picked.
		for _, tr := range []Tier{TierFast, TierBalanced, TierFrontier} {
			if got := ModelForTier(h, tr); got.IsDisabled() {
				t.Errorf("ModelForTier(%s, %s) recommends retired model %s", h, tr, got)
			}
		}
		// AllModels is the deliberate exception — it answers capability.
		if len(h.AllModels()) < len(h.Models()) {
			t.Errorf("%s.AllModels() is smaller than Models(); it must be the superset", h)
		}
	}
}

// A default that is retired would preselect a model no picker lists.
func TestDisabled_NoAgentDefaultsToARetiredModel(t *testing.T) {
	for _, h := range AllAgentTypes() {
		if h.DefaultModel().IsDisabled() {
			t.Errorf("%s defaults to retired model %s", h, h.DefaultModel())
		}
	}
}

func TestModelsFor_OffersOnlyWhatTheOrgCanServe(t *testing.T) {
	bedrockOnly := func(p Provider) bool { return p == AWSBedrock }
	got := ModelsFor(Opencode, bedrockOnly)
	names := map[Model]bool{}
	for _, m := range got {
		names[m] = true
	}
	// Qwen3 Coder is Bedrock-only, so a Bedrock org must see it. This was
	// Nova's job until Nova was retired — which is what the next case checks.
	if !names[Qwen3Coder480B] {
		t.Errorf("opencode/Bedrock omits qwen3-coder-480b; got %v", got)
	}
	// ...and a RETIRED Bedrock-only model must not appear, however well the org
	// is credentialed. Being serveable is not being offered.
	if names[NovaProV1] {
		t.Errorf("opencode/Bedrock still offers retired nova-pro-v1; got %v", got)
	}
	// Claude models are served by Bedrock too.
	if !names[ClaudeSonnet46] {
		t.Errorf("opencode/Bedrock omits claude-sonnet-4-6; got %v", got)
	}
	// GPT-5.5 is NOT on Bedrock — it is OpenAI's own model, so it must be
	// withheld even though opencode can reach Bedrock perfectly well.
	if names[Gpt55] {
		t.Errorf("opencode/Bedrock offers gpt-5.5, which no Bedrock provider serves; got %v", got)
	}
	// Codex has no Bedrock transport at all, so nothing is offerable.
	if got := ModelsFor(Codex, bedrockOnly); len(got) != 0 {
		t.Errorf("codex on a Bedrock-only org offers %v, but it cannot reach Bedrock", got)
	}
	// Nothing configured offers nothing — never a fallback to "all models".
	if got := ModelsFor(ClaudeCode, func(Provider) bool { return false }); len(got) != 0 {
		t.Errorf("an unconfigured org was offered %v", got)
	}
}

// An agent with no display name renders as its wire value — legible, but it
// means AGENT_TYPE leaked into the UI.
func TestAgentType_EveryAgentHasADisplayName(t *testing.T) {
	for _, h := range AllAgentTypes() {
		if _, named := agentTypeToDisplayName[h]; !named {
			t.Errorf("%v has no display name", h)
		}
	}
	// At least one agent must back interactive sessions, or the Assistant has
	// nothing to offer.
	interactive := 0
	for _, h := range AllAgentTypes() {
		if h.SupportsInteractiveSession() {
			interactive++
		}
	}
	if interactive == 0 {
		t.Error("no agent supports interactive sessions; the Assistant picker would be empty")
	}
	// opencode is batch-only: agentbox runs it and it exits, so a session would
	// have nothing to talk to.
	if Opencode.SupportsInteractiveSession() {
		t.Error("opencode has no interactive mode in agentbox")
	}
}

// Bedrock availability lags the direct API and varies by account and region.
// The version-pinned prefixes mean a 4.6 request will NOT fall back to a 4.5
// profile — correct, but it leaves an org whose Bedrock only has 4.5 unable to
// run any Claude model there unless the catalogue carries the older version
// too. This pins that the pair coexist without shadowing each other.
func TestCatalog_GenerationsCoexistWithoutShadowing(t *testing.T) {
	pairs := []struct{ older, newer Model }{
		{ClaudeSonnet45, ClaudeSonnet46},
		{ClaudeOpus45, ClaudeOpus48},
	}
	for _, p := range pairs {
		o, n := p.older.BedrockProfilePrefix(), p.newer.BedrockProfilePrefix()
		if o == "" || n == "" {
			t.Errorf("%v/%v: both generations need a Bedrock prefix", p.older, p.newer)
			continue
		}
		// Neither may be a prefix of the other, or discovery's Contains match
		// would let one generation resolve to the other's profile — the silent
		// wrong-model swap the pinning exists to prevent.
		if strings.HasPrefix(o, n) || strings.HasPrefix(n, o) {
			t.Errorf("%q and %q overlap; discovery could resolve one to the other", o, n)
		}
		// Both must be offered by the same agents, or picking the older one
		// silently changes which harness runs.
		for _, at := range AllAgentTypes() {
			if at.Supports(p.newer) != at.Supports(p.older) {
				t.Errorf("%v supports %v but not %v; the generations must be interchangeable", at, p.newer, p.older)
			}
		}
	}
}

// Legacy models sort last in pickers, but must NOT be reordered in
// agentTypeToModels — kit's recommendedByComplexity indexes into that list, so
// putting a superseded model at the end would make "high complexity" pick it.
func TestModelsFor_LegacySortsLastWithoutDisturbingCapabilityOrder(t *testing.T) {
	all := func(Provider) bool { return true }
	got := ModelsFor(ClaudeCode, all)

	seenLegacy := false
	for _, m := range got {
		if m.IsLegacy() {
			seenLegacy = true
			continue
		}
		if seenLegacy {
			t.Errorf("current model %v appears after a legacy one in %v", m, got)
		}
	}
	// Order WITHIN each group is preserved, so the picker still reads
	// low->high inside the current models.
	var current []Model
	for _, m := range got {
		if !m.IsLegacy() {
			current = append(current, m)
		}
	}
	var expected []Model
	for _, m := range agentTypeToModels[ClaudeCode] {
		if !m.IsLegacy() {
			expected = append(expected, m)
		}
	}
	for i := range expected {
		if current[i] != expected[i] {
			t.Errorf("current models reordered: got %v, want %v", current, expected)
			break
		}
	}

	// The capability list itself must stay ascending — the last entry is what
	// a high-complexity session gets, and it must not be superseded.
	lineup := agentTypeToModels[ClaudeCode]
	if lineup[len(lineup)-1].IsLegacy() {
		t.Errorf("agentTypeToModels[ClaudeCode] ends with legacy %v; high complexity would pick it", lineup[len(lineup)-1])
	}
}

// NOTHING may depend on enum declaration order. Every ranking is an explicit
// rule, so renumbering or inserting a constant cannot silently change which
// agent runs a Task, which model a complexity resolves to, or which provider
// serves a model.
//
// The renumbering here is the test: it swaps the values the old
// order-dependent code leaned on, and every answer must be unchanged.
func TestNothingDependsOnDeclarationOrder(t *testing.T) {
	// Agent priority is a map, not the enum's order.
	if ClaudeCode.Priority() >= Opencode.Priority() {
		t.Error("claude-code must outrank opencode for a shared model")
	}
	// A model both can run resolves to the higher-priority agent regardless of
	// which constant is numerically smaller.
	agents := AgentTypesFor(ClaudeSonnet46)
	if len(agents) == 0 || agents[0] != ClaudeCode {
		t.Errorf("AgentTypesFor(Sonnet 4.6) = %v, want claude-code first", agents)
	}

	// Tier is a property, so a complexity maps to a band rather than to a list
	// index — the whole point being that today's frontier model is tomorrow's
	// balanced one, and re-tagging it is an edit rather than a reshuffle.
	if got := ModelForTier(ClaudeCode, TierFast); got != ClaudeHaiku45 {
		t.Errorf("fast tier = %v, want Haiku 4.5", got)
	}
	if got := ModelForTier(ClaudeCode, TierFrontier); got.IsLegacy() {
		t.Errorf("frontier tier = %v, a superseded model", got)
	}
	// An untiered lineup resolves to the agent's explicit default rather than
	// to whichever model is listed first.
	if got := ModelForTier(Codex, TierBalanced); got != Codex.DefaultModel() {
		t.Errorf("codex balanced = %v, want its default %v", got, Codex.DefaultModel())
	}

	// Provider preference is a rule; passing candidates in the WORST order
	// must not change the answer.
	worst := []Provider{AnthropicDirect, AWSBedrock, AnthropicSubscription}
	got, ok := PreferredProvider(worst, func(p Provider) bool {
		return p == AnthropicDirect || p == AnthropicSubscription
	})
	if !ok || got != AnthropicSubscription {
		t.Errorf("PreferredProvider = %v, want the subscription regardless of order", got)
	}
}

// AllAgentTypes feeds the agent picker, so its order is USER-VISIBLE. Ranking
// it by value would put the enum's declaration order on screen — the same
// dependency AgentTypesFor stopped carrying.
func TestAllAgentTypes_IsPriorityOrderedNotNumeric(t *testing.T) {
	got := AllAgentTypes()
	for i := 1; i < len(got); i++ {
		if got[i-1].Priority() > got[i].Priority() {
			t.Errorf("AllAgentTypes is not priority-ordered: %v", got)
		}
	}
	// And AgentTypesFor inherits it — a model several agents run must report
	// the highest-priority one first, since callers take [0] as the harness.
	for _, m := range []Model{ClaudeSonnet46, ClaudeHaiku45} {
		agents := AgentTypesFor(m)
		for i := 1; i < len(agents); i++ {
			if agents[i-1].Priority() > agents[i].Priority() {
				t.Errorf("AgentTypesFor(%v) is not priority-ordered: %v", m, agents)
			}
		}
	}
}

// The agent enum is NOT persisted — Task.AgentType stores the string
// "claude-code" — so the numbers are free to move. What is not free is a
// sentinel: a hand-bumped MaxAgentType made adding an agent read as INVALID
// until someone remembered, and a reserved value read as valid. Validity is
// membership now, exactly as Provider's is.
func TestAgentType_ValidityIsMembershipNotARange(t *testing.T) {
	for _, h := range AllAgentTypes() {
		if !h.IsValid() {
			t.Errorf("%v is declared but not valid", h)
		}
	}
	// A value past the declared set must be invalid — a range check with a
	// stale sentinel would have accepted it.
	if AgentType(99).IsValid() {
		t.Error("an undeclared agent must not be valid")
	}
	if AgentType(0).IsValid() {
		t.Error("the zero value must not be valid")
	}
	// And the string is what persists, so it must round-trip.
	for _, h := range AllAgentTypes() {
		got, err := ResolveAgentType(h.String())
		if err != nil || got != h {
			t.Errorf("%v does not round-trip through its wire string: got %v (%v)", h, got, err)
		}
	}
}

// These strings are the PERSISTED form and a container switch, in that order
// of danger:
//
//	Mongo     Task.AgentType stores "claude-code"; a rename orphans every
//	          existing Task, which then resolves to the default harness.
//	agentbox  its config switches on these exact literals to pick a driver,
//	          and it is a separate repo — a rename here compiles fine and
//	          fails at spawn.
//	wire      the AGENT_TYPE env var and the dashboard's agentType field.
//
// The round-trip test above only proves self-consistency; it would pass if
// every literal changed together. This pins the literals themselves.
func TestAgentType_WireStringsAreStable(t *testing.T) {
	pinned := map[AgentType]string{
		ClaudeCode: "claude-code",
		Codex:      "codex",
		Opencode:   "opencode",
	}
	if len(pinned) != len(agentTypeToString) {
		t.Errorf("%d agents declared, %d pinned — agentbox switches on these literals, so a new one must be pinned deliberately", len(agentTypeToString), len(pinned))
	}
	for h, want := range pinned {
		if got := h.String(); got != want {
			t.Errorf("%v.String() = %q, want %q — stored Tasks and agentbox both read this literal", h, got, want)
		}
	}
	// Empty resolves to claude-code: Tasks created before the agent type
	// existed have no value stored, and agentbox applies the same default.
	if got, err := ResolveAgentType(""); err != nil || got != ClaudeCode {
		t.Errorf(`ResolveAgentType("") = %v (%v), want claude-code — legacy Tasks store nothing`, got, err)
	}
}
