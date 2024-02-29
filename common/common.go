package common

type Stock struct {
	Code      string
	BreakUp   float64 // 上破价
	BreakDown float64 // 下破价
}

type StockData struct {
	DataType  DataType
	Code      string
	Name      string
	LastPrice float64 // 最新价
	LastDate  string  // 最新时间
	LastTime  string  // 最新时间
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

func GetDataTypeByCode(code string) DataType {
	switch {
	case code[:1] == "6" || code[:1] == "0" || code[:1] == "3" || code[:1] == "9":
		return STOCK
	case code[:1] == "1" || code[:1] == "2" || code[:1] == "5":
		return FUND
	default:
		// 无法识别该证券代码的类别
		return -1
	}
}

func GetPrefixTypeByCode(code string) PrefixType {
	switch {
	case code[:1] == "6" || code[:1] == "5" || code[:1] == "9":
		return SH
	case code[:1] == "0" || code[:1] == "1" || code[:1] == "2" || code[:1] == "3":
		return SZ
	default:
		// 无法识别该代码所属的交易所
		return ""
	}
}
