package dataSource

import "TraderHelperCore/api"

type Source int

const (
	SOURCE_SINA Source = iota
	SOURCE_TENCENT
	SOURCE_EASTMONEY
)

type DataSourceInfo struct {
	FetchData      func(string) string
	Formatter      api.Formatter
	DataUrlBuilder api.UrlBuilder
}
