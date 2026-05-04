package jobs

import "io"

type Runner interface {
	Run(parameters map[string]interface{}, logsWriter io.Writer) (map[string]interface{}, error)
}

type Command interface {
	Runner
}

// StoppableCommand is implemented by commands that can be cancelled
// mid-run via an external signal. The runner's outer execution loop
// type-asserts each Command and, if it satisfies StoppableCommand,
// calls SetStopSignal *before* invoking Run. Long-running commands
// (e.g., the Tasks agent-step container spawn) honor the channel by
// killing their in-flight subprocess; short-lived commands don't need
// to implement.
//
// SetStopSignal is invoked exactly once per command instance. The
// supplied channel is shared with the runner's heartbeat poller —
// it CLOSES (no value sent) when the server reports the Job has been
// moved to the Stopping state. Closed-channel-as-broadcast wakes ALL
// readers atomically and idempotently; implementations must not
// close or write to the channel.
type StoppableCommand interface {
	Command
	SetStopSignal(stop <-chan struct{})
}

// ProgressEmittingCommand is implemented by commands that produce
// incremental progress (live turn / token counters from a long-running
// agent subprocess, today). The runner's outer execution loop
// type-asserts each Command and, if it satisfies ProgressEmitting-
// Command, calls SetProgressSink before invoking Run with a callback
// the command invokes whenever it has new progress. The runner caches
// the latest snapshot and forwards it on the next heartbeat call.
//
// SetProgressSink is invoked exactly once per command instance, before
// Run. Implementations may call sink any number of times (including
// zero). sink is safe to call from any goroutine — its implementation
// stores into an atomic pointer that the heartbeat poller reads.
//
// Distinct interface from StoppableCommand because not every long-
// running command emits structured progress (e.g., a slow build
// container has no equivalent of "current turn / token usage"), and
// not every progress-emitting command necessarily wants stop wiring
// (though in practice both apply to the same commands today). Splitting
// keeps the contracts orthogonal and the type-assertions cheap.
type ProgressEmittingCommand interface {
	Command
	SetProgressSink(sink func(LiveProgressV1))
}
