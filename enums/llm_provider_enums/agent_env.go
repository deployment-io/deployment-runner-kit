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
	// EnvSubscriptionAuthMode marks an org as using its own Claude subscription.
	// NON-SECRET: the OAuth token itself lives in the customer's own Secrets
	// Manager and is read runner-side, never transiting the control plane.
	//
	// The name still says CLAUDE because it is claude-code's own concept and the
	// value is deployed; the Go identifier no longer references the deleted
	// ClaudeAuth struct.
	EnvSubscriptionAuthMode = "CLAUDE_AUTH_MODE"

	// SubscriptionAuthModeValue is what EnvSubscriptionAuthMode is set to. Any
	// other value means "not subscription mode".
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
	return env[EnvSubscriptionAuthMode] == SubscriptionAuthModeValue
}

// IsBedrockMode reports whether an injected env bundle marks Bedrock.
func IsBedrockMode(env map[string]string) bool {
	return env[EnvBedrockMode] == BedrockModeValue
}
