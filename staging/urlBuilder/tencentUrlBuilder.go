package urlBuilder

import "TraderHelperCore/common"

type TencentUrlBuilder struct {
	UrlBase string
	Ext     string
}

func NewTencentUrlBuilder(urlBase string, ext string) *TencentUrlBuilder {
	return &TencentUrlBuilder{
		UrlBase: urlBase,
		Ext:     ext,
	}
}

func (t *TencentUrlBuilder) Build(dataType common.DataType, prefixType common.PrefixType, code string) string {
	var prefix string
	switch prefixType {
	case common.SH:
		prefix = "sh"
	case common.SZ:
		prefix = "sz"
	default:
		break
	}
	return t.UrlBase + prefix + code + t.Ext
}
