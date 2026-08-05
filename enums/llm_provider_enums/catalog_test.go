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
	h, err := ResolveHarness("")
	if err != nil || h != ClaudeCode {
		t.Errorf("ResolveHarness(\"\") = %v, %v; want ClaudeCode, nil", h, err)
	}
	if _, err := ResolveHarness("nonexistent"); err == nil {
		t.Error("ResolveHarness must reject an unknown harness rather than defaulting it")
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
	for h := ClaudeCode; h < MaxHarness; h++ {
		for _, m := range h.Models() {
			if got := ProvidersFor(h, m); len(got) == 0 {
				t.Errorf("%s lists model %s but no provider can serve it — harnessToProviders and modelToProviders disagree", h, m)
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
	for h := ClaudeCode; h < MaxHarness; h++ {
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
	harnesses := HarnessesFor(NovaProV1)
	if len(harnesses) != 1 || harnesses[0] != Opencode {
		t.Errorf("HarnessesFor(NovaProV1) = %v, want exactly [Opencode]", harnesses)
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
		for _, h := range HarnessesFor(m) {
			if !h.Supports(m) {
				t.Errorf("HarnessesFor(%s) returned %s, which does not support it", m, h)
			}
		}
	}
	if len(HarnessesFor(Gpt55)) < 2 {
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
	// unsupported. harnessToProviders already excludes it; returning "" here is
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
	// If harnessToProviders says opencode can reach a provider, that provider
	// must have an opencode name — otherwise the catalogue offers a combination
	// no id can be rendered for, and the task fails at spawn with a malformed
	// model.
	for _, p := range Opencode.Providers() {
		if p.OpencodeName() == "" {
			t.Errorf("opencode can reach %s but it has no opencode provider name", p)
		}
	}
}
