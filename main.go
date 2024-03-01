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
var ds = dataSource.NewDataSource(dataSource.SOURCE_TENCENT) // 数据源
var notifier = notifiers.NewLogNotifier()
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

	// 配置Pushdeer通知
	if pushdeerBaseUrl != nil && *pushdeerBaseUrl != "" && pushdeerKey != nil && *pushdeerKey != "" {
		notifier = notifiers.NewPushdeerNotifier(*pushdeerBaseUrl, *pushdeerKey)
		notifier.Notify("[TraderHelper] 预警服务已成功启动", "")
		fmt.Println("Pushdeer notify configuration is enabled. Test notification sent.")
	} else {
		fmt.Println("Pushdeer notify configuration is disabled.")
	}

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

	// 设置核心路由和处理函数
	mux := http.NewServeMux()
	mux.HandleFunc("/add_favorite_stock", addFavoriteStock)
	mux.HandleFunc("/update_break_price", updateBreakPrice)
	mux.HandleFunc("/get_favorite_stocks", getFavoriteStocks)
	mux.HandleFunc("/get_favorite_stocks_data", getFavoriteStocksData)
	mux.HandleFunc("/remove_favorite_stock", removeFavoriteStock)
	mux.HandleFunc("/test/force_fetch", test_forceFetch)
	// 跨域处理
	corsHandler := corsWrapper(mux)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: corsHandler,
	}

	// 启动核心服务
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

func corsWrapper(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有源或指定具体的源地址
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	}
}

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
		notifier.Notify(fmt.Sprintf("[%s] 数据获取失败", stock.Code), "请检查网络连接")
		return
	}
	// update
	stocksDataMutex.Lock()
	stocksData[stock.Code] = newData
	stocksDataMutex.Unlock()
	// 判断上破下破
	stockConfig := favoriteStocks[stock.Code]
	if stockConfig.BreakUp > 0 && newData.LastPrice >= stockConfig.BreakUp {
		notifier.Notify(fmt.Sprintf("[%s] %s 触发上破", newData.Code, newData.Name), fmt.Sprintf("现价：%f，上破 %f", newData.LastPrice, stockConfig.BreakUp))
	}
	if stockConfig.BreakDown > 0 && newData.LastPrice <= stockConfig.BreakDown {
		notifier.Notify(fmt.Sprintf("[%s] %s 触发下破", newData.Code, newData.Name), fmt.Sprintf("现价：%f，下破 %f", newData.LastPrice, stockConfig.BreakDown))
	}
}

func isCodeValid(code string) bool {
	data := ds.GetData(code)
	return data.LastPrice != 0
}
