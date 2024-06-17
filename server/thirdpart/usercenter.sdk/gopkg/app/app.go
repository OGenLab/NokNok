// Package app is a usercenter app helper. Other service can get config by appid.
package app

import (
	"errors"

	"qoobing.com/gomod/cache"
	"qoobing.com/gomod/log"
)

var (
	appCacher       cache.Cacher[App] = nil
	ErrAppNotExist                    = errors.New("app is not exist")
	ErrAppidInvalid                   = errors.New("app is invalid")
)

type App struct {
	AppId            string                 `json:"appid"`
	AppName          string                 `json:"appname"`
	AppType          string                 `json:"apptype"`
	AppTitle         string                 `json:"apptitle"`
	AppSecureKey     string                 `json:"appseckey"`
	AppSpecialConfig map[string]interface{} `json:"appconfigs"`
}

// GetAppByAppId get app config by appid.
func GetAppByAppId(appid string) (*App, error) {
	// Step 1. check appid checksum
	// TODO: check appid checksum
	if len(appid) < 8 {
		return nil, ErrAppidInvalid
	}

	// Step 2. get from cache
	var app, err = appCacher.GetFromCache(appid)
	if err != nil {
		return nil, ErrAppidInvalid
	}

	// Step 3. return
	log.PrintPretty("success cached/return app:", app)
	return app, nil
}

func (app *App) Config(name string) interface{} {
	if confstr, ok := app.AppSpecialConfig[name]; !ok {
		return nil
	} else {
		return confstr
	}
}

func (app *App) ConfigInt(name string, defaultvalue int) int {
	if c := app.Config(name); c != nil {
		return c.(int)
	}
	return defaultvalue
}

func (app *App) ConfigBool(name string, defaultvalue bool) bool {
	if c := app.Config(name); c != nil {
		return c.(bool)
	}
	return defaultvalue
}

func (app *App) ConfigString(name string, defaultvalue string) string {
	if c := app.Config(name); c != nil {
		return c.(string)
	}
	return defaultvalue
}

type Options = cache.Config
type AppStorager = cache.Getter[App]

func Init(storager AppStorager, options Options) {
	appCacher = cache.NewCache[App](storager, options)
}
