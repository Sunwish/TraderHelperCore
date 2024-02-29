package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Stock struct {
	Code      string
	BreakUp   float64 // 上破价
	BreakDown float64 // 下破价
}

type StockData struct {
	Code      string
	BreakUp   float64   // 上破价
	BreakDown float64   // 下破价
	LastPrice float64   // 最新价
	LastTime  time.Time // 最新时间
}

var favoriteStocks = make(map[string]Stock) // 用于存储自选股列表

func addFavoriteStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var newStock Stock
	err := json.NewDecoder(r.Body).Decode(&newStock)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 验证并添加新自选股到favoriteStocks
	favoriteStocks[newStock.Code] = newStock
	// 打印：成功添加 newStock 到自选股列表
	fmt.Println(len(favoriteStocks), "成功添加自选股", newStock)
}

func updateBreakPrice(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数，找到指定自选股并更新其上破价或下破价
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var newStock Stock
	err := json.NewDecoder(r.Body).Decode(&newStock)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 从favoriteStocks中找到指定自选股并更新其上破价或下破价
	if stock, ok := favoriteStocks[newStock.Code]; ok {
		// 备份更新前的stock
		oldStock := stock
		// 如果上破价存在则更新上破，如果下破价存在则更新下破
		if newStock.BreakUp > 0 {
			stock.BreakUp = newStock.BreakUp
		}
		if newStock.BreakDown > 0 {
			stock.BreakDown = newStock.BreakDown
		}
		favoriteStocks[stock.Code] = stock
		fmt.Println(len(favoriteStocks), "成功更新自选股", oldStock, "→", stock)
	} else {
		// 打印：指定的自选股不存在
		fmt.Println(len(favoriteStocks), "指定要更新的自选股", newStock.Code, "不存在")
	}
}

func getFavoriteStocks(w http.ResponseWriter, r *http.Request) {
	// 序列化并返回favoriteStocks
	stocksJson, _ := json.Marshal(favoriteStocks)
	w.Header().Set("Content-Type", "application/json")
	w.Write(stocksJson)
}

func getFavoriteStockData(stock Stock) StockData {
	return StockData{}
}

func getFavoriteStocksData(w http.ResponseWriter, r *http.Request) {
	// 遍历favoriteStocks，调用getFavoriteStockData获得
	var stockDatas []StockData
	for _, stock := range favoriteStocks {
		// 调用getStockData获取stockData
		stockData := getFavoriteStockData(stock)
		stockDatas = append(stockDatas, stockData)
	}
	// 序列化并返回favoriteStocks
	stockDatasJson, _ := json.Marshal(stockDatas)
	w.Header().Set("Content-Type", "application/json")
	w.Write(stockDatasJson)
}

func removeFavoriteStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var newStock Stock
	err := json.NewDecoder(r.Body).Decode(&newStock)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	delete(favoriteStocks, newStock.Code)
	fmt.Println(len(favoriteStocks), "成功删除自选股", newStock.Code)
}

func fetchAndUpdateStockPrice(stock Stock) {
	// 实现获取最新股价并更新到stock.CurrentPrice的功能
	// ...
}

func main() {
	// 初始化ticker
	ticker := time.NewTicker(1 * time.Second)

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
