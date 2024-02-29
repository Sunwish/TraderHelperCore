package api

import (
	"TraderHelperCore/common"
)

type UrlBuilder interface {
	Build(dataType common.DataType, prefixType common.PrefixType, code string) string
}
