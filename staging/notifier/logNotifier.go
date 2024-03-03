package notifier

import (
	"TraderHelperCore/api"
	"fmt"
)

type logNotifier struct {
}

func NewLogNotifier() api.Notifier {
	return &logNotifier{}
}

func (logNotifier) Notify(notification api.Notification) {
	if notification.Ext.Type == api.NotificationTypeLog {
		fmt.Print(notification.Title)
	} else {
		fmt.Printf("title: %s\nmessage: %s\n", notification.Title, notification.Content)
	}
}
