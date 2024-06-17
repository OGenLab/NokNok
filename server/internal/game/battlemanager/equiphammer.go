package battlemanager

import (
	"github.com/lonng/nano/session"
	"lightspeed.2dao3.com/nokserver/consts/errcode"
	"lightspeed.2dao3.com/nokserver/internal/protocol"
	"lightspeed.2dao3.com/nokserver/model"
	"lightspeed.2dao3.com/nokserver/model/t_user_hammer"
)

type EquipHammerRequest struct {
	HammerId uint64 `json:"hammerId"`
}

type EquipHammerResponse struct {
	protocol.BaseResponse
}

func (u *BattleManager) EquipHammer(s *session.Session, req *EquipHammerRequest) error {
	var (
		m   = model.NewModelDefault()
		uid = s.UID()
	)
	defer m.Close()

	if uid <= 0 {
		return s.Response(&EquipHammerResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.NOT_LOGIN, nil),
		})
	}

	m.Begin()

	hammer, err := t_user_hammer.SelectHammerById(m.DB, uid, req.HammerId)
	if err != nil {
		return s.Response(&EquipHammerResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	} else if hammer == nil {
		return s.Response(&EquipHammerResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.PARAMETER_INVALID, nil),
		})
	}

	if hammer.F_equipment_count == 1 {
		return s.Response(&EquipHammerResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.SUCCESS, nil),
		})
	}

	if err := t_user_hammer.NotEquipHammer(m.DB, uid); err != nil {
		return s.Response(&EquipHammerResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}
	if err := t_user_hammer.EquipHammer(m.DB, uid, req.HammerId); err != nil {
		return s.Response(&EquipHammerResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	if err := m.Commit(); err != nil {
		return s.Response(&EquipHammerResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	return s.Response(&EquipHammerResponse{
		BaseResponse: protocol.NewBaseResponse(errcode.SUCCESS, nil),
	})
}
