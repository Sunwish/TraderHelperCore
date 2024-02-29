package formatter

import (
	"TraderHelperCore/common"
)

type Formatter interface {
	Format(task common.FormatTask) common.StockData
}
