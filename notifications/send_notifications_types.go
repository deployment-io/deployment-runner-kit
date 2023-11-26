package notifications

import (
	"github.com/deployment-io/deployment-runner-kit/enums/notification_enums"
	"github.com/deployment-io/deployment-runner-kit/types"
)

type SendNotificationDtoV1 struct {
	PreviewID    string
	DeploymentID string
	BuildID      string
	Type         notification_enums.Type
}

type SendNotificationsArgsV1 struct {
	types.AuthArgsV1
	Notifications []SendNotificationDtoV1
}

type SendNotificationsReplyV1 struct {
	Done bool
}
