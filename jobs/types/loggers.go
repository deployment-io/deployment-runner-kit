package types

type Logger interface {
	Log(messages []string) error
}
