package llm_provider_enums

// The env-var contract between the control plane and the agent container.
//
// These names and values are a CROSS-REPO WIRE CONTRACT. deployment-server
// puts them into a Job's AgentEnvVars; deployment-runner reads them to decide
// what to inject; agentbox reads some of them directly. They live here because
// deployment-runner-kit is the only module all three can import — the runner
// cannot import kit, which is why it previously carried its own copies with a
// comment reading "change them in both or subscription auth silently stops
// engaging and every task quietly falls back to the API key". That comment
// described a bug waiting to happen; this removes the possibility.
//
// ⚠️ THE VALUES ARE THE CONTRACT, not the identifiers. Renaming a Go constant
// is free and compile-checked. Changing a string here changes what a deployed
// runner is looking for, so a mismatched deploy fails silently: the marker is
// simply never recognised and the task falls back rather than erroring.
const (
	// EnvAgentAuthMode marks an org as authenticating with its own subscription
	// rather than an API key. NON-SECRET: the OAuth token itself lives in the
	// customer's own Secrets Manager and is read runner-side, never transiting
	// the control plane.
	//
	// AGENT-NEUTRAL BY DESIGN. Unlike EnvBedrockMode below, this marker is
	// purely ours — the runner consumes it and deletes it, so no CLI ever sees
	// it and the name was never constrained. It composes with AGENT_TYPE:
	// AGENT_TYPE=claude-code + AGENT_AUTH_MODE=subscription today, and
	// AGENT_TYPE=codex + AGENT_AUTH_MODE=subscription when a Codex
	// subscription lands — no second marker, no second code path deciding
	// which one to read.
	EnvAgentAuthMode = "AGENT_AUTH_MODE"

	// EnvAgentAuthModeLegacy is the name this marker had while it was
	// Claude-specific.
	//
	// TRANSITIONAL — remove one release after every runner is upgraded.
	// Readers must accept BOTH: a control plane sending the new name to a
	// runner that only knows the old one would silently stop engaging
	// subscription auth and quietly fall back to the API key — metered billing,
	// no error, tasks still passing. Accepting both costs one map lookup and
	// removes the coordinated-deploy requirement entirely.
	EnvAgentAuthModeLegacy = "CLAUDE_AUTH_MODE"

	// SubscriptionAuthModeValue is what EnvAgentAuthMode is set to. Any other
	// value means "not subscription mode" — the comparison is exact, so a
	// near-miss must not engage subscription auth.
	SubscriptionAuthModeValue = "subscription"

	// EnvBedrockMode marks an org as reaching models through AWS Bedrock, and
	// tells the runner to assume dr-bedrock-role and inject short-lived
	// credentials at spawn.
	//
	// It does double duty: claude-code reads this exact variable as its own
	// Bedrock switch, which is why the runner keeps it for claude-code and
	// strips it for every other agent — the name is misleading for opencode,
	// which cannot use it.
	EnvBedrockMode = "CLAUDE_CODE_USE_BEDROCK"

	// BedrockModeValue is what EnvBedrockMode is set to.
	BedrockModeValue = "1"
)

// IsSubscriptionAuthMode reports whether an injected env bundle marks
// subscription auth. Provided so callers compare through one place rather than
// each writing their own string comparison against the value above.
func IsSubscriptionAuthMode(env map[string]string) bool {
	if env[EnvAgentAuthMode] == SubscriptionAuthModeValue {
		return true
	}
	// Legacy name, for jobs created before the rename or picked up by a runner
	// mid-upgrade. Drop with EnvAgentAuthModeLegacy.
	return env[EnvAgentAuthModeLegacy] == SubscriptionAuthModeValue
}

// ConsumeSubscriptionAuthMode reports whether env marks subscription auth and
// removes the marker, so it never reaches the agent container. It is our
// control-plane signal, not a CLI's — the runner acts on it by injecting the
// OAuth token, and the CLI has no use for it.
//
// Detect-and-strip is ONE call on purpose. Two accepted names means a
// hand-written `delete` eventually removes one and leaks the other, and a
// leaked marker is invisible: the task runs fine, the variable is just sitting
// in the container. Callers that only want to test the bundle without mutating
// it should use IsSubscriptionAuthMode.
func ConsumeSubscriptionAuthMode(env map[string]string) bool {
	on := IsSubscriptionAuthMode(env)
	// Unconditional: a marker set to some other value is still ours, and still
	// must not reach the container.
	delete(env, EnvAgentAuthMode)
	delete(env, EnvAgentAuthModeLegacy)
	return on
}

// IsBedrockMode reports whether an injected env bundle marks Bedrock.
func IsBedrockMode(env map[string]string) bool {
	return env[EnvBedrockMode] == BedrockModeValue
}
