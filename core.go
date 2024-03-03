package main

import (
	"TraderHelperCore/api"
	"TraderHelperCore/common"
	n "TraderHelperCore/staging/notification"
	"fmt"
)

func fetchAndUpdateStockPrice(stock common.Stock) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from a panic:", r)
			// notifier.Notify("[TraderHelper] Recovered from a panic!", fmt.Sprint(r))
		}
	}()

	// fetch
	newData := ds.GetData(stock.Code)
	if newData.LastPrice == 0 {
		title := fmt.Sprintf("[%s] 数据获取失败", stock.Code)
		content := "请检查网络连接"
		notification := n.MakeNotification(title, content, n.WithExtType(api.NotificationTypeNetwork))
		notifier.Notify(*notification)
		return
	}

	// 判断上破下破，配置通知内容并记录激活状态
	notifyTitle := ""
	notifyContent := ""
	stockConfig := favoriteStocks[stock.Code]
	breakDirection := 0
	if stockConfig.BreakUp > 0 && newData.LastPrice >= stockConfig.BreakUp {
		activeStocksMutex.Lock()
		activeStocks[stock.Code] = true
		activeStocksMutex.Unlock()
		breakDirection = 1
		notifyTitle = fmt.Sprintf("[%s] %s 触发上破", newData.Code, newData.Name)
		notifyContent = fmt.Sprintf("现价：%f，上破 %f", newData.LastPrice, stockConfig.BreakUp)
	}
	if stockConfig.BreakDown > 0 && newData.LastPrice <= stockConfig.BreakDown {
		activeStocksMutex.Lock()
		activeStocks[stock.Code] = true
		activeStocksMutex.Unlock()
		notifyTitle = fmt.Sprintf("[%s] %s 触发下破", newData.Code, newData.Name)
		notifyContent = fmt.Sprintf("现价：%f，下破 %f", newData.LastPrice, stockConfig.BreakDown)
	}

	// update
	stocksDataMutex.Lock()
	stocksData[stock.Code] = newData
	stocksDataMutex.Unlock()

	// notify
	if notifyTitle != "" || notifyContent != "" {
		notification := n.MakeNotification(notifyTitle, notifyContent, n.WithExtType(api.NotificationTypeBreak), n.WithExtCode(breakDirection))
		notifier.Notify(*notification)
	}
}

func isCodeValid(code string) bool {
	data := ds.GetData(code)
	return data.LastPrice != 0
}
