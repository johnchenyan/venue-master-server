package cache

import (
	"strings"

	"github.com/deatil/lakego-doak/lakego/array"
	"github.com/deatil/lakego-doak/lakego/cache"
	redisDriver "github.com/deatil/lakego-doak/lakego/cache/driver/redis"
	"github.com/deatil/lakego-doak/lakego/cache/interfaces"
	"github.com/deatil/lakego-doak/lakego/facade/config"
	"github.com/deatil/lakego-doak/lakego/facade/logger"
	"github.com/deatil/lakego-doak/lakego/register"
)

/**
 * 缓存
 *
 * cache.Default.Put("lakego-cache", "lakego-cache-data", 122222)
 * cache.Default.Forever("lakego-cache-forever", "lakego-cache-Forever-data")
 * cacheData, err := cache.Default.Get("lakego-cache")
 *
 * @create 2021-7-3
 * @author deatil
 */

// 默认
var Default *cache.Cache

// 初始化
func init() {
	// 注册默认
	registerDriver()

	// 默认
	Default = New()
}

// 实例化
func New(once ...bool) *cache.Cache {
	name := GetDefaultCache()

	return Cache(name, once...)
}

// 实例化
func NewWithType(name string, once ...bool) *cache.Cache {
	return Cache(name, once...)
}

func Cache(name string, once ...bool) *cache.Cache {
	conf := config.New("cache")

	// 缓存列表
	cfg := array.ArrayFrom(conf.GetStringMap("caches"))

	// 转为小写
	name = strings.ToLower(name)

	if !cfg.Has(name) {
		panic("缓存驱动[" + name + "]配置不存在")
	}

	// 配置
	driverConf := cfg.Value(name).ToStringMap()
	driverType := cfg.Value(name + ".type").ToString()

	driver := register.
		NewManagerWithPrefix("cache").
		GetRegister(driverType, driverConf, once...)
	if driver == nil {
		panic("缓存驱动[" + driverType + "]没有被注册")
	}

	c := cache.New(driver.(interfaces.Driver), driverConf)

	// 前缀
	keyPrefix := conf.GetString("key-prefix")
	c.WithPrefix(keyPrefix)

	return c
}

func GetDefaultCache() string {
	return config.New("cache").GetString("default")
}

// 注册
func registerDriver() {
	// 注册缓存驱动
	register.
		NewManagerWithPrefix("cache").
		Register("redis", func(conf map[string]any) any {
			cfg := array.ArrayFrom(conf)

			driver := redisDriver.New(redisDriver.Config{
				DB:       cfg.Value("db").ToInt(),
				Addr:     cfg.Value("addr").ToString(),
				Password: cfg.Value("password").ToString(),

				MinIdleConn:  cfg.Value("minidle-conn").ToInt(),
				DialTimeout:  cfg.Value("dial-timeout").ToDuration(),
				ReadTimeout:  cfg.Value("read-timeout").ToDuration(),
				WriteTimeout: cfg.Value("write-timeout").ToDuration(),

				PoolSize:    cfg.Value("pool-size").ToInt(),
				PoolTimeout: cfg.Value("pool-timeout").ToDuration(),

				EnableTrace: cfg.Value("enabletrace").ToBool(),

				Logger: logger.New(),
			})

			return driver
		})
}
