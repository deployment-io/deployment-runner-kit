package llm_provider_enums

// The env-var contract for spawning an agent container.
//
// ⚠️ ONLY TWO KINDS OF THING BELONG IN AN AGENT'S ENV:
//
//  1. secrets and settings the agent itself consumes (ANTHROPIC_API_KEY,
//     AWS_* credentials);
//  2. a CLI's own documented switches, whose names we do not choose.
//
// A decision *we* make about how to run the agent is NEITHER. It travels as a
// typed job parameter — parameters_enums.AgentProvider, the sibling of
// AgentType — because that is a control signal for the runner, not payload for
// the container.
//
// This distinction used to be blurred: the provider was encoded as a
// CLAUDE_AUTH_MODE string, stuffed into the same bundle as the API key, and
// then read back and deleted by the runner. That round trip threw away the
// type, put a non-secret in the secrets channel, and created a contract whose
// breakage was silent — a name mismatch meant subscription auth simply never
// engaged and every task fell back to metered API-key billing, with no error
// and no failing task. Reading Provider directly removes the marker, the round
// trip, and that whole failure mode.
const (
	// EnvBedrockMode is claude-code's OWN switch for routing through Bedrock —
	// its name and value are Anthropic's, not ours, which is exactly why this
	// one legitimately lives in the container env.
	//
	// The runner WRITES it (see ApplyBedrockMode); nothing on our side reads it
	// back. Derive from Provider instead.
	EnvBedrockMode = "CLAUDE_CODE_USE_BEDROCK"

	// BedrockModeValue is what EnvBedrockMode is set to.
	BedrockModeValue = "1"
)

// ApplyBedrockMode sets claude-code's Bedrock switch when — and only when — the
// org's provider routes through Bedrock AND the agent is claude-code. No other
// agent understands this variable: codex ignores it, and opencode selects
// Bedrock through its own "amazon-bedrock/…" model id (see OpencodeModelID).
//
// Write-if-needed, deliberately, rather than the strip-if-wrong-agent it
// replaces. Removing a variable that should not be there fails open — miss one
// path and it leaks into the container; adding one only where it belongs cannot
// leak at all.
func ApplyBedrockMode(env map[string]string, p Provider, agentType AgentType) {
	if p == AWSBedrock && agentType == ClaudeCode {
		env[EnvBedrockMode] = BedrockModeValue
	}
}
