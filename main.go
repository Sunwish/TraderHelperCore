package main

import (
	"TraderHelperCore/common"
	"TraderHelperCore/staging/dataSource"
	notifiers "TraderHelperCore/staging/notifier"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var favoriteStocks = make(map[string]common.Stock)           // 用于存储自选股列表
var stocksData = make(map[string]common.StockData)           // 用于存储自选股实时数据
var ds = dataSource.NewDataSource(dataSource.SOURCE_TENCENT) // 数据源
var tickerDuration = 3 * time.Second
var notifier = notifiers.NewLogNotifier()

func main() {
	// 初始化ticker
	ticker := time.NewTicker(tickerDuration)

	// 创建一个goroutine处理定时任务
	go func() {
		for range ticker.C {
			for _, stock := range favoriteStocks {
				go fetchAndUpdateStockPrice(stock)
			}
		}
	}()

	// 设置HTTP路由和启动服务器
	mux := http.NewServeMux()
	// ...设置路由处理函数...
	mux.HandleFunc("/add_favorite_stock", addFavoriteStock)
	mux.HandleFunc("/update_break_price", updateBreakPrice)
	mux.HandleFunc("/get_favorite_stocks", getFavoriteStocks)
	mux.HandleFunc("/get_favorite_stocks_data", getFavoriteStocksData)
	mux.HandleFunc("/remove_favorite_stock", removeFavoriteStock)
	mux.HandleFunc("/test/force_fetch", test_forceFetch)

	server := &http.Server{
		Addr:    ":9888",
		Handler: mux,
	}

	// 启动HTTP服务器
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	fmt.Println("Server is running...")

	// 确保程序在接收到信号时优雅退出
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// 停止ticker
	ticker.Stop()

	// 关闭HTTP服务器
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	fmt.Println("Server stopped.")
}

func fetchAndUpdateStockPrice(stock common.Stock) {
	// fetch
	newData := ds.GetData(stock.Code)
	// update
	stocksData[stock.Code] = newData
	// 判断上破下破
	stockConfig := favoriteStocks[stock.Code]
	if stockConfig.BreakUp > 0 && newData.LastPrice >= stockConfig.BreakUp {
		notifier.Notify(fmt.Sprintf("[%s] %s 触发上破", newData.Code, newData.Name), fmt.Sprintf("现价：%f，上破 %f", newData.LastPrice, stockConfig.BreakUp))
	}
	if stockConfig.BreakDown > 0 && newData.LastPrice <= stockConfig.BreakDown {
		notifier.Notify(fmt.Sprintf("[%s] %s 触发下破", newData.Code, newData.Name), fmt.Sprintf("现价：%f，下破 %f", newData.LastPrice, stockConfig.BreakDown))
	}
}
