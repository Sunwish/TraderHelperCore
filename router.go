package main

import "net/http"

func configRoute(mux *http.ServeMux) {
	// 前端
	mux.Handle("/", http.FileServer(http.Dir("./www")))

	// 数据 API
	mux.HandleFunc("/add_favorite_stock", addFavoriteStock)
	mux.HandleFunc("/update_break_price", updateBreakPrice)
	mux.HandleFunc("/get_favorite_stocks", getFavoriteStocks)
	mux.HandleFunc("/get_favorite_stocks_data", getFavoriteStocksData)
	mux.HandleFunc("/get_active_stocks", getActiveStocks)
	mux.HandleFunc("/remove_favorite_stock", removeFavoriteStock)
	mux.HandleFunc("/get_favorite_data_pack", getFavoriteDataPack)

	// 配置 API
	mux.HandleFunc("/config_notifier_pushdeer", configNotifierPushdeer)

	// 测试 API
	mux.HandleFunc("/test/force_fetch", test_forceFetch)
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
