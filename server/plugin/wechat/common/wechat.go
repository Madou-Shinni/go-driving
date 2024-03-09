package common

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/pkg/tools"
	"go.uber.org/zap"
)

// GetWechatAccessToken 从中控服务器获取微信AccessToken
func GetWechatAccessToken() string {
	ctx := context.Background()
	rdb := global.GlobalConfig.Rdb

	// 获取微信配置
	config, err := GetWechatConfig()
	if err != nil {
		global.GlobalConfig.Log.Error("GetWechatAccessToken 获取微信配置失败", zap.Error(err))
		return ""
	}

	// 禁用
	if !config.MiniProgramEnabled {
		global.GlobalConfig.Log.Warn("GetWechatAccessToken 小程序配置禁用", zap.Bool("MiniProgramEnabled", config.MiniProgramEnabled))
		return ""
	}

	appId := config.MiniProgram.AppId

	// 从中控服务器获取access token
	accessToken, err := tools.GetRedisStrResult[string](rdb, ctx, GetWeChatAccessTokenKey(appId))
	if err != nil {
		return ""
	}
	return accessToken
}

// GetWechatPublicAccessToken 从中控服务器获取微信公众号AccessToken
func GetWechatPublicAccessToken() string {
	ctx := context.Background()
	rdb := global.GlobalConfig.Rdb

	// 获取微信配置
	config, err := GetWechatConfig()
	if err != nil {
		global.GlobalConfig.Log.Error("GetWechatAccessToken 获取微信配置失败", zap.Error(err))
		return ""
	}

	// 禁用
	if !config.OfficialAccountEnabled {
		global.GlobalConfig.Log.Warn("GetWechatPublicAccessToken 公众号配置禁用", zap.Bool("OfficialAccountEnabled", config.OfficialAccountEnabled))
		return ""
	}

	appId := config.OfficialAccount.AppId

	// 从中控服务器获取access token
	accessToken, err := tools.GetRedisStrResult[string](rdb, ctx, GetWeChatAccessTokenKey(appId))
	if err != nil {
		return ""
	}
	return accessToken
}

// GetWechatPublicJsApiTicket 从中控服务器获取微信公众号JsApiTicket
func GetWechatPublicJsApiTicket() string {
	ctx := context.Background()
	rdb := global.GlobalConfig.Rdb

	// 获取微信配置
	config, err := GetWechatConfig()
	if err != nil {
		global.GlobalConfig.Log.Error("GetWechatAccessToken 获取微信配置失败", zap.Error(err))
		return ""
	}

	// 禁用
	if !config.OfficialAccountEnabled {
		global.GlobalConfig.Log.Warn("GetWechatPublicJsApiTicket 公众号配置禁用", zap.Bool("OfficialAccountEnabled", config.OfficialAccountEnabled))
		return ""
	}

	appId := config.OfficialAccount.AppId

	// 从中控服务器获取JsApiTicket
	jsApiTicket, err := tools.GetRedisStrResult[string](rdb, ctx, GetWeChatJsApiTicketKey(appId))
	if err != nil {
		return ""
	}
	return jsApiTicket
}

// GetWechatConfig 获取微信配置
func GetWechatConfig() (config model.Wechat, err error) {
	ctx := context.Background()
	rdb := global.GlobalConfig.Rdb
	// 获取微信配置
	config, err = tools.GetRedisStrResult[model.Wechat](rdb, ctx, GetWechatConfigKey())
	if err != nil {
		return
	}

	return config, nil
}
