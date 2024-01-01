package notification_enums

type Type uint

const (
	BuildStarted            Type = 1
	PreviewStarted          Type = 2
	PreviewCompletedSuccess Type = 3
	PreviewCompletedError   Type = 4
	PreviewDeletedSuccess   Type = 5
	PreviewDeletionStarted  Type = 6
	PreviewDeletionUndone   Type = 7
	BuildCreated            Type = 8
	BuildCompletedSuccess   Type = 9
	BuildCompletedError     Type = 10
	HeadersUpdatedSuccess   Type = 11
	HeadersUpdatedError     Type = 12
)
