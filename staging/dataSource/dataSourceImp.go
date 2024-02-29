package dataSource

import (
	"TraderHelperCore/api"
	"TraderHelperCore/common"
)

type DataSourceImp struct {
	FetchData      func(string) string
	Formatter      api.Formatter
	DataUrlBuilder api.UrlBuilder
}

func NewDataSource(source Source) api.DataSource {
	var dataSourceInfo DataSourceInfo
	switch source {
	case SOURCE_TENCENT:
		dataSourceInfo = TencentDataSourceInfo()
	default:
		// 未知的数据源
		return nil
	}
	return DataSourceImp(dataSourceInfo)
}

func (dsi DataSourceImp) GetData(code string) common.StockData {
	dataType := common.GetDataTypeByCode(code)
	prefixType := common.GetPrefixTypeByCode(code)
	url := dsi.DataUrlBuilder.Build(dataType, prefixType, code)
	originData := dsi.FetchData(url)
	data := dsi.Formatter.Format(common.FormatTask{
		Code:       code,
		OriginData: originData,
		TargetType: dataType,
	})
	return data
}
