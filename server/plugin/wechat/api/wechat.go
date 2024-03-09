package api

import (
	"context"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/common"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/pkg/tools"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type WechatApi struct{}

// GetJsApiUsingPermissions 获取用于调用微信JS接口的使用权限
// @Tags WeChat
// @Summary 微信JS接口的使用权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param url query string true "url"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /wechat/public/jsapi [get]
func (cuApi *WechatApi) GetJsApiUsingPermissions(c *gin.Context) {
	var url string
	url = c.Query("url")
	if url == "" {
		response.FailWithMessage("empty url", c)
		return
	}

	// 获取微信配置
	config, err := common.GetWechatConfig()
	if err != nil {
		response.FailWithMessage("empty code", c)
		return
	}
	wechatConfig := config.ToWxConfig()

	noncestr := tools.NewRandGenerator(tools.WithLength(12)).Generate()
	jsapi_ticket := common.GetWechatPublicJsApiTicket()
	data := tools.GetJsApiUsingPermissions{
		Noncestr:    noncestr,
		JsapiTicket: jsapi_ticket,
		Timestamp:   time.Now().Unix(),
		Url:         url,
	}

	sha1 := data.Sha1()
	appId := wechatConfig.PubWxConfig.AppID // 公众号appId

	result := tools.GetJsApiUsingPermissionsResp{
		GetJsApiUsingPermissions: data,
		AppId:                    appId,
		Signature:                sha1,
	}

	response.OkWithData(result, c)
}

// GetSnsapiUserInfo 获取(公众号)微信用户信息
// @Tags WeChat
// @Summary 获取(公众号)微信用户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param code query string true "用户同意授权，获取code"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /wechat/public/userInfo [get]
func (cuApi *WechatApi) GetSnsapiUserInfo(c *gin.Context) {
	// 获取code
	code := c.Query("code")
	if code == "" {
		response.FailWithMessage("empty code", c)
		return
	}

	// 获取微信配置
	config, err := common.GetWechatConfig()
	if err != nil {
		response.FailWithMessage("empty code", c)
		return
	}
	wechatConfig := config.ToWxConfig()

	// 获取web access token
	result, err := wechatConfig.PubWxConfig.GetWebAccessToken(code)
	if err != nil {
		response.FailWithMessage("get access token failed", c)
		return
	}
	if result.ErrCode != 0 {
		response.FailWithMessage(result.ErrMsg, c)
		return
	}

	// 获取用户信息
	userInfo, err := tools.GetSnsapiUserInfo(result.AccessToken, result.OpenId)
	if err != nil {
		response.FailWithMessage("get user info failed", c)
		return
	}

	response.OkWithData(userInfo, c)
}

// GetConfig 获取配置
// @Tags WeChat
// @Summary 获取配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /wechat/private/config [get]
func (cuApi *WechatApi) GetConfig(c *gin.Context) {
	// 获取微信配置
	config, err := common.GetWechatConfig()
	if err != nil && !errors.Is(err, redis.Nil) {
		response.FailWithMessage("获取微信配置失败", c)
		return
	}

	if errors.Is(err, redis.Nil) {
		response.OkWithData(model.Wechat{}, c)
		return
	}

	response.OkWithData(config, c)
}

// UpdateConfig 更新配置
// @Tags WeChat
// @Summary 更新配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param  data  body  model.Wechat true  "wechat配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /wechat/private/config [put]
func (cuApi *WechatApi) UpdateConfig(c *gin.Context) {
	var wechatConfig model.Wechat
	rdb := global.GlobalConfig.Rdb
	log := global.GlobalConfig.Log
	ctx := context.Background()

	err := c.ShouldBindJSON(&wechatConfig)
	if err != nil {
		log.Error("更新微信配置失败，请求参数异常", zap.Error(err))
		response.FailWithMessage("参数异常", c)
		return
	}

	_, err = tools.SetRedisStrResult[model.Wechat](rdb, ctx, common.GetWechatConfigKey(), wechatConfig, -time.Second)
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Error("更新微信配置失败", zap.Error(err))
		response.FailWithMessage("更新微信配置失败", c)
		return
	}

	response.Ok(c)
}

// GetToken 获取微信令牌
// @Tags WeChat
// @Summary 获取微信令牌
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param  type  query  string true  "令牌类型类型 miniProgram | officialAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /wechat/private/token [put]
func (cuApi *WechatApi) GetToken(c *gin.Context) {
	tp := c.Query("type")

	var data = struct {
		AccessToken string `json:"accessToken"`
	}{}

	if tp == "miniProgram" {
		data.AccessToken = common.GetWechatAccessToken()
		response.OkWithData(data, c)
		return
	}

	if tp == "officialAccount" {
		data.AccessToken = common.GetWechatPublicAccessToken()
		response.OkWithData(data, c)
		return
	}
}
