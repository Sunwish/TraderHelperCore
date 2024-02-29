package dataSource

import (
	formatter "TraderHelperCore/staging/formatter"
	urlBuilder "TraderHelperCore/staging/urlBuilder"
	"bytes"
	"io"
	"net/http"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
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
	req.Header.Set("Referer", "https://gu.qq.com/")
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyBytes, _ = gbkToUtf8(bodyBytes)
	if err != nil {
		return ""
	}
	return string(bodyBytes)
}

// GBK 转 UTF-8
func gbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
