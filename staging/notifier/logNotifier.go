package notifier

import "fmt"

type logNotifier struct {
}

func NewLogNotifier() *logNotifier {
	return &logNotifier{}
}

func (logNotifier) Notify(title string, content string) {
	fmt.Printf("title: %s\nmessage: %s\n", title, content)
}
