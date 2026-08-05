// Package llm_provider_enums is the single source of truth for the three
// axes of running a coding agent: which PROVIDER serves a model, which MODEL
// is being run, and which agent HARNESS runs it — plus the relationships
// between them.
//
// It lives in deployment-runner-kit because that is the only module every
// repo imports (kit, app-server, deployment-server and deployment-runner all
// depend on it, while deployment-runner cannot import kit). Anything defined
// here is visible to the control plane and the runner alike, which is what
// stops the two drifting — the failure mode this package exists to end. See
// the duplicated CLAUDE_AUTH_MODE / CLAUDE_CODE_USE_BEDROCK literals in
// deployment-runner for what the alternative looks like.
package llm_provider_enums

// Provider identifies WHOSE models are being called and over which API.
//
// ⚠️ PERSISTED — DO NOT RENUMBER. Values 1-4 are stored in live MongoDB
// documents at Organization.ClaudeAuth.Provider. This enum was moved here
// from kit/enums/claude_auth_enums with its numbering intact precisely so no
// data migration was needed; reordering for tidiness would silently
// reinterpret every existing org's credential configuration. Append only,
// before MaxProvider.
//
// ⚠️ NOT EVERY VALUE IS PERSISTED. 1-4 appear in ClaudeAuth.Provider.
// OpenAIDirect does not appear anywhere: OpenAIAuth has no Provider field and
// infers its provider from the presence of an API key, so nothing writes 5.
// It exists for the catalogue below, which needs to reason about OpenAI even
// though the org document never names it. Do not read "no org has provider 5"
// as a bug. PLAN_provider_centric_llm_keys.md §3 makes this symmetric by
// replacing both structs with LLMConfig.Providers[].
//
// ⚠️ AnthropicSubscription IS AN AUTH MODE WEARING A PROVIDER'S CLOTHES.
// It is Anthropic-the-vendor reached with an OAuth subscription token rather
// than an API key, so it is not a sibling of AnthropicDirect in the way
// AWSBedrock is. Keeping the two fused is a deliberate, recorded decision
// (see the plan §2 vs §11.3): splitting auth mode into its own axis would
// have forced a migration of the persisted values above for a distinction
// only the settings UI actually cares about, and the UI can render one
// "Anthropic" card with a mode toggle regardless. The rename from the former
// "Subscription" at least makes the fusion visible instead of hidden.
type Provider uint

const (
	AnthropicDirect Provider = iota + 1 // 1 — Anthropic API key
	AWSBedrock                          // 2 — runner-assumed IAM role, no stored secret
	GoogleVertex                        // 3 — reserved, not wired
	// AnthropicSubscription is the customer's own Claude Code subscription
	// (Pro/Max) OAuth token. No credential is held control-plane-side: the
	// token lives only in the customer's own AWS Secrets Manager and is read
	// by the runner at agentbox spawn. claude-code only — the runner refuses
	// it for other harnesses, which agentTypeToProviders below encodes.
	AnthropicSubscription // 4
	// OpenAIDirect is catalogue-only; see the note above.
	OpenAIDirect // 5

	MaxProvider // always add providers before MaxProvider
)

// providerToString values are USER-VISIBLE. They are surfaced as
// ProviderName on GET /organizations/current/claude-auth and rendered in the
// dashboard, so they are deliberately unchanged by the Go-level rename of
// Subscription -> AnthropicSubscription. Renaming an identifier should not
// silently alter what a customer reads.
var providerToString = map[Provider]string{
	AnthropicDirect:       "Anthropic Direct",
	AWSBedrock:            "AWS Bedrock",
	GoogleVertex:          "Google Vertex",
	AnthropicSubscription: "Claude Subscription",
	OpenAIDirect:          "OpenAI Direct",
}

func (p Provider) String() string {
	return providerToString[p]
}

func (p Provider) IsZero() bool {
	return p == 0
}

// IsValid reports whether p names a declared provider. The zero value means
// "unconfigured" and is not valid.
func (p Provider) IsValid() bool {
	return p > 0 && p < MaxProvider
}

// AuthMode returns HOW this provider is authenticated to.
//
// This recovers the axis PLAN_provider_centric_llm_keys.md §2 wanted split out
// of Provider — without renumbering anything. Auth mode is DERIVED from the
// provider rather than stored beside it, so the persisted values above stay
// untouched while callers that genuinely need the distinction (which fields a
// settings card renders, whether a credential is stored at all) can still ask.
//
// Deliberately NOT a stored field: a provider fully determines its auth mode,
// so persisting both would create a pair that can disagree.
func (p Provider) AuthMode() AuthMode {
	switch p {
	case AnthropicDirect, OpenAIDirect:
		return AuthAPIKey
	case AWSBedrock, GoogleVertex:
		return AuthCloudRole
	case AnthropicSubscription:
		return AuthSubscription
	}
	return AuthUnknown
}

// AuthMode is how a provider is authenticated to — the axis fused into
// Provider for persistence reasons, recoverable via Provider.AuthMode().
type AuthMode uint

const (
	AuthUnknown AuthMode = iota
	// AuthAPIKey — a secret held control-plane-side, encrypted at rest.
	AuthAPIKey
	// AuthCloudRole — no stored secret at all; the runner assumes a role in the
	// customer's own cloud and vends short-lived credentials at spawn.
	AuthCloudRole
	// AuthSubscription — an OAuth token held in the CUSTOMER's secret store and
	// read runner-side, never transiting the control plane.
	AuthSubscription

	MaxAuthMode
)

var authModeToString = map[AuthMode]string{
	AuthAPIKey:       "API key",
	AuthCloudRole:    "Cloud IAM role",
	AuthSubscription: "Subscription",
}

func (a AuthMode) String() string {
	return authModeToString[a]
}

// StoresSecret reports whether configuring this provider means holding a
// credential control-plane-side. False for Bedrock/Vertex (role-assumed) and
// for subscriptions (token lives in the customer's own secret store), which is
// why "configured" cannot be a key-presence check for those.
func (p Provider) StoresSecret() bool {
	return p.AuthMode() == AuthAPIKey
}
