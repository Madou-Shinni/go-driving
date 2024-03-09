package wechat

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/initialize"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/jobs"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wechat/router"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type wechatPlugin struct{}

func CreateWechatPlug(rdb *redis.Client, db *gorm.DB, log *zap.Logger) *wechatPlugin {
	global.GlobalConfig.Rdb = rdb
	global.GlobalConfig.DB = db
	global.GlobalConfig.Log = log
	return &wechatPlugin{}
}

func (*wechatPlugin) Register(group *gin.RouterGroup) {
	initialize.InitWechatConfig()
	initialize.InitMenuAuthority()
	jobs.CronInit()
	router.RouterGroupApp.InitWechatRouter(group)
}

func (*wechatPlugin) RouterPath() string {
	return "wechat"
}
