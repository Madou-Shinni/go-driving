package global

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Config struct {
	Rdb *redis.Client
	Log *zap.Logger
	DB  *gorm.DB
}

var GlobalConfig = new(Config)
