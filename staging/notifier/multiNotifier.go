package notifier

import "TraderHelperCore/api"

type NotifierType int

const (
	NotifierTypeLog NotifierType = iota
	NotifierTypePushDeer
)

type notifierConfig struct {
	Notifier api.Notifier
	Enable   bool
}

type multiNotifier struct {
	notifiers map[NotifierType]notifierConfig
}

func NewMultiNotifier(notifiers ...api.Notifier) *multiNotifier {
	m := &multiNotifier{
		notifiers: make(map[NotifierType]notifierConfig),
	}
	for _, n := range notifiers {
		m.AddNotifier(n)
	}
	return m
}

func (m multiNotifier) Notify(title string, message string) {
	for _, notifier := range m.notifiers {
		if notifier.Enable {
			notifier.Notifier.Notify(title, message)
		}
	}
}

func getNotifierType(n api.Notifier) NotifierType {
	switch n.(type) {
	case *pushdeerNotifier:
		return NotifierTypePushDeer
	case *logNotifier:
		return NotifierTypeLog
	default:
		panic("unknown notifier type")
	}
}

// 添加通知器，每种类型的通知器只允许添加一个，新添加的会覆盖旧添加的
func (m *multiNotifier) AddNotifier(notifier api.Notifier) {
	t := getNotifierType(notifier)
	m.notifiers[t] = notifierConfig{Notifier: notifier, Enable: true}
}

func (m *multiNotifier) RemoveNotifier(notifier api.Notifier) {
	t := getNotifierType(notifier)
	delete(m.notifiers, t)
}

func (m *multiNotifier) EnableNotifier(notifierType NotifierType) {
	if n, ok := m.notifiers[notifierType]; ok {
		n.Enable = true
		m.notifiers[notifierType] = n
	}
}

func (m *multiNotifier) DisableNotifier(notifierType NotifierType) {
	if n, ok := m.notifiers[notifierType]; ok {
		n.Enable = false
		m.notifiers[notifierType] = n
	}
}
