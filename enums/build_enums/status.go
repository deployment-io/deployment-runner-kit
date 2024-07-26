package build_enums

type Status uint

const (
	Pending  Status = 1
	Success  Status = 2
	Error    Status = 3
	Running  Status = 4
	Stopping Status = 5
	Stopped  Status = 6
	TimedOut Status = 7
)

var typeToString = map[Status]string{
	Pending:  "Deployment Pending",
	Success:  "Deployment Successful",
	Error:    "Error",
	Running:  "Deploying",
	Stopping: "Stopping Deployment",
	Stopped:  "Deployment Stopped",
	TimedOut: "Deployment Timed out",
}

func (t Status) String() string {
	return typeToString[t]
}
