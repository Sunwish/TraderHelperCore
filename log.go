package main

import (
	"TraderHelperCore/api"
	n "TraderHelperCore/staging/notification"
	"fmt"
	"os"
)

// 通过notifier分发日志

func Logln(a ...any) {
	if notifier == nil {
		return
	}
	notifier.Notify(*n.MakeNotification(fmt.Sprintln(a...), "", n.WithExtType(api.NotificationTypeLog)))
}

func Logf(format string, a ...any) {
	if notifier == nil {
		return
	}
	notifier.Notify(*n.MakeNotification(fmt.Sprintf(format, a...), "", n.WithExtType(api.NotificationTypeLog)))
}

func LogFatal(v ...any) {
	if notifier == nil {
		return
	}
	notifier.Notify(*n.MakeNotification(fmt.Sprint(v...), "", n.WithExtType(api.NotificationTypeSystem)))
	os.Exit(1)
}
