package initialize

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/service"
	"go.uber.org/zap"
)

var (
	menu = system.SysBaseMenu{
		MenuLevel: 0,
		ParentId:  "24",
		Path:      "wechatConfigManager",
		Name:      "wechatConfigManager",
		Hidden:    false,
		Component: "plugin/wechat/view/index.vue",
		Sort:      5,
		Meta: system.Meta{
			Title: "微信模块插件",
			Icon:  "aim",
		},
	}
	apis = []system.SysApi{
		{Path: "/wechat/private/config", Description: "获取微信配置", ApiGroup: "微信配置", Method: "GET"},
		{Path: "/wechat/private/config", Description: "修改微信配置", ApiGroup: "微信配置", Method: "PUT"},
		{Path: "/wechat/private/token", Description: "获取微信令牌", ApiGroup: "微信配置", Method: "GET"},
	}
	authorityId = uint(888)
)

// InitMenuAuthority 微信配置初始化
func InitMenuAuthority() {
	// 菜单权限初始化
	menuid, err := service.AddBaseMenu(menu)
	if err != nil && !errors.Is(err, service.ErrorMenuExits) {
		global.GlobalConfig.Log.Error("InitWechatConfig error", zap.Error(err))
		return
	}

	if menuid == 0 {
		return
	}

	err = service.SetMenuAuthority(menuid, authorityId)
	if err != nil {
		global.GlobalConfig.Log.Error("InitWechatConfig error", zap.Error(err))
		return
	}

	// api权限
	err = service.AddApiAuthority(authorityId, apis)
	if err != nil {
		return
	}
}
