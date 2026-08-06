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
	// EnvClaudeCodeUseBedrock is claude-code's OWN switch for routing through
	// Bedrock. Its name and value are Anthropic's, not ours — which is why this
	// one legitimately lives in the container env, and why the identifier names
	// the CLI rather than the concept. There is no general "Bedrock mode" env
	// var to abstract over: codex has no equivalent, and opencode selects
	// Bedrock through its model id. A neutral name here would invite handing a
	// second agent a variable it does not understand.
	//
	// The runner WRITES it (see ApplyClaudeCodeUseBedrock); nothing on our side
	// reads it back. Derive from Provider instead.
	EnvClaudeCodeUseBedrock = "CLAUDE_CODE_USE_BEDROCK"

	// ClaudeCodeUseBedrockValue is what EnvClaudeCodeUseBedrock is set to.
	// Meaningless on its own — it exists only to name that variable's value,
	// which is what binds it to claude-code too.
	ClaudeCodeUseBedrockValue = "1"
)

// ApplyClaudeCodeUseBedrock sets claude-code's Bedrock switch when — and only
// when — the org's provider routes through Bedrock AND the agent is
// claude-code. No other agent understands this variable: codex ignores it, and
// opencode selects Bedrock through its own "amazon-bedrock/…" model id (see
// OpencodeModelID).
//
// It keeps taking agentType despite being claude-code-specific, so callers can
// call it unconditionally and let it refuse. That is the point: it replaces a
// strip-if-wrong-agent step, and removing a variable that should not be there
// fails open — miss one path and it leaks into the container. Pushing the
// agent check back to the caller would reintroduce exactly that.
func ApplyClaudeCodeUseBedrock(env map[string]string, p Provider, agentType AgentType) {
	if p == AWSBedrock && agentType == ClaudeCode {
		env[EnvClaudeCodeUseBedrock] = ClaudeCodeUseBedrockValue
	}
}
