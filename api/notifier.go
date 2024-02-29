package api

type Notifier interface {
	Notify(title string, message string)
}
