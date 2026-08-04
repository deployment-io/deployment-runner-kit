package llm_provider_enums

import "testing"

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
		NovaPro:        "nova-pro",
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

func TestCatalog_BedrockModelsAllHaveAFamilyToken(t *testing.T) {
	// modelToProviders claiming AWSBedrock without a family token in
	// modelToBedrockFamily means discovery has nothing to match on, and the
	// logical id reaches Bedrock verbatim — the exact failure the live smoke
	// test produced ("The provided model identifier is invalid").
	for m := ClaudeHaiku45; m < MaxModel; m++ {
		onBedrock := false
		for _, p := range m.Providers() {
			if p == AWSBedrock {
				onBedrock = true
			}
		}
		if onBedrock && m.BedrockFamily() == "" {
			t.Errorf("model %s lists AWSBedrock as a provider but has no Bedrock family token; discovery cannot resolve it", m)
		}
		if !onBedrock && m.BedrockFamily() != "" {
			t.Errorf("model %s has a Bedrock family token but does not list AWSBedrock as a provider", m)
		}
	}
}

func TestCatalog_BedrockFamiliesAreFamiliesNotConcreteIDs(t *testing.T) {
	// A version or revision baked in here would need a release of this module
	// plus a bump in four repos per Bedrock model launch, and would fail at task
	// time in between. Concrete ids carry dots and colons; families must not.
	for m := ClaudeHaiku45; m < MaxModel; m++ {
		f := m.BedrockFamily()
		for _, bad := range []string{".", ":"} {
			if f != "" && contains(f, bad) {
				t.Errorf("Bedrock family %q for %s looks like a concrete profile id; it must be a family token only", f, m)
			}
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
	got := NovaPro.Providers()
	if len(got) != 1 || got[0] != AWSBedrock {
		t.Errorf("NovaPro.Providers() = %v, want exactly [AWSBedrock]", got)
	}
	if NovaPro.Vendor() != VendorAmazon {
		t.Errorf("NovaPro.Vendor() = %v, want VendorAmazon", NovaPro.Vendor())
	}

	// Only opencode can run it — claude-code speaks the Anthropic API and
	// codex is OpenAI-only, so neither can drive a Nova model however it is
	// reached.
	harnesses := HarnessesFor(NovaPro)
	if len(harnesses) != 1 || harnesses[0] != Opencode {
		t.Errorf("HarnessesFor(NovaPro) = %v, want exactly [Opencode]", harnesses)
	}
	if ClaudeCode.Supports(NovaPro) || Codex.Supports(NovaPro) {
		t.Error("neither claude-code nor codex can run a Nova model")
	}

	// And it must be resolvable through opencode, or offering it is a lie.
	if len(ProvidersFor(Opencode, NovaPro)) == 0 {
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
