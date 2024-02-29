package common

import (
	"errors"
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

type DataType int

const (
	STOCK DataType = iota
	FUND
)

type PrefixType string

const (
	SH PrefixType = "SH" // 上交所
	SZ PrefixType = "SZ" // 深交所
)

type FormatTask struct {
	Code       string
	OriginData string
	TargetType DataType
}

func getDataTypeByCode(code string) (DataType, error) {
	switch {
	case code[:1] == "6" || code[:1] == "0" || code[:1] == "3" || code[:1] == "9":
		return STOCK, nil
	case code[:1] == "1" || code[:1] == "2" || code[:1] == "5":
		return FUND, nil
	default:
		return 0, errors.New("无法识别该证券代码的类别")
	}
}

func getPrefixTypeByCode(code string) (PrefixType, error) {
	switch {
	case code[:1] == "6" || code[:1] == "5" || code[:1] == "9":
		return SH, nil
	case code[:1] == "0" || code[:1] == "1" || code[:1] == "2" || code[:1] == "3":
		return SZ, nil
	default:
		return "", errors.New("无法识别该代码所属的交易所")
	}
}
