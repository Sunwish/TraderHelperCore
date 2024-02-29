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

func (logNotifier) Notify(title string, content string) {
	fmt.Printf("title: %s\nmessage: %s\n", title, content)
}
