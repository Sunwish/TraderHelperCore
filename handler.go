package main

import (
	"TraderHelperCore/common"
	notifiers "TraderHelperCore/staging/notifier"
	"encoding/json"
	"fmt"
	"net/http"
)

func addFavoriteStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var newStock common.Stock
	err := json.NewDecoder(r.Body).Decode(&newStock)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("非法的请求数据"))
		return
	}

	// 检查favoriteStocks[newStock.Code]是否已存在
	if _, ok := favoriteStocks[newStock.Code]; ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("该自选股已存在"))
		return
	}

	// 验证并添加新自选股到favoriteStocks
	if !isCodeValid(newStock.Code) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("自选股代码不合法"))
		return
	}
	favoriteStocks[newStock.Code] = newStock

	common.SaveFavoriteStocksToFile(favoriteStocks, *dataDirectory, dataFileName)

	// 打印：成功添加 newStock 到自选股列表
	fmt.Println(len(favoriteStocks), "成功添加自选股", newStock)
}

func updateBreakPrice(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数，找到指定自选股并更新其上破价或下破价
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var newStock common.Stock
	err := json.NewDecoder(r.Body).Decode(&newStock)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("非法的请求数据"))
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

		common.SaveFavoriteStocksToFile(favoriteStocks, *dataDirectory, dataFileName)
		delete(activeStocks, stock.Code)

		fmt.Println(len(favoriteStocks), "成功更新自选股", oldStock, "→", stock)
	} else {
		fmt.Println(len(favoriteStocks), "指定要更新的自选股", newStock.Code, "不存在")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("该自选股不存在"))
		return
	}
}

func getFavoriteStocks(w http.ResponseWriter, r *http.Request) {
	// 序列化并返回favoriteStocks
	stocksJson, _ := json.Marshal(favoriteStocks)
	w.Header().Set("Content-Type", "application/json")
	w.Write(stocksJson)
}

func getFavoriteStocksData(w http.ResponseWriter, r *http.Request) {
	stocksDataMutex.RLock()
	stocksJson, _ := json.Marshal(stocksData)
	defer stocksDataMutex.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	w.Write(stocksJson)
}

func getActiveStocks(w http.ResponseWriter, r *http.Request) {
	activeStocksMutex.RLock()
	stocksJson, _ := json.Marshal(activeStocks)
	defer activeStocksMutex.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	w.Write(stocksJson)
}

func getFavoriteDataPack(w http.ResponseWriter, r *http.Request) {
	dataPack := &common.DataPack{
		ActiveStocks:   activeStocks,
		FavoriteStocks: favoriteStocks,
		StocksData:     stocksData,
	}
	dataJson, _ := json.Marshal(dataPack)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dataJson)
}

func removeFavoriteStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("非法的请求方式"))
		return
	}

	var stock common.Stock
	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("非法的请求数据"))
		return
	}

	delete(favoriteStocks, stock.Code)
	stocksDataMutex.Lock()
	delete(stocksData, stock.Code)
	stocksDataMutex.Unlock()

	common.SaveFavoriteStocksToFile(favoriteStocks, *dataDirectory, dataFileName)

	fmt.Println(len(favoriteStocks), "成功删除自选股", stock.Code)
}

func configNotifierPushdeer(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数，找到指定自选股并更新其上破价或下破价
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var config common.PushDeerConfig
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("非法的请求数据"))
		return
	}

	if config.BaseURL == "" || config.Key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BaseURL和Key不能为空"))
		return
	}

	notifier.AddNotifier(notifiers.NewPushdeerNotifier(config.BaseURL, config.Key))
}

func test_forceFetch(w http.ResponseWriter, r *http.Request) {
	stocksData["600123"] = ds.GetData("600123")
}
