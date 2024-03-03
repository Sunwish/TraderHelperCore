package api

type Notifier interface {
	Notify(Notification)
}
