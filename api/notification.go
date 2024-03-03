package api

type Notification struct {
	Title   string
	Content string
	Ext     NotificationExt
}

type NotificationExt struct {
	Type  NotificationType
	Level int
	Code  int
	Time  int64
}

type NotificationType int

const (
	NotificationTypeLog NotificationType = iota
	NotificationTypeInfo
	NotificationTypeSystem
	NotificationTypeBreak
	NotificationTypeApiResponse
	NotificationTypeNetwork
)

type NotificationCode int

const (
	// NotificationType: System
	NotificationCodeServerStartupFail NotificationCode = iota
)
