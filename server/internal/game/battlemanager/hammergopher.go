package battlemanager

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/lonng/nano/session"
	"lightspeed.2dao3.com/nokserver/common/boost"
	"lightspeed.2dao3.com/nokserver/common/util"
	"lightspeed.2dao3.com/nokserver/config"
	"lightspeed.2dao3.com/nokserver/consts/errcode"
	"lightspeed.2dao3.com/nokserver/internal/protocol"
	"lightspeed.2dao3.com/nokserver/model"
	"lightspeed.2dao3.com/nokserver/model/cache"
	"lightspeed.2dao3.com/nokserver/model/t_hammer_entities"
	"lightspeed.2dao3.com/nokserver/model/t_player"
	"qoobing.com/gomod/log"
)

type HammerGopherRequest struct {
	Round string `json:"round"`
	Index int    `json:"index"`
}

type HammerGopherResponse struct {
	protocol.BaseResponse
	Coins       uint64                `json:"coins"`
	NokInfo     *protocol.NokInfo     `json:"nokInfo"`
	AccountInfo *protocol.AccountInfo `json:"accountInfo"`
}

func (u *BattleManager) HammerGopher(s *session.Session, req *HammerGopherRequest) error {
	var (
		m   = model.NewModelDefault()
		uid = s.UID()
	)
	defer m.Close()

	if uid <= 0 {
		return s.Response(&HammerGopherResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.NOT_LOGIN, nil),
		})
	}

	m.Begin()
	// 1. 从redis读取对应round的数据
	gopherLocations, err := getGopherLocation(m.Redis, req.Round)
	if err != nil {
		return s.Response(&HammerGopherResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.ROUND_QUERY_ERR, err),
		})
	}

	defer func() {
		delGopherLocation(m.Redis, req.Round)
	}()

	// 2. 比较是否击中
	hit := false
	var gopherId uint64 = 0
	for _, v := range gopherLocations.Matrix {
		if req.Index == v.Index {
			hit = true
			gopherId = v.Id
			break
		}
	}

	if !hit {
		return s.Response(&HammerGopherResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.NOT_HAMMER_GOPHER, nil),
		})
	}

	// 3. 计算coins奖励, 更新体力值
	getCoins, err := hammerGopher(m, uid)
	if err != nil {
		return s.Response(&HammerGopherResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.HAMMER_GOPHER_UPDATE, err),
		})
	}

	player, err := t_player.SelectPlayer(m.DB, uid)
	if err != nil {
		return s.Response(&HammerGopherResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	if _, err := t_hammer_entities.CreateOrUpdateCount(m.DB, gopherId, uint64(boost.BoostSet[player.F_boost].F_consume_samina)); err != nil {
		return s.Response(&HammerGopherResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}
	if err := m.Commit(); err != nil {
		return s.Response(&HammerGopherResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	return s.Response(&HammerGopherResponse{
		BaseResponse: protocol.NewBaseResponse(errcode.SUCCESS, nil),
		Coins:        getCoins,
		NokInfo: &protocol.NokInfo{
			Stamina:     config.Instance().NokConfig.MaxStamina - player.F_stamina,
			HitCount:    player.F_hit_count,
			LastHitTime: player.F_last_hit_time.Unix(),
		},
		AccountInfo: &protocol.AccountInfo{
			Coins: player.F_coins,
		},
	})
}

func hammerGopher(m *model.Model, uid int64) (uint64, error) {
	var (
		db = m.DB
	)

	player, err := t_player.SelectPlayer(db, uid)
	if err != nil {
		return 0, err
	}

	// 计算consumeStamina
	boostConfig := boost.BoostSet[player.F_boost]
	consumeStamina := util.Min(config.Instance().NokConfig.MaxStamina-player.F_stamina, int(boostConfig.F_consume_samina))

	// 生成1-5的随机数, 计算获得的coins
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := r.Intn(5) + 1

	// 计算coins
	getCoins := uint64(consumeStamina * randomNum)
	if err := t_player.UpdateHammerGopher(db, uid, getCoins, consumeStamina, boostConfig.F_consume_samina); err != nil {
		return 0, err
	}

	return getCoins, nil
}

func getGopherLocation(rds redis.Conn, round string) (*GopherLocation, error) {
	res, err := cache.GetKey(rds, []interface{}{round})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, errors.New("round not exist or expired")
	}

	val, err := redis.String(res, err)
	if err != nil {
		log.Errorf("faucet value to string err: %v", err)
		return nil, err
	}

	var gl *GopherLocation
	if err := json.Unmarshal([]byte(val), &gl); err != nil {
		return nil, err
	}

	return gl, nil
}

func delGopherLocation(rds redis.Conn, round string) error {
	return cache.DelKey(rds, []interface{}{round})
}
