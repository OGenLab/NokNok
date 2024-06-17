package usermanager

import (
	"errors"

	"github.com/lonng/nano/session"
	"lightspeed.2dao3.com/nokserver/consts/errcode"
	"lightspeed.2dao3.com/nokserver/internal/protocol"
	"lightspeed.2dao3.com/nokserver/model"
	"lightspeed.2dao3.com/nokserver/model/t_player"
)

type BoostRequest struct {
}

type BoostResponse struct {
	protocol.BaseResponse
	Coins       uint64                `json:"coins"`
	AccountInfo *protocol.AccountInfo `json:"accountInfo"`
}

func (u *UserManager) Boost(s *session.Session, req *BoostRequest) error {
	var (
		m   = model.NewModelDefault()
		uid = s.UID()
	)
	defer m.Close()
	if uid <= 0 {
		return s.Response(&BoostResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.NOT_LOGIN, nil),
		})
	}

	// 1. 查询用户信息
	player, err := t_player.SelectPlayer(m.DB, uid)
	if err != nil {
		return s.Response(&BoostResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	// 2. 查询升级消耗
	neededCoins, err := t_player.NextBoost(player.F_boost)
	if err != nil {
		return s.Response(&BoostResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.MAX_BOOST_ERR, err),
		})
	}

	// 3. 判断是否满足条件
	if player.F_coins < neededCoins {
		return s.Response(&BoostResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.UPDATE_NEXT_BOOST_ERR, errors.New("there aren't enough coins")),
		})
	}

	// 4. 更新数据
	if err := t_player.UpdateNextBoost(m.DB, uid, neededCoins); err != nil {
		return s.Response(&BoostResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.UPDATE_NEXT_BOOST_ERR, err),
		})
	}

	return s.Response(&BoostResponse{
		BaseResponse: protocol.NewBaseResponse(errcode.SUCCESS, nil),
		Coins:        neededCoins,
		AccountInfo: &protocol.AccountInfo{
			Boost: player.F_boost + 1,
			Coins: player.F_coins - neededCoins,
		},
	})
}
