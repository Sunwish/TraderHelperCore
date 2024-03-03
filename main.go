package main

import (
	"TraderHelperCore/api"
	"TraderHelperCore/common"
	"TraderHelperCore/staging/dataSource"
	n "TraderHelperCore/staging/notification"
	notifiers "TraderHelperCore/staging/notifier"
	"context"
	"flag"
	"fmt"
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
			notification := n.MakeNotification("Recovered from a panic!", fmt.Sprint(r), n.WithExtType(api.NotificationTypeSystem))
			notifier.Notify(*notification)
		}
	}()

	// 解析启动参数
	Logln("Parsing flags...")
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
			notification := n.MakeNotification("Fail to ListenAndServe", err.Error(), n.WithExtType(api.NotificationTypeSystem), n.WithExtCode(int(api.NotificationCodeServerStartupFail)))
			notifier.Notify(*notification)
			LogFatal("ListenAndServe: ", err)
		}
	}()
	Logf("Server is running at :%d\n", *port)

	// 确保程序在接收到信号时优雅退出
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	// 停止ticker
	ticker.Stop()
	// 关闭HTTP服务器
	if err := server.Shutdown(context.Background()); err != nil {
		LogFatal("Server Shutdown: ", err)
	}

	Logln("Server stopped.")
}

func configNotifier() {
	// Pushdeer Notifier
	if pushdeerBaseUrl != nil && *pushdeerBaseUrl != "" && pushdeerKey != nil && *pushdeerKey != "" {
		notifier.AddNotifier(notifiers.NewPushdeerNotifier(*pushdeerBaseUrl, *pushdeerKey))
		notifier.Notify(*n.MakeNotification("Pushdeer 预警服务成功启动", ""))
		Logln("Pushdeer notify configuration is enabled. Test notification sent.")
	} else {
		Logln("Pushdeer notify configuration is disabled.")
	}

	// TCP Notifier
	notifier.AddNotifier(notifiers.NewTcpNotifier(49121))
}
