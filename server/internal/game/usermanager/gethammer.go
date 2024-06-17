package usermanager

import (
	"sort"

	"github.com/lonng/nano/session"
	"lightspeed.2dao3.com/nokserver/common/entities"
	"lightspeed.2dao3.com/nokserver/consts/errcode"
	"lightspeed.2dao3.com/nokserver/internal/protocol"
	"lightspeed.2dao3.com/nokserver/model"
	"lightspeed.2dao3.com/nokserver/model/t_user_hammer"
)

type GetHammerRequest struct {
}

type GetHammerResponse struct {
	protocol.BaseResponse
	HammerInfos []*protocol.HammerInfo `json:"hammers"`
}

func (u *UserManager) GetHammer(s *session.Session, req *GetHammerRequest) error {
	var (
		m   = model.NewModelDefault()
		uid = s.UID()
	)
	defer m.Close()
	if uid <= 0 {
		return s.Response(&GetHammerResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.NOT_LOGIN, nil),
		})
	}

	// 1. 查询用户锤子表
	hammers, err := t_user_hammer.SelectHammer(m.DB, uid)
	if err != nil {
		return s.Response(&GetHammerResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	var hammerInfos []*protocol.HammerInfo
	for _, v := range hammers {
		hammerInfos = append(hammerInfos, &protocol.HammerInfo{
			EntitiesInfo: protocol.EntitiesInfo{
				Id: v.F_hammer_id,
			},
			EquipmentCount: v.F_equipment_count,
			Count:          v.F_count,
		})
	}

	hammerSet := entities.EntSet.HammerSet()
	sort.Slice(hammerInfos, func(i, j int) bool {
		eCountI := hammerInfos[i].EquipmentCount
		eCountJ := hammerInfos[j].EquipmentCount
		qualityI := entities.TypExtToHammerQuality[hammerSet[uint64(hammerInfos[i].Id)].F_type_ext]
		qualityJ := entities.TypExtToHammerQuality[hammerSet[uint64(hammerInfos[j].Id)].F_type_ext]
		return eCountI > eCountJ || (eCountI == eCountJ && qualityI > qualityJ)
	})

	return s.Response(&GetHammerResponse{
		BaseResponse: protocol.NewBaseResponse(errcode.SUCCESS, nil),
		HammerInfos:  hammerInfos,
	})
}
