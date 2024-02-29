package dataSource

import (
	formatter "TraderHelperCore/staging/formatter"
	urlBuilder "TraderHelperCore/staging/urlBuilder"
	"io"
	"net/http"
)

type TencentDataSource struct{}

func TencentDataSourceInfo() DataSourceInfo {
	return DataSourceInfo{
		FetchData:      TencentDataFetcher,
		Formatter:      formatter.NewTencentDataFormatter(),
		DataUrlBuilder: urlBuilder.NewTencentUrlBuilder("https://qt.gtimg.cn/?q=", ""),
	}
}

func TencentDataFetcher(url string) string {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept-Encoding", "GB2312")
	req.Header.Set("Referer", "https://gu.qq.com/")
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(bodyBytes)
}
