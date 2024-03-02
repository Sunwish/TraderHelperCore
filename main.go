package main

import (
	"TraderHelperCore/common"
	"TraderHelperCore/staging/dataSource"
	notifiers "TraderHelperCore/staging/notifier"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"
)

var (
	pushdeerBaseUrl = flag.String("pushdeerBaseUrl", "", "Pushdeer Notify Base URL")
	pushdeerKey     = flag.String("pushdeerKey", "", "Pushdeer Notify Key")
	dataDirectory   = flag.String("dataDirectory", "./data", "Data directory path")
	checkInterval   = flag.Int("checkInterval", 2, "Interval of checking latest stock data")
	port            = flag.Int("port", 9888, "Port of the server")
)

var dataFileName = "favoriteStocks.json"
var favoriteStocks = make(map[string]common.Stock) // 用于存储自选股列表
var stocksData = make(map[string]common.StockData) // 用于存储自选股实时数据
var stocksDataMutex = sync.RWMutex{}
var activeStocks = make(map[string]bool) // 用于存储触发了预警但未确认的股票列表
var activeStocksMutex = sync.RWMutex{}
var ds = dataSource.NewDataSource(dataSource.SOURCE_TENCENT) // 数据源
var notifier = notifiers.NewMultiNotifier(notifiers.NewLogNotifier())
var tickerDuration time.Duration

func main() {
	// 异常处理
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from a panic:", r)
			notifier.Notify("[TraderHelper] Recovered from a panic!", fmt.Sprint(r))
		}
	}()

	// 解析启动参数
	fmt.Println("Parsing flags...")
	flag.Parse()
	tickerDuration = time.Duration(*checkInterval) * time.Second

	// 配置通知
	configNotifier()

	// 载入自选股列表
	favoriteStocks, _ = common.LoadFavoriteStocksFromFile(path.Join(*dataDirectory, dataFileName))

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

	// 设置路由和处理函数
	mux := http.NewServeMux()
	configRoute(mux)
	// 跨域处理
	corsHandler := corsWrapper(mux)
	// 构建服务
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: corsHandler,
	}
	// 启动服务
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
	fmt.Printf("Server is running at :%d\n", *port)

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

func configNotifier() {
	if pushdeerBaseUrl != nil && *pushdeerBaseUrl != "" && pushdeerKey != nil && *pushdeerKey != "" {
		notifier.AddNotifier(notifiers.NewPushdeerNotifier(*pushdeerBaseUrl, *pushdeerKey))
		notifier.Notify("[TraderHelper] Pushdeer 预警成功启动", "")
		fmt.Println("Pushdeer notify configuration is enabled. Test notification sent.")
	} else {
		fmt.Println("Pushdeer notify configuration is disabled.")
	}
}
