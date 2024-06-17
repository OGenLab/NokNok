/******************************************/
//
////////////////////////////////////////////
// Author:
// CreateTime:
/******************************************/
package main

import (
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"strings"
	"sync"
	"time"

	"lightspeed.2dao3.com/nokserver/common/app"
	"lightspeed.2dao3.com/nokserver/config"
	"lightspeed.2dao3.com/nokserver/internal/game"

	"qoobing.com/gomod/api"
	"qoobing.com/gomod/log"
)

func main() {
	// Step 1. init config
	cfg := config.Instance()
	app.Init()
	log.SetLogLevel(log.DEBUG)

	// Step 2. premetheus metric exporter & perfermance recorder
	startPrometheusExporter()
	startPprofRuntimeDebugger()

	// Step 3. run run run
	listen := cfg.Address + ":" + cfg.Port
	log.Infof("server start at [%s]", listen)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() { defer wg.Done(); game.Startup() }() // 开启游戏服
	wg.Wait()
}

func startPprofRuntimeDebugger() {
	go func() {
		pprofHostPort := "0.0.0.0:6060"
		log.Warningf("start pprof http server at [%s]", pprofHostPort)
		if err := http.ListenAndServe(pprofHostPort, nil); err != nil {
			log.Warningf("pprof http server(%s) stop with error:%s", pprofHostPort, err)
		}
	}()
}

// startPrometheusExporter
func startPrometheusExporter() {
	cfg := config.Instance()
	promcfg := cfg.Prometheus
	prommode := strings.ToUpper(promcfg.Mode)
	instance := cfg.Address + ":" + cfg.Port

	// Case 1: prometheus is off
	if prommode == "" || prommode == "OFF" {
		log.Infof("prometheus exporter is trun off, prommode=[%s]", prommode)
		return
	}
	api.InitMetrics(instance)

	// Case 2: prometheus pull mode support
	if strings.Contains(prommode, "PULL") {
		api.SetupMetricsExporterHandler("")
		log.Infof("prometheus exporter PULL trun on")
		go func() {
			l := promcfg.PullAddress + ":" + strconv.Itoa(promcfg.PullPort)
			engine := api.NewEngine()
			engine.POST("/metrics", api.MetricExporterHandler)
			err := engine.Run(l)
			log.Infof("prometheus push server exit with error:%s", err)
		}()
	}

	// Case 3: prometheus push mode support
	if strings.Contains(prommode, "PUSH") {
		api.SetupMetricsExporterPusher(promcfg.PushGateway, "")
		log.Infof("prometheus exporter PUSH trun on")
		go func() {
			for {
				log.Debugf("start push metrics ...")
				if err := api.MetricExporterPusher(); err != nil {
					log.Errorf("failed push metrics, retry 15 second later. err:[%s]", err)
					time.Sleep(15 * time.Second)
					continue
				}
				log.Debugf("success push metrics, do next push 15 second later")
				time.Sleep(15 * time.Second)
			}
		}()
	}

	log.Infof("startPrometheusExporter done")
}
