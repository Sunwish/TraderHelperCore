package urlBuilder

import (
	"TraderHelperCore/api"
	"TraderHelperCore/common"
)

type tencentUrlBuilder struct {
	urlBase string
	ext     string
}

func NewTencentUrlBuilder(urlBase string, ext string) api.UrlBuilder {
	return &tencentUrlBuilder{
		urlBase: urlBase,
		ext:     ext,
	}
}

func (t *tencentUrlBuilder) Build(dataType common.DataType, prefixType common.PrefixType, code string) string {
	var prefix string
	switch prefixType {
	case common.SH:
		prefix = "sh"
	case common.SZ:
		prefix = "sz"
	default:
		break
	}
	return t.urlBase + prefix + code + t.ext
}
