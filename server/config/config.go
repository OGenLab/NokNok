// Copyright 2022 The LSC team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package config

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pelletier/go-toml/v2"
	"qoobing.com/gomod/log"
	"qoobing.com/gomod/redis/sentinel"
	"qoobing.com/gomod/str"
)

// config	服务配置结构体，
//
//	是文件conf/$SERVERNAME.conf的内存表现形式。
type config struct {
	Debug      bool             `toml:"debug"   default:false`     //危险调试开关
	Port       string           `toml:"port"    default:"8787"`    //服务监听端口
	Address    string           `toml:"address" default:"0.0.0.0"` //服务监听地址
	NokDb      databaseConfig   `toml:"database.nokdb"`            //交易核心数据库
	NokRedis   sentinel.Config  `toml:"redis.nokredis"`            //交易核心缓存
	Prometheus prometheusConfig `toml:"prometheus"`                //prometheus
	NanoConfig nanoConfig       `toml:"nanoConfig"`
	NokConfig  nokConfig        `toml:"nokConfig"`
	TgConfig   tgConfig         `toml:"tgConfig"`
}

type nanoConfig struct {
	Heartbeat int `toml:"heartbeat"`
}

type nokConfig struct {
	Season          int `toml:"season"`
	GopherCnt       int `toml:"gopherCnt"`
	GopherRound     int `toml:"gopherRound"`
	MaxStamina      int `toml:"maxStamina"`
	RefreshInterval int `toml:"refreshInterval"` // 静态表刷新间隔
	RefreshTime     int `toml:"refreshTime"`     // 全局状态刷新时间点, 例如 8 表示早 8
}

type tgConfig struct {
	BotToken string `toml:"botToken"`
}
type usercenterConfig struct {
	URL                  string `toml:"url"`
	IsUserInfoValidation string `toml:"is_userinfo_validation"`
	AppId                string `toml:"appid"`
	AppSecretKey         string `toml:"appsecretkey"`
}

type databaseConfig struct {
	Host            string `toml:"host"`             //数据库名称
	Port            int    `toml:"port"`             //数据库名称
	Debug           bool   `toml:"debug"`            //调试开关（会在日志打印SQL)
	Dbname          string `toml:"dbname"`           //数据库名称
	Username        string `toml:"username"`         //数据库用户名
	Password        string `toml:"password"`         //数据库连接密码
	ExtraParameters string `toml:"extra_parameters"` //数据库连接扩展参数
}

type prometheusConfig struct {
	Debug       bool   `toml:"debug"`        //调试开关
	Mode        string `toml:"mode"`         //prometheus数据上报模式
	PullPort    int    `toml:"pull_port"`    //promhttp服务监听端口
	PullAddress string `toml:"pull_address"` //promhttp服务监听地址
	PushGateway string `toml:"push_gateway"` //prometheus监控推送pushgateway
}

func (c config) IsJustDebugCodeOpen() bool {
	return c.Debug
}

func getConfigFile() string {
	suffix := ""
	curenv := strings.ToUpper(str.GetEnvDefault("DEVTOENV", ""))
	switch curenv {
	case "":
		suffix = ""
	case "DEV":
		suffix = ".dev"
	case "TST":
		suffix = ".tst"
	case "PRE":
		suffix = ".pre"
	case "PRD":
		suffix = ".prd"
	default:
		panic("enviroment value 'DEVTOENV' should be oneof {DEV,TST,PRE,PRD}")
	}

	_, exename := filepath.Split(os.Args[0])
	conffile := flag.String("config", "./etc/"+exename+".conf"+suffix, "config file")
	return *conffile
}

var (
	cfg  config    //配置数据实例
	once sync.Once //保证数据实例只初始化一次
)

func initConfig() {
	defer log.PrintPretty("config:", &cfg)
	cfg.Debug = false

	conffile := getConfigFile()
	doc, err := ioutil.ReadFile(conffile)
	if err != nil {
		panic("initial config, read config file error:" + err.Error())
	}
	if err := toml.Unmarshal(doc, &cfg); err != nil {
		panic("initial config, unmarshal config file error:" + err.Error())
	}
}

// 获取全局唯一Config实例
func Instance() *config {
	once.Do(initConfig)
	return &cfg
}
