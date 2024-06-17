package usermanager

import (
	"github.com/lonng/nano/session"
	"lightspeed.2dao3.com/nokserver/common/boost"
	"lightspeed.2dao3.com/nokserver/consts/errcode"
	"lightspeed.2dao3.com/nokserver/internal/protocol"
	"lightspeed.2dao3.com/nokserver/model"
)

type GetBoostRequest struct {
}

type GetBoostResponse struct {
	protocol.BaseResponse
	BoostInfos []*protocol.BoostInfo `json:"boosts"`
}

func (u *UserManager) GetBoost(s *session.Session, req *GetBoostRequest) error {
	var (
		m   = model.NewModelDefault()
		uid = s.UID()
	)
	defer m.Close()
	if uid <= 0 {
		return s.Response(&GetBoostResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.NOT_LOGIN, nil),
		})
	}

	var boostInfos []*protocol.BoostInfo
	for _, v := range boost.BoostSet {
		boostInfos = append(boostInfos, &protocol.BoostInfo{
			Boost:         v.F_level,
			NeededCoins:   v.F_needed_coins,
			ConsumeSamina: v.F_consume_samina,
			CoinsRate:     v.F_coins_rate,
		})
	}

	return s.Response(&GetBoostResponse{
		BaseResponse: protocol.NewBaseResponse(errcode.SUCCESS, nil),
		BoostInfos:   boostInfos,
	})
}
