package jobs

import "github.com/deployment-io/deployment-runner-kit/types"

type UpsertJobHeartbeatArgsV1 struct {
	types.AuthArgsV1
	JobID string
	// LiveProgress carries the latest in-flight snapshot from the
	// command being executed (e.g., turn count + token usage from
	// agentbox during a Tasks Step Job). Optional — nil for command
	// types that don't emit progress, or for the first few heartbeats
	// of a Tasks Step Job before the agent's parser has produced
	// anything to snapshot. See jobs.ProgressEmittingCommand for the
	// runner-side opt-in interface that drives population of this
	// field.
	LiveProgress *LiveProgressV1
}

type UpsertJobHeartbeatReplyV1 struct {
	Stopping bool
}

// LiveProgressV1 is the on-the-wire shape for a live progress
// snapshot. Used by:
//   - agentbox in its progress.json file
//   - runner reading progress.json and forwarding via heartbeat
//   - server storing on Job.LiveProgress
//   - dashboard rendering live counters
//
// Flat (no nested TokenUsage) so mongo persistence and JSON
// serialization stay trivial. UpdatedAtUnix is stamped by the
// producer (agentbox) so consumers can detect stale snapshots if
// the producer goes silent without exiting.
type LiveProgressV1 struct {
	Turns           int   `json:"turns" bson:"turns"`
	InputTokens     int   `json:"input_tokens" bson:"inputTokens"`
	OutputTokens    int   `json:"output_tokens" bson:"outputTokens"`
	CacheReadTokens int   `json:"cache_read_tokens" bson:"cacheReadTokens"`
	UpdatedAtUnix   int64 `json:"updated_at_unix" bson:"updatedAtUnix"`
}
