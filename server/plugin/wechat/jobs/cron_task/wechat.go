package cron_task

import (
	"context"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/common"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/pkg/tools"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

// RefreshAccessToken 刷新微信AccessToken
func RefreshAccessToken() {
	defer recoverPanic()

	ctx := context.Background()
	rdb := global.GlobalConfig.Rdb
	// 获取微信配置
	config, err := common.GetWechatConfig()
	if err != nil {
		global.GlobalConfig.Log.Error("RefreshAccessToken 获取微信配置失败", zap.Error(err))
		return
	}

	// 禁用
	if !config.MiniProgramEnabled {
		return
	}

	wechatConfig := config.ToWxConfig()
	appId := wechatConfig.AppID

	// 从中控服务器获取access token
	accessToken, err := tools.GetRedisStrResult[string](rdb, ctx, common.GetWeChatAccessTokenKey(appId))
	if err != nil && !errors.Is(err, redis.Nil) {
		return
	}

	// 获取过期时间
	ttl, err := tools.TTLResult(rdb, ctx, common.GetWeChatAccessTokenKey(appId))
	if err != nil {
		global.GlobalConfig.Log.Error("获取access token过期时间失败", zap.Error(err))
		return
	}

	if accessToken != "" && ttl.Minutes() > (time.Minute*5).Minutes() {
		// 存在，且过期时间大于5分钟 不刷新
		return
	}

	// 获取微信access token
	result, err := wechatConfig.GetStableAccessToken(false)
	if err != nil {
		global.GlobalConfig.Log.Error("刷新微信AccessToken失败", zap.Error(err))
		return
	}
	if result.ErrCode != 0 {
		global.GlobalConfig.Log.Error("刷新微信AccessToken失败", zap.String("errmsg", result.ErrMsg))
		return
	}

	// 将access token缓存之中控服务器
	_, err = tools.SetRedisStrResult[string](rdb, ctx, common.GetWeChatAccessTokenKey(appId), result.AccessToken, time.Second*time.Duration(result.ExpiresIn))
	if err != nil {
		global.GlobalConfig.Log.Error("access token缓存之中控服务器失败", zap.Error(err))
		return
	}
}

// RefreshPublicAccessToken 刷新微信公众号AccessToken
func RefreshPublicAccessToken() {
	defer recoverPanic()

	ctx := context.Background()
	rdb := global.GlobalConfig.Rdb
	// 获取微信配置
	config, err := common.GetWechatConfig()
	if err != nil {
		global.GlobalConfig.Log.Error("RefreshAccessToken 获取微信配置失败", zap.Error(err))
		return
	}

	// 禁用
	if !config.OfficialAccountEnabled {
		return
	}

	wechatConfig := config.ToWxConfig()
	appId := wechatConfig.PubWxConfig.AppID

	// 从中控服务器获取access token
	accessToken, err := tools.GetRedisStrResult[string](rdb, ctx, common.GetWeChatAccessTokenKey(appId))
	if err != nil && !errors.Is(err, redis.Nil) {
		return
	}

	// 获取过期时间
	ttl, err := tools.TTLResult(rdb, ctx, common.GetWeChatAccessTokenKey(appId))
	if err != nil {
		global.GlobalConfig.Log.Error("获取access token过期时间失败", zap.Error(err))
		return
	}

	if accessToken != "" && ttl.Minutes() > (time.Minute*5).Minutes() {
		// 存在，且过期时间大于5分钟 不刷新
		return
	}

	// 获取微信access token
	result, err := wechatConfig.PubWxConfig.GetStableAccessToken(false)
	if err != nil {
		global.GlobalConfig.Log.Error("刷新微信公众号AccessToken失败", zap.Error(err))
		return
	}
	if result.ErrCode != 0 {
		global.GlobalConfig.Log.Error("刷新微信公众号AccessToken失败", zap.String("errmsg", result.ErrMsg))
		return
	}

	// 将access token缓存之中控服务器
	_, err = tools.SetRedisStrResult[string](rdb, ctx, common.GetWeChatAccessTokenKey(appId), result.AccessToken, time.Second*time.Duration(result.ExpiresIn))
	if err != nil {
		global.GlobalConfig.Log.Error("微信公众号access token缓存之中控服务器失败", zap.Error(err))
		return
	}
}

// RefreshPublicJsApiTicket 刷新微信公众号JsApiTicket
func RefreshPublicJsApiTicket() {
	defer recoverPanic()

	ctx := context.Background()
	rdb := global.GlobalConfig.Rdb
	// 获取微信配置
	config, err := common.GetWechatConfig()
	if err != nil {
		global.GlobalConfig.Log.Error("RefreshAccessToken 获取微信配置失败", zap.Error(err))
		return
	}

	// 禁用
	if !config.OfficialAccountEnabled {
		return
	}

	wechatConfig := config.ToWxConfig()
	appId := wechatConfig.PubWxConfig.AppID

	// 从中控服务器获取js api ticket
	jsApiTicket, err := tools.GetRedisStrResult[string](rdb, ctx, common.GetWeChatJsApiTicketKey(appId))
	if err != nil && !errors.Is(err, redis.Nil) {
		return
	}

	// 获取过期时间
	ttl, err := tools.TTLResult(rdb, ctx, common.GetWeChatJsApiTicketKey(appId))
	if err != nil {
		global.GlobalConfig.Log.Error("获取js api ticket过期时间失败", zap.Error(err))
		return
	}

	if jsApiTicket != "" && ttl.Minutes() > (time.Minute*5).Minutes() {
		// 存在，且过期时间大于5分钟 不刷新
		return
	}

	// 从中控服务器获取access token
	accessToken, err := tools.GetRedisStrResult[string](rdb, ctx, common.GetWeChatAccessTokenKey(appId))
	if err != nil && !errors.Is(err, redis.Nil) {
		return
	}
	if accessToken == "" {
		// AccessToken不存在无法刷新
		return
	}

	// 获取微信js api ticket
	result, err := tools.GetJsapiTicket(accessToken)
	if err != nil {
		global.GlobalConfig.Log.Error("获取微信js api ticket失败", zap.Error(err))
		return
	}
	if result.ErrCode != 0 {
		global.GlobalConfig.Log.Error("获取微信js api ticket失败", zap.String("errmsg", result.ErrMsg))
		return
	}

	// 将js api ticket缓存之中控服务器
	_, err = tools.SetRedisStrResult[string](rdb, ctx, common.GetWeChatJsApiTicketKey(appId), result.Ticket, time.Second*time.Duration(result.ExpiresIn))
	if err != nil {
		global.GlobalConfig.Log.Error("js api ticket缓存之中控服务器失败", zap.Error(err))
		return
	}
}
