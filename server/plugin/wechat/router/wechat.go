package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/api"
	"github.com/gin-gonic/gin"
)

type WechatRouter struct{}

func (s *WechatRouter) InitWechatRouter(Router *gin.RouterGroup) {
	wechatRouter := Router.Group("/public")
	wechatAuthRouter := Router.Group("/private").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	wechatApi := api.ApiGroupApp.WechatApi
	{
		wechatRouter.GET("jsapi", wechatApi.GetJsApiUsingPermissions) // 获取jssdk调用权限
		wechatRouter.GET("userInfo", wechatApi.GetSnsapiUserInfo)     // 获取(公众号)微信用户信息
	}

	{
		wechatAuthRouter.GET("config", wechatApi.GetConfig)    // 获取微信配置
		wechatAuthRouter.PUT("config", wechatApi.UpdateConfig) // 修改微信配置
		wechatAuthRouter.GET("token", wechatApi.GetToken)      // 获取微信令牌
	}

}
