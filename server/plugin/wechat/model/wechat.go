package model

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/pkg/tools"

// Wechat 微信配置
type Wechat struct {
	// 小程序
	MiniProgramEnabled bool `json:"miniProgramEnabled"`
	MiniProgram        struct {
		AppId     string `json:"appId"`
		AppSecret string `json:"appSecret"`
	} `json:"miniProgram"`

	// 公众号
	OfficialAccountEnabled bool `json:"officialAccountEnabled"`
	OfficialAccount        struct {
		AppId     string `json:"appId"`
		AppSecret string `json:"appSecret"`
	} `json:"officialAccount"`
}

func (w Wechat) ToWxConfig() *tools.WxConfig {
	miniProgram := w.MiniProgram
	officialAccount := w.OfficialAccount
	return tools.NewWxConfig(miniProgram.AppId, miniProgram.AppSecret, officialAccount.AppId, officialAccount.AppSecret)
}
