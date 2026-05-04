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
// it sends one value (then closes) when the server reports the Job
// has been moved to the Stopping state. Implementations must not
// close or write to the channel.
type StoppableCommand interface {
	Command
	SetStopSignal(stop <-chan struct{})
}
