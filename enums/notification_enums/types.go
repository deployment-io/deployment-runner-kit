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
)
