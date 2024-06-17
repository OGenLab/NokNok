package game

import (
	"strings"
	"time"

	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/pipeline"
	"github.com/lonng/nano/scheduler"
	"github.com/lonng/nano/serialize/json"
	"github.com/lonng/nano/session"
	log "github.com/sirupsen/logrus"
	"lightspeed.2dao3.com/nokserver/config"
	"lightspeed.2dao3.com/nokserver/internal/game/battlemanager"
	"lightspeed.2dao3.com/nokserver/internal/game/usermanager"
	"lightspeed.2dao3.com/nokserver/model"
	"lightspeed.2dao3.com/nokserver/model/t_player"
)

func Startup() {
	components := &component.Components{}
	components.Register(
		usermanager.DefaultUserManger,
		component.WithName("user"), // rewrite component and handler name
		component.WithNameFunc(strings.ToLower),
	)
	components.Register(
		battlemanager.DefaultBattleManager,
		component.WithName("battle"), // rewrite component and handler name
		component.WithNameFunc(strings.ToLower),
	)

	pip := pipeline.New()
	stats := &stats{}
	pip.Inbound().PushBack(stats.inbound)
	pip.Outbound().PushBack(stats.outbound)
	nano.Listen(config.Instance().Address+":"+config.Instance().Port,
		nano.WithIsWebsocket(true),
		nano.WithPipeline(pip),
		nano.WithHeartbeatInterval(time.Duration(config.Instance().NanoConfig.Heartbeat)*time.Second),
		nano.WithLogger(log.WithField("component", "nano")),
		nano.WithSerializer(json.NewSerializer()),
		nano.WithComponents(components),
		nano.WithDebugMode(),
		// nano.WithWSPath("/nok"),
	)

	go dailyRefresh()
}

func dailyRefresh() {
	var m = model.NewModel()
	if err := t_player.RefreshPlayer(m.DB); err != nil {
		m.Close()
		panic(err)
	}
	m.Close()

	// 计算从现在开始到下一次执行任务的时间需要等待多久
	duration := computeDuration()

	// 创建一个定时器，当到达指定的时间时，定时器就会向它的C字段发送一个时间值
	timer := time.NewTimer(duration)
	failFlag := false
	for {
		select {
		case <-timer.C:
			// 到达指定的时间，执行任务
			m := model.NewModel()
			if err := t_player.RefreshPlayer(m.DB); err != nil {
				log.Errorf("DailyRefresh faild: %v", err)
				m.Close()
				failFlag = true
				timer.Reset(1 * time.Minute)
				continue
			}
			m.Close()

			if failFlag {
				failFlag = false
				duration = computeDuration()
			} else {
				duration = 24 * time.Hour
			}

			// 任务执行完毕，重置定时器
			log.Infof("DailyRefresh succeeded, bext DailyRefresh scheduled after %v", duration)
			timer.Reset(duration)
		}
	}
}

func computeDuration() time.Duration {
	loc, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(loc)

	var next time.Time
	refreshTime := config.Instance().NokConfig.RefreshTime
	if now.Hour() < refreshTime {
		// 如果现在的时间还没到早上8点，那么下一次执行任务的时间就是今天的早上8点
		next = time.Date(now.Year(), now.Month(), now.Day(), refreshTime, 0, 0, 0, loc)
	} else {
		// 如果现在的时间已经过了早上8点，那么下一次执行任务的时间就是明天的早上8点
		next = time.Date(now.Year(), now.Month(), now.Day()+1, refreshTime, 0, 0, 0, loc)
	}

	return next.Sub(now)
}

type stats struct {
	component.Base
	timer         *scheduler.Timer
	outboundBytes int
	inboundBytes  int
}

func (stats *stats) outbound(s *session.Session, msg *pipeline.Message) error {
	stats.outboundBytes += len(msg.Data)
	return nil
}

func (stats *stats) inbound(s *session.Session, msg *pipeline.Message) error {
	stats.inboundBytes += len(msg.Data)
	return nil
}
