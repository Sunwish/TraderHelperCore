package notification

import (
	"TraderHelperCore/api"
	"time"
)

type ExtChanger func(notification *api.Notification)

func MakeNotification(title string, content string, extChangers ...ExtChanger) *api.Notification {
	notification := &api.Notification{
		Title:   title,
		Content: content,
		Ext: api.NotificationExt{
			Code:  0,
			Time:  time.Now().Unix(),
			Type:  api.NotificationTypeInfo,
			Level: 0,
		},
	}
	for _, changer := range extChangers {
		changer(notification)
	}
	return notification
}

func WithExtCode(code int) ExtChanger {
	return func(notification *api.Notification) {
		notification.Ext.Code = code
	}
}

func WithExtTime(time int64) ExtChanger {
	return func(notification *api.Notification) {
		notification.Ext.Time = time
	}
}

func WithExtType(type_ api.NotificationType) ExtChanger {
	return func(notification *api.Notification) {
		notification.Ext.Type = type_
	}
}

func WithExtLevel(level int) ExtChanger {
	return func(notification *api.Notification) {
		notification.Ext.Level = level
	}
}
