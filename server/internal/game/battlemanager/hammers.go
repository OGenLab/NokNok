package battlemanager

import (
	"github.com/lonng/nano/session"
	"lightspeed.2dao3.com/nokserver/common/entities"
	"lightspeed.2dao3.com/nokserver/consts/errcode"
	"lightspeed.2dao3.com/nokserver/internal/protocol"
	"lightspeed.2dao3.com/nokserver/model"
)

type HammersRequest struct {
}

type HammersResponse struct {
	protocol.BaseResponse
	Hammers []*protocol.HammerInfo `json:"hammers"`
}

func (u *BattleManager) Hammers(s *session.Session, req *HammersRequest) error {
	var (
		m         = model.NewModelDefault()
		uid       = s.UID()
		hammerSet = entities.EntSet.HammerSet()
	)
	defer m.Close()

	if uid <= 0 {
		return s.Response(&HammersResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.NOT_LOGIN, nil),
		})
	}

	hammers := make([]*protocol.HammerInfo, 0, len(hammerSet))
	for _, hammer := range hammerSet {
		hammers = append(hammers, &protocol.HammerInfo{
			EntitiesInfo: protocol.EntitiesInfo{
				Id:          hammer.F_entities_id,
				Season:      hammer.F_season,
				Name:        hammer.F_name,
				ImageUrl:    hammer.F_image_url,
				Probability: hammer.F_probability,
			},
			Quality: entities.TypExtToHammerQuality[hammer.F_type_ext],
		})
	}

	return s.Response(&HammersResponse{
		BaseResponse: protocol.NewBaseResponse(errcode.SUCCESS, nil),
		Hammers:      hammers,
	})
}
