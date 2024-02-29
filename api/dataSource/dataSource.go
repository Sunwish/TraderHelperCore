package dataSource

import "TraderHelperCore/common"

type DataSource interface {
	GetData(code string) common.StockData
}
