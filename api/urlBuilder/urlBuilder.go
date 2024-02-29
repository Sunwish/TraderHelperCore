package urlBuilder

import (
	"TraderHelperCore/common"
)

type UrlBuilderInterface interface {
	Build(dataType common.DataType, prefixType common.PrefixType, code string) string
}
