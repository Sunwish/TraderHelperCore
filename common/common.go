package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

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

type DataPack struct {
	FavoriteStocks map[string]Stock
	StocksData     map[string]StockData
	ActiveStocks   map[string]bool
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

type PushDeerConfig struct {
	BaseURL string
	Key     string
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

func SaveFavoriteStocksToFile(favoriteStocks map[string]Stock, outputDirectory string, outputFileName string) error {
	// 将map转换为json字节流
	jsonBytes, err := json.MarshalIndent(favoriteStocks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal favorite stocks: %w", err)
	}

	// 创建目录（如果不存在）
	err = os.MkdirAll(outputDirectory, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 将json字节写入到指定的文件路径
	err = os.WriteFile(path.Join(outputDirectory, outputFileName), jsonBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func LoadFavoriteStocksFromFile(filePath string) (map[string]Stock, error) {
	emptyMap := make(map[string]Stock)
	// 读取json文件内容
	jsonBytes, err := os.ReadFile(filePath)
	if err != nil {
		return emptyMap, fmt.Errorf("failed to read file: %w", err)
	}

	var favoriteStocks map[string]Stock
	// 反序列化json字节到map
	err = json.Unmarshal(jsonBytes, &favoriteStocks)
	if err != nil {
		return emptyMap, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return favoriteStocks, nil
}
