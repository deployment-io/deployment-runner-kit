package jobs

type Logger interface {
	Log(messages []string) error
}
