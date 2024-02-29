package fomatter

import (
	"TraderHelperCore/common"
	"regexp"
	"strconv"
	"strings"
)

type tencentDataFormatter struct {
	regex *regexp.Regexp
}

func NewTencentDataFormatter() *tencentDataFormatter {
	return &tencentDataFormatter{
		regex: regexp.MustCompile("\".+?\""),
	}
}

func (tdf *tencentDataFormatter) Format(task common.FormatTask) common.StockData {
	matches := tdf.regex.FindStringSubmatch(task.OriginData)
	if len(matches) < 1 {
		return common.StockData{}
	}

	var result common.StockData
	switch task.TargetType {
	case common.STOCK, common.FUND:
		result = tdf.format(task.Code, matches[0])
	default:
		return common.StockData{}
	}

	return result
}

func (tdf *tencentDataFormatter) format(code string, dataStr string) common.StockData {
	slices := strings.Split(dataStr, "~")
	lastPrice, _ := strconv.ParseFloat(slices[3], 64)
	return common.StockData{
		DataType:  common.STOCK,
		Code:      code,
		Name:      slices[1],
		LastPrice: lastPrice,
		LastTime:  slices[30][8:10] + ":" + slices[30][10:12] + ":" + slices[30][12:14],
		LastDate:  slices[30][:4] + "/" + slices[30][4:6] + "/" + slices[30][6:8],
	}
}
