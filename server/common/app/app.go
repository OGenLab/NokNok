package app

import (
	"lightspeed.2dao3.com/nokserver/model"
	"lightspeed.2dao3.com/wallet/usercenter/sdk/app"
	"qoobing.com/gomod/cache"
)

var (
	GetAppByAppId = app.GetAppByAppId
)

func Init() {
	var cfg = cache.Config{
		UseLocalCache:            true,
		LocalCacheLifetimeSecond: 60,
		UseRedisCache:            true,
		RedisCacheLifetimeSecond: -1,
		RedisCacheKeyPrefix:      "appconfig-sdk-info",
	}
	m := model.NewModelDefaultRedis()
	cfg.RedisCacheConnPool = m.RedisPool
	app.Init(nil, cfg)
}
