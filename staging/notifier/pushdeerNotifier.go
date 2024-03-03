package notifier

import (
	"TraderHelperCore/api"
	"fmt"
	"net/http"
	"net/url"
)

type pushdeerNotifier struct {
	baseUrl string
	appKey  string
}

func NewPushdeerNotifier(baseUrl string, appKey string) api.Notifier {
	return &pushdeerNotifier{
		baseUrl: baseUrl,
		appKey:  appKey,
	}
}

func (p pushdeerNotifier) Notify(notification api.Notification) {
	fullUrl := p.baseUrl + "push?pushkey=" + p.appKey + "&text=" + url.QueryEscape(notification.Title) + "&desp=" + url.QueryEscape(notification.Content) + "&type=markdown"
	response, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("消息推送失败：", err)
		return
	}
	defer response.Body.Close()
}
